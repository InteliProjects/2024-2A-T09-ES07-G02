package main

import (
	"encoding/base64"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"webscrapping/kafka"
	"webscrapping/logs"
	"webscrapping/s33"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gocolly/colly"
	//"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

var (
	startTime time.Time
	s3Client  *s3.Client
	bucket    string
	stopScan  = false
	duration  = 3 * time.Minute // Define the time limit for the scraping process (5 minutes in this case)
)

// encodeFileToBase64 encodes the contents of the file at filePath into a Base64 string.
func encodeFileToBase64(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// Log the error for better tracking
		logs.LogError(fmt.Errorf("failed to open file: %v", err))
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file's content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logs.LogError(fmt.Errorf("failed to read file: %v", err))
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Encode the file's content to Base64
	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	return encoded, nil
}

func main() {
	// Load environment variables
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	brokers := []string{"kafka:29092"}
	err := kafka.InitKafkaProducer(brokers)
	if err != nil {
		log.Fatalln("Failed to initialize Kafka producer:", err)
	}

	awsRegion := "us-east-1"
	bucket = "law-hunters1"

	// Configs for s3
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		log.Fatalf("Unable to load AWS config, %v", err)
		return
	}
	s3Client = s3.NewFromConfig(cfg)
	s33.InitS3(s3Client, bucket)

	// Initialize logs
	logs.InitLogger()

	c := cron.New(cron.WithLocation(time.FixedZone("BRT", -3*60*60)))

	_, err = c.AddFunc("0 0 * * *", func() {
		fmt.Println("Executando web scraping à meia-noite...")
		err := executeScraping()
		if err != nil {
			log.Println("Erro ao executar o web scraping:", err)
		}
	})
	if err != nil {
		log.Fatalf("Erro ao agendar a tarefa: %v", err)
	}

	c.Start()

	// Configure the HTTP server with timeout settings
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      http.HandlerFunc(scrapeHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Start the server on port 8081
	fmt.Println("Server started on port 8081")
	logs.LogInfo("Server started on port 8081")
	log.Fatal(srv.ListenAndServe())
}

func executeScraping() error {
	stopScan = false
	startTime = time.Now()
	fmt.Println("Scraping in progress...")


	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)


	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10, 
	})


	go func() {
		time.Sleep(duration) 
		fmt.Println("Time limit exceeded. Stopping scraper and uploading files to S3...")
		stopScan = true
		c.Wait() 
		err := uploadAllFilesToS3()
		if err != nil {
			log.Println("Error uploading files to S3:", err)
			return
		}
		fmt.Println("Files successfully uploaded to S3.")
	}()

	c.OnHTML("div.newsList a[href]", func(e *colly.HTMLElement) {
		if stopScan {
			return
		}
		link := e.Attr("href")
		fmt.Println("News found:", link)
		logs.LogInfo(link)

		if !stopScan {
			e.Request.Visit(link)
		}
	})


	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if stopScan {
			return
		}
		link := e.Attr("href")
		if e.Text == "clique aqui" {
			fmt.Println("File link found:", link)
			logs.LogInfo(fmt.Sprintf("File link found: %s", link))
			err := saveFile(link)
			if err != nil {
				logs.LogError(err)
				log.Println("Error saving file:", err)
			}
		}
	})


	err := c.Visit("https://www.apimecbrasil.com.br/noticias/apimec-brasil/")
	if err != nil {
		logs.LogError(err)
		return fmt.Errorf("error visiting page: %v", err)
	}


	c.Wait()

	fmt.Println("Scraping process completed!")
	return nil
}

// Handler for the / endpoint
func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	startTime = time.Now() // Mark the start time of the execution
	fmt.Fprintln(w, "Scraping in progress...")

	// Create a new Colly collector
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)

	// Limit the number of concurrent requests
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10, // Execute one request at a time for simplicity
	})

	// Start a goroutine to monitor the time and stop scraping after the duration
	go func() {
		time.Sleep(duration) // Wait for the duration to pass
		fmt.Println("Time limit exceeded. Stopping scraper and uploading files to S3...")
		stopScan = true
		c.Wait() // Stop the scraper
		err := uploadAllFilesToS3()
		if err != nil {
			log.Println("Error uploading files to S3:", err)
			return
		}
		fmt.Println("Files successfully uploaded to S3.")
	}()

	// Callback for extracting news links
	c.OnHTML("div.newsList a[href]", func(e *colly.HTMLElement) {
		if stopScan {
			return
		}
		link := e.Attr("href")
		fmt.Println("News found:", link)
		logs.LogInfo(link)

		//kafka.SendMessageToKafka("scrapped-data", link)

		// Visit the found link
		if !stopScan {
			e.Request.Visit(link)
		}
	})

	// Callback for extracting file links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if stopScan {
			return
		}
		link := e.Attr("href")
		if e.Text == "clique aqui" {
			fmt.Println("File link found:", link)
			logs.LogInfo(fmt.Sprintf("File link found: %s", link))
			err := saveFile(link)
			if err != nil {
				logs.LogError(err)
				log.Println("Error saving file:", err)
			}
		}
	})

	// Visit the starting URL
	err := c.Visit("https://www.apimecbrasil.com.br/noticias/apimec-brasil/")
	if err != nil {
		http.Error(w, "Error visiting page", http.StatusInternalServerError)
		logs.LogError(err)
		return
	}

	// Wait for scraping to finish
	c.Wait()

	// After scraping finishes, inform the user
	fmt.Fprintln(w, "Scraping process completed!")
}

// Function to save files locally
func saveFile(url string) error {
	// Get the file name from the URL
	fileName := filepath.Base(url)

	// Sanitize the file name
	fileName = sanitizeFileName(fileName)

	// Check the file extension; ignore if not PDF
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".pdf" {
		fmt.Println("Ignoring non-PDF file:", fileName)
		return nil
	}

	// Create the "downloads" folder if it does not exist
	err := os.MkdirAll("downloads", os.ModePerm)
	if err != nil {
		return err
	}

	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Make the HTTP request to get the file
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file in the "downloads" folder
	filePath := filepath.Join("downloads", fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the HTTP response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("PDF file saved:", filePath)
	logs.LogInfo(fmt.Sprintf("PDF file saved: %s", filePath))

	encodedFile, err := encodeFileToBase64(filePath)
	if err != nil {
		return err
	}
	kafka.SendEncodedFileToKafka("scrapped-files", filePath, encodedFile)

	return nil
}

// Function to sanitize file names by replacing invalid characters
func sanitizeFileName(fileName string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*]+`)
	return re.ReplaceAllString(fileName, "_")
}

// Function to upload all downloaded files to S3
func uploadAllFilesToS3() error {
	files, err := os.ReadDir("downloads")
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join("downloads", file.Name())
		fileContent, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileContent.Close()

		_, err = s33.UploadFileToS3(file.Name(), fileContent)
		if err != nil {
			return err
		}
		logs.LogInfo(fmt.Sprintf("File uploaded to S3: %s", file.Name()))
	}

	return nil
}

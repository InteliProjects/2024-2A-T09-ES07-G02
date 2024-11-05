package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

var (
	logger   *log.Logger
	producer sarama.SyncProducer
)

// KafkaLogWriter is a custom log writer that sends log entries to a Kafka topic.
type KafkaLogWriter struct {
	topic    string
	producer sarama.SyncProducer
}

func (w *KafkaLogWriter) Write(p []byte) (n int, err error) {
	if w.producer == nil {
		return 0, fmt.Errorf("Kafka producer is nil") // Check if producer is nil
	}
	msg := &sarama.ProducerMessage{
		Topic: w.topic,
		Value: sarama.StringEncoder(p),
	}

	_, _, err = w.producer.SendMessage(msg)
	if err != nil {
		return 0, err // Return an error if Kafka message fails
	}
	return len(p), nil
}

// Middleware to enable CORS for all routes
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Global S3 client initialized for reuse
var s3Client *s3.Client
var logFile *os.File

// init function is called before main() to initialize AWS SDK, Kafka, and log files
func init() {
	// Ensure that the logs directory exists
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755) // Create logs directory if it doesn't exist
		if err != nil {
			log.Fatalf("Error creating log directory: %v", err)
		}
	}

	// Open or create the log file
	var err error
	logFile, err = os.OpenFile(filepath.Join(logDir, "backend_logs.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Publish error to Kafka if log file cannot be opened
		if producer != nil {
			kafkaLogWriter := &KafkaLogWriter{
				topic:    "SpeechService-logs",
				producer: producer,
			}
			kafkaLogWriter.Write([]byte(fmt.Sprintf("Error opening log file: %v", err)))
		}
		log.Fatalf("Error opening log file: %v", err)
	}

	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Error loading SDK configuration: %v", err)
		return
	}

	// Enable detailed AWS logging
	cfg.ClientLogMode = aws.LogRequestWithBody | aws.LogResponseWithBody
	cfg.Logger = logging.NewStandardLogger(os.Stdout)

	// Initialize S3 client
	s3Client = s3.NewFromConfig(cfg)
	log.Println("S3 client initialized successfully")
}

// uploadFileToS3 uploads a file to the specified S3 bucket
func uploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, bucketName string) (string, error) {
	// Generate a unique filename based on the current timestamp
	fileKey := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), fileHeader.Filename)
	// Perform the file upload to the S3 bucket
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileKey),
		Body:        file,
		ContentType: aws.String(http.DetectContentType([]byte(fileHeader.Filename))),
	}, func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions, middleware.AddUserAgentKeyValue("Authorization-Logger", "v1.0"))
	})
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}
	return fileKey, nil // Return the file key in S3
}

// Handles audio and transcript upload, saving them to S3
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload request received")
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the audio file from the form
	audioFile, audioHeader, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Failed to retrieve the audio file", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()
	// Get the transcript from the form
	transcript := r.FormValue("transcript")
	if transcript == "" {
		http.Error(w, "Transcript is required", http.StatusBadRequest)
		return
	}
	// Name of the S3 bucket
	bucketName := "speech-service-inteli"
	// Upload the audio file to S3
	fileKey, err := uploadFileToS3(audioFile, audioHeader, bucketName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to upload the file", http.StatusInternalServerError)
		return
	}
	// Generate a response with the S3 file key and the transcript
	response := fmt.Sprintf("File uploaded successfully: %s\nTranscript: %s", fileKey, transcript)
	w.Write([]byte(response))
}

// logHandler handles log data sent from the frontend and writes it to a log file.
func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming log data
	var logMessage struct {
		Message string `json:"message"`
		Level   string `json:"level"`
	}

	err := json.NewDecoder(r.Body).Decode(&logMessage)
	if err != nil {
		http.Error(w, "Failed to parse log data", http.StatusBadRequest)
		return
	}

	// Format the log entry
	formattedLog := fmt.Sprintf("%s - %s: %s\n", time.Now().Format(time.RFC3339), logMessage.Level, logMessage.Message)

	// Write the log entry to the log file
	_, err = logFile.WriteString(formattedLog)
	if err != nil {
		http.Error(w, "Failed to write log to file", http.StatusInternalServerError)
		return
	}

	// Respond to the frontend with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Log received and saved"))
}

func main() {
	// Set up Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	brokers := []string{"kafka:29092"}
	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	
	if producer != nil {
		kafkaLogWriter := &KafkaLogWriter{
			topic:    "SpeechService-logs",
			producer: producer,
		}
		_, err := kafkaLogWriter.Write([]byte("Application is fully running"))
		if err != nil {
			log.Printf("Error publishing 'Application is fully running' log to Kafka: %v", err)
		}
	}

  if err != nil {
		log.Panic("Error creating kafka producer:", err)
	}
	defer producer.Close()

	// Map the "/upload" endpoint to the uploadHandler function
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/log", logHandler) // Endpoint to receive logs from the frontend

	log.Println("Server started on port :6969")
	log.Fatal(http.ListenAndServe(":6969", enableCors(mux)))
}


package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
  "time"
  "context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Postgres string
	KafkaBrokers string
	KafkaTopic string
}

func createKafkaTopic(broker string, topic string, partitions int, replicationFactor int) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return fmt.Errorf("Failed to create Kafka admin client: %v", err)
	}
	defer adminClient.Close()

	topicSpecification := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{topicSpecification},
	)
	if err != nil {
		return fmt.Errorf("Failed to create topic: %v", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("Failed to create topic: %v", result.Error)
		}
		fmt.Printf("Topic '%s' created successfully\n", result.Topic)
	}

	return nil
}

func createKafkaProducer(brokers string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create Kafka producer: %v", err)
	}
	return p, nil
}

func sendLogToKafka(p *kafka.Producer, topic string, message string) {
    deliveryChan := make(chan kafka.Event)
    defer close(deliveryChan)

    err := p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(message),
    }, deliveryChan)

    if err != nil {
        log.Printf("Failed to produce message: %v", err)
        return
    }

    // Aguardar a entrega da mensagem
    e := <-deliveryChan

    // Verificar se o evento é do tipo *kafka.Message
    msg, ok := e.(*kafka.Message)
    if !ok {
        log.Printf("Failed to deliver message: unexpected event type")
        return
    }

    // Verificar se houve algum erro ao enviar a mensagem
    if msg.TopicPartition.Error != nil {
        log.Printf("Failed to deliver message: %v", msg.TopicPartition.Error)
    } else {
        fmt.Printf("Message delivered to topic %s [partition %d] at offset %v\n",
            *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
    }
}

func consumeFromKafka(brokers, topic string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "synonym-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	fmt.Println("Consuming messages from Kafka...")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message from Kafka: %s\n", string(msg.Value))

			// Process the message (it could be a word to generate synonyms for)
			tag := string(msg.Value)
			synonyms, err := generateSynonymsForTag(tag)
			if err != nil {
				log.Printf("Failed to generate synonyms for tag %s: %v", tag, err)
			} else {
				fmt.Printf("Synonyms for '%s': %s\n", tag, synonyms)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func queryOllama(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate" // standard url for ollama API

	requestBody, err := json.Marshal(OllamaRequest{
		Model:  "llama3.1",
		Prompt: prompt,
	})

	if err != nil {
		return "", fmt.Errorf("Failed to marshall request: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder // to accumulate the response

	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var partialResponse map[string]interface{}

		err := decoder.Decode(&partialResponse)
		if err != nil {
			return "", fmt.Errorf("Failed to decode response: %v", err)
		}

		fmt.Printf("Decoded chunk: %v\n", partialResponse) // Debugging

		if respChunk, exists := partialResponse["response"]; exists {
			fullResponse.WriteString(fmt.Sprintf("%v", respChunk))
		}
	}

	return fullResponse.String(), nil
}

type OllamaResponse struct {
	Completion string `json:"completion"`
}

func getAllTagsAndSubtags(db *sql.DB) ([]string, error) {
	var tags []string

	rows, err := db.Query("SELECT tag FROM tag")
	if err != nil {
		return nil, fmt.Errorf("Failed to query tags: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("Failed to query subtags: %v", err)
		}
		tags = append(tags, tag)
	}

	subtagRows, err := db.Query("SELECT name FROM subtag")
	if err != nil {
		return nil, fmt.Errorf("Failed to query subtags: %v", err)
	}
	defer subtagRows.Close()

	for subtagRows.Next() {
		var subtag string
		if err := subtagRows.Scan(&subtag); err != nil {
			return nil, fmt.Errorf("Failed to scan subtag: %v", err)
		}
		tags = append(tags, subtag)
	}

	return tags, nil
}

func generateSynonymsForTag(tag string) (string, error) {
	prompt := fmt.Sprintf("Gere uma lista de sinônimos para a palavra, e não diga absolutamente mais nada além da lista de sinônimos: '%s'.", tag)

	response, err := queryOllama(prompt)
	if err != nil {
		return "", fmt.Errorf("Failed to generate synonyms for %s: %v", tag, err)
	}

	return response, nil
}

func processTagsAndGenerateSynonyms(db *sql.DB) (map[string]string, error) {
	tags, err := getAllTagsAndSubtags(db)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tags and subtags: %v", err)
	}

	synonymsMap := make(map[string]string)

	for _, tag := range tags {
		synonyms, err := generateSynonymsForTag(tag)
		if err != nil {
			log.Printf("Error generating synonyms for tag %s: %v", tag, err)
			continue
		}

		fmt.Println(synonyms)
		synonymsMap[tag] = synonyms
	}

	return synonymsMap, nil
}

// Load .env
func initConfig() Config {
	godotenv.Load()

	return Config{
		Postgres: getEnv("POSTGRES_CONNECTION", ""),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "synonym_topic"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// TreeNode represents a node in the tag/subtag hierarchy, holding synonyms and child nodes.
type TreeNode struct {
	Synonyms map[string]bool      // A set of synonyms for this tag/subtag.
	Children map[string]*TreeNode // A map of child nodes (subtags).
}

// NewTreeNode creates and initializes a new TreeNode with the given synonyms.
func NewTreeNode(synonyms ...string) *TreeNode {
	node := &TreeNode{
		Synonyms: make(map[string]bool),
		Children: make(map[string]*TreeNode),
	}
	for _, synonym := range synonyms {
		node.Synonyms[synonym] = true
	}
	return node
}

// RemoveEmptySubTags recursively removes child nodes that have no synonyms and no further children.
func (n *TreeNode) RemoveEmptySubTags() {
	for subtag, child := range n.Children {
		child.RemoveEmptySubTags()
		if len(child.Synonyms) == 0 && len(child.Children) == 0 {
			delete(n.Children, subtag)
		}
	}
}

// AddSubTag adds a subtag to the current node. If it already exists, it adds synonyms to it.
func (n *TreeNode) AddSubTag(subtag string, synonyms ...string) *TreeNode {
	child, exists := n.Children[subtag]
	if !exists {
		child = NewTreeNode(synonyms...)
		n.Children[subtag] = child
	} else {
		for _, synonym := range synonyms {
			child.Synonyms[synonym] = true
		}
	}
	return child
}

// AddSubTagsPath adds a series of subtags along a path, creating nodes as necessary.
func (n *TreeNode) AddSubTagsPath(path []string, synonyms ...string) *TreeNode {
	current := n
	for i, subtag := range path {
		if i == len(path)-1 {
			current = current.AddSubTag(subtag, synonyms...)
		} else {
			current = current.AddSubTag(subtag)
		}
	}
	return current
}

// CountOccurrences counts how many times the synonyms of this node and its children appear in the content.
func (n *TreeNode) CountOccurrences(content string) (int, string) {
	count := 0
	used := make(map[string]bool)

	for synonym := range n.Synonyms {
		if strings.Contains(content, synonym) && !used[synonym] {
			count += strings.Count(content, synonym)
			used[synonym] = true
		}
	}

	maxCount := count
	maxName := ""

	for subtag, child := range n.Children {
		childCount, childBestSubtag := child.CountOccurrences(content)
		count += childCount

		if childCount > maxCount {
			maxCount = childCount
			maxName = subtag
		}

		if childBestSubtag != "" {
			maxName = childBestSubtag
		}
	}

	return count, maxName
}

// findBestSynonym finds the best matching synonym for this node within the content.
func (n *TreeNode) findBestSynonym(content string) string {
	for synonym := range n.Synonyms {
		if strings.Contains(content, synonym) {
			return synonym
		}
	}
	return ""
}

// findBestMatch finds the tag or subtag with the most occurrences in the content.
func findBestMatch(tree map[string]*TreeNode, content string) (string, bool, int) {
	maxOccurrences := 0
	bestMatch := ""
	isSubtag := false

	for tagName, tree := range tree {
		count, bestSubtag := tree.CountOccurrences(content)
		if count > maxOccurrences {
			maxOccurrences = count
			bestMatch = bestSubtag
			isSubtag = (bestSubtag != "")
			if !isSubtag {
				bestMatch = tagName
			}
		}
	}
	return bestMatch, isSubtag, maxOccurrences
}

// findTopMatches finds the top N tags or subtags with the most occurrences in the content.
func findTopMatches(tree map[string]*TreeNode, content string, topN int) []struct {
	Name     string
	IsSubtag bool
	Count    int
} {
	matches := make([]struct {
		Name     string
		IsSubtag bool
		Count    int
	}, 0)

	for tagName, tree := range tree {
		count, bestSubtag := tree.CountOccurrences(content)
		if bestSubtag != "" {
			matches = append(matches, struct {
				Name     string
				IsSubtag bool
				Count    int
			}{bestSubtag, true, count})
		} else {
			matches = append(matches, struct {
				Name     string
				IsSubtag bool
				Count    int
			}{tagName, false, count})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Count > matches[j].Count
	})

	if len(matches) > topN {
		return matches[:topN]
	}
	return matches
}

// associateLLRWithBestMatches associates an LRR (record) with the best matching tags or subtags in the database.
func associateLLRWithBestMatches(db *sql.DB, lrrID int, matches []struct {
	Name     string
	IsSubtag bool
	Count    int
}) {
	for _, match := range matches {
		var err error
		if match.IsSubtag {
			_, err = db.Exec("INSERT INTO lrr_tag (lrr_id, subtag_id) VALUES ($1, (SELECT id FROM subtag WHERE name = $2))", lrrID, match.Name)
		} else {
			_, err = db.Exec("INSERT INTO lrr_tag (lrr_id, tag_id) VALUES ($1, (SELECT id FROM tag WHERE tag = $2))", lrrID, match.Name)
		}
		if err != nil {
			log.Fatalf("Failed to insert LRR association: %v", err)
		}
		fmt.Printf("Associated LRR %d with %s (isSubtag: %t)\n", lrrID, match.Name, match.IsSubtag)
	}
}

func (n *TreeNode) PrintTree(indent string, nodeName string, isLast bool) {
	fmt.Printf("%s└── [%s]:\n", indent, nodeName)
	newIndent := indent + "    "
	if !isLast {
		newIndent = indent + "│   "
	}

	synonyms := make([]string, 0, len(n.Synonyms))
	for synonym := range n.Synonyms {
		synonyms = append(synonyms, synonym)
	}

	for i, synonym := range synonyms {
		if i == len(synonyms)-1 && len(n.Children) == 0 {
			fmt.Printf(newIndent+"└── [Synonym]: %s\n", synonym)
		} else {
			fmt.Printf(newIndent+"├── [Synonym]: %s\n", synonym)
		}
	}

	children := make([]string, 0, len(n.Children))
	for subtag := range n.Children {
		children = append(children, subtag)
	}

	for i, child := range children {
		isLastChild := i == len(children)-1
		n.Children[child].PrintTree(newIndent, child, isLastChild)
	}
}

// buildTreeFromDB builds a tree structure of tags and subtags by querying the database.
func buildTreeFromDB(db *sql.DB) map[string]*TreeNode {
	tags, err := db.Query("SELECT id, tag FROM tag")
	if err != nil {
		log.Fatalf("Failed to query tags: %v", err)
	}
	defer tags.Close()

	tree := make(map[string]*TreeNode)

	tagIDMap := make(map[int]string)
	for tags.Next() {
		var id int
		var tagName string
		if err := tags.Scan(&id, &tagName); err != nil {
			log.Fatalf("Failed to scan tag: %v", err)
		}
		tree[tagName] = NewTreeNode()
		tagIDMap[id] = tagName
	}

	subtags, err := db.Query("SELECT id, name, parent_tag_id, parent_subtag_id FROM subtag")
	if err != nil {
		log.Fatalf("Failed to query subtags: %v", err)
	}
	defer subtags.Close()

	subtagNodeMap := make(map[int]*TreeNode)

	for subtags.Next() {
		var id int
		var name string
		var parentTagID, parentSubtagID sql.NullInt64

		if err := subtags.Scan(&id, &name, &parentTagID, &parentSubtagID); err != nil {
			log.Fatalf("Failed to scan subtag: %v", err)
		}

		var parentNode *TreeNode
		if parentTagID.Valid {
			parentNode = tree[tagIDMap[int(parentTagID.Int64)]]
		} else if parentSubtagID.Valid {
			parentNode = subtagNodeMap[int(parentSubtagID.Int64)]
		}

		if parentNode != nil {
			subtagNode := parentNode.AddSubTag(name)
			subtagNodeMap[id] = subtagNode
		}
	}

	synonyms, err := db.Query("SELECT synonym, parent_tag_id, parent_subtag_id FROM synonym")
	if err != nil {
		log.Fatalf("Failed to query synonyms: %v", err)
	}
	defer synonyms.Close()

	for synonyms.Next() {
		var synonym string
		var parentTagID, parentSubtagID sql.NullInt64

		if err := synonyms.Scan(&synonym, &parentTagID, &parentSubtagID); err != nil {
			log.Fatalf("Failed to scan synonym: %v", err)
		}

		var parentNode *TreeNode
		if parentTagID.Valid {
			parentNode = tree[tagIDMap[int(parentTagID.Int64)]]
		} else if parentSubtagID.Valid {
			parentNode = subtagNodeMap[int(parentSubtagID.Int64)]
		}

		if parentNode != nil {
			parentNode.Synonyms[synonym] = true
		}
	}

	return tree
}

func printAssociatedTagOrSubtag(db *sql.DB, lrrID int) {
	var tagName, subtagName sql.NullString

	err := db.QueryRow(`
		SELECT t.tag, s.name
		FROM lrr_tag lt
		LEFT JOIN tag t ON lt.tag_id = t.id
		LEFT JOIN subtag s ON lt.subtag_id = s.id
		WHERE lt.lrr_id = $1
	`, lrrID).Scan(&tagName, &subtagName)
	if err != nil {
		log.Fatalf("Failed to retrieve associated tag or subtag: %v", err)
	}

	if tagName.Valid {
		fmt.Printf("LRR %d is associated with Tag: %s\n", lrrID, tagName.String)
	} else if subtagName.Valid {
		fmt.Printf("LRR %d is associated with Subtag: %s\n", lrrID, subtagName.String)
	} else {
		fmt.Printf("LRR %d has no association\n", lrrID)
	}
}

func main() {
	app := fiber.New()

	config := initConfig()
	logTopic := "tag-service-logs"

	err := createKafkaTopic(config.KafkaBrokers, logTopic, 1, 1)
	if err != nil {
		log.Fatalf("Failed to create Kafka topic: %v", err)
	}

	producer, err := createKafkaProducer(config.KafkaBrokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

  serviceStartMsg := "Service started successfully at " + time.Now().Format(time.RFC3339)
  sendLogToKafka(producer, logTopic, serviceStartMsg)
	go consumeFromKafka(config.KafkaBrokers, config.KafkaTopic)

	app.Get("/tag", func(c *fiber.Ctx) error {
		llr := "llr.txt"
		dbName := config.Postgres

		db, err := sql.Open("postgres", dbName)
		if err != nil {
			logMsg := fmt.Sprintf("Failed to open DB: %s", err)
			sendLogToKafka(producer, logTopic, logMsg)
			return c.Status(500).SendString(logMsg)
		}
		defer db.Close()

		content, err := os.ReadFile(llr)
		if err != nil {
			logMsg := fmt.Sprintf("Error reading file: %v", err)
			sendLogToKafka(producer, logTopic, logMsg)
			return c.Status(500).SendString(logMsg)
		}

		tree := buildTreeFromDB(db)
		topMatches := findTopMatches(tree, string(content), 3)
		if len(topMatches) > 0 {
			associateLLRWithBestMatches(db, 1, topMatches)
			printAssociatedTagOrSubtag(db, 1)
			return c.SendString("We gucci\n")
		} else {
			logMsg := "No correspondence found."
			sendLogToKafka(producer, logTopic, logMsg)
			return c.SendString(logMsg + "\n")
		}
	})

	err = app.Listen(":3000")
	if err != nil {
		logMsg := fmt.Sprintf("Failed to start server: %v", err)
		sendLogToKafka(producer, logTopic, logMsg)
		log.Fatalf(logMsg)
	}
}

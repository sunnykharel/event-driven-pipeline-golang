package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	NUM_VCPUS                 = 3
	DYNAMO_DB_MAX_BULK_INSERT = 25
	UNKNOWN_DOMAIN            = "empty_domain"
	UNKNOWN_EMAIL             = "empty_email"
)

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type CompromisedCredential struct {
	Email    string `dynamodbav:"email"`
	Username string `dynamodbav:"username"`
	Password string `dynamodbav:"password"`
	Domain   string `dynamodbav:"domain"`
	Id       string `dynamodbav:"id"`
}

func fetchFileFromS3(ctx context.Context, client *s3.Client, bucket, key string) ([]byte, error) {
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching object from S3: %w", err)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading file content: %w", err)
	}
	return buf.Bytes(), nil
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

func bulkInsertToDynamoDB(ctx context.Context, basics TableBasics, credentials []CompromisedCredential) error {
	items := make([]map[string]types.AttributeValue, len(credentials))
	for i, credential := range credentials {
		item, err := attributevalue.MarshalMap(credential)
		if err != nil {
			return fmt.Errorf("failed to marshal credential: %w", err)
		}
		items[i] = item
	}
	for _, item := range items {
		_, err := basics.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basics.TableName),
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed to insert batch: %w", err)
		}
	}
	return nil
}

func parseEmailUsernameDomainPasswordFromLine(line string) (string, string, string, string) {
	line = strings.ReplaceAll(line, ",", ":")
	line = strings.ReplaceAll(line, ";", ":")

	parts := strings.Split(line, ":")
	email, username, domain, password := "", "", "", ""

	if len(parts) == 2 {
		if strings.Contains(parts[0], "@") {
			email = parts[0]
			password = parts[1]
			username = strings.Split(email, "@")[0]
			domain = strings.Split(email, "@")[1]
		} else {
			password = parts[0]
			email = parts[1]
			if strings.Contains(email, "@") {
				username = strings.Split(email, "@")[0]
				domain = strings.Split(email, "@")[1]
			}
		}
	} else if len(parts) == 1 && strings.Contains(parts[0], "@") {
		email = parts[0]
		username = strings.Split(email, "@")[0]
		domain = strings.Split(email, "@")[1]
	}
	if domain == "" {
		domain = UNKNOWN_DOMAIN
	}
	if email == "" {
		email = UNKNOWN_EMAIL
	}

	return email, username, domain, password
}

func worker(ctx context.Context, wg *sync.WaitGroup, basics TableBasics, taskChan <-chan []string) {
	defer wg.Done()
	for lines := range taskChan {
		var credentials []CompromisedCredential
		for _, line := range lines {
			email, username, domain, password := parseEmailUsernameDomainPasswordFromLine(line)
			credentials = append(credentials, CompromisedCredential{
				Email:    email,
				Username: username,
				Password: hashPassword(password),
				Domain:   domain,
				Id:       uuid.New().String(),
			})
		}
		if err := bulkInsertToDynamoDB(ctx, basics, credentials); err != nil {
			log.Printf("Failed to bulk insert: %v", err)
		} else {
			log.Printf("Successfully inserted %d records", len(credentials))
		}
	}
}

func processAndInsert(ctx context.Context, basics TableBasics, content []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	taskChan := make(chan []string, NUM_VCPUS)
	var wg sync.WaitGroup

	// Start worker threads
	for i := 0; i < NUM_VCPUS; i++ {
		wg.Add(1)
		go worker(ctx, &wg, basics, taskChan)
	}

	var batch []string
	for scanner.Scan() {
		batch = append(batch, scanner.Text())
		if len(batch) == DYNAMO_DB_MAX_BULK_INSERT {
			taskChan <- batch
			batch = nil
		}
	}
	// Send remaining records
	if len(batch) > 0 {
		taskChan <- batch
	}

	close(taskChan)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file content %w", err)
	}
	return nil
}

func handleS3Event(ctx context.Context, s3Event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	dynamoClient := dynamodb.NewFromConfig(cfg)
	basics := TableBasics{
		DynamoDbClient: dynamoClient,
		TableName:      os.Getenv("COMPROMISEDCREDENTIALS_TABLE_NAME"),
	}
	if basics.TableName == "" {
		return fmt.Errorf("COMPROMISEDCREDENTIALS_TABLE_NAME environment variable is not set")
	}

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		log.Printf("Processing file from bucket: %s, key: %s", bucket, key)

		fileContent, err := fetchFileFromS3(ctx, s3Client, bucket, key)
		if err != nil {
			return fmt.Errorf("failed to fetch file from S3: %w", err)
		}

		err = processAndInsert(ctx, basics, fileContent)
		if err != nil {
			return fmt.Errorf("failed to process file content %w", err)
		}
	}
	return nil
}

func main() {
	lambda.Start(handleS3Event)
}

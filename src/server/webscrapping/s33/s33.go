package s33

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucket string

// InitS3 initializes the S3 client and sets the bucket name
func InitS3(client *s3.Client, bucketName string) {
	s3Client = client
	bucket = bucketName
}

// UploadFileToS3 uploads the file content to a specified S3 bucket.
func UploadFileToS3(fileName string, fileContent io.Reader) (string, error) {
	fileKey := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), filepath.Base(fileName))

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
		Body:   fileContent,
	})
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}

	return fileKey, nil
}

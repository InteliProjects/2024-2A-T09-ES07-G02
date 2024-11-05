package main

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

// Configuração para inicializar o cliente S3 no teste
func setupS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return s3.NewFromConfig(cfg)
}

func TestUploadHandler(t *testing.T) {
	// Configura o cliente S3 e outros elementos necessários
	s3Client = setupS3Client()
	bucketName := "speech-service-inteli"

	// Criação de um arquivo fake para testar o upload
	audioContent := []byte("test audio content")

	// Criação do formulário multipart para a requisição HTTP
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", "test_audio.wav")
	part.Write(audioContent)
	writer.WriteField("transcript", "This is a test transcription")
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// Grava o tempo antes de enviar a requisição
	startTime := time.Now().UTC()

	// Grava a resposta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uploadHandler)
	handler.ServeHTTP(rr, req)

	// Verifica se a resposta é 200 OK
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")

	// Gera a chave esperada do arquivo com base na lógica de nomeação
	fileKey := fmt.Sprintf("uploads/%d-%s", startTime.Unix(), "test_audio.wav")

	// Verifica se o arquivo foi enviado ao S3
	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("Failed to list objects in bucket: %v", err)
	}

	// Verifica se o arquivo está no S3 e tem o tamanho correto
	found := false
	for _, object := range output.Contents {
		if *object.Key == fileKey {
			found = true
			// Adiciona logs para o objeto esperado e o encontrado
			t.Logf("Expected file key: %s", fileKey)
			t.Logf("Found file key: %s", *object.Key)
			t.Logf("Expected file size: %d", len(audioContent))
			t.Logf("Found file size: %d", *object.Size)
			t.Logf("Expected last modified time: %v", startTime)
			t.Logf("Found last modified time: %v", object.LastModified)

			// Convertemos o ponteiro para int64 diretamente
			assert.Equal(t, int64(len(audioContent)), *object.Size, "File size should match")
			assert.WithinDuration(t, startTime, object.LastModified.UTC(), time.Second*10, "Upload time should be within expected range")
		}
	}

	if !found {
		t.Errorf("Uploaded file not found in S3 bucket")
	}
}

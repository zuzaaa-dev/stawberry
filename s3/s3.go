package objectstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"marketplace/config"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketBasics struct {
	BucketName string
	S3Client   *s3.Client
}

func ObjectStorageConn(cfg *config.Config) *BucketBasics {
	sdkCfg, err := sdkConfig.LoadDefaultConfig(context.TODO(),
		sdkConfig.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
		sdkConfig.WithRegion(cfg.SigningRegion),
		sdkConfig.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.AccessKey,
				SecretAccessKey: cfg.SecretKEY,
			}, nil
		})),
		sdkConfig.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           cfg.URL,
				SigningRegion: cfg.SigningRegion,
			}, nil
		})),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &BucketBasics{cfg.BucketName, s3.NewFromConfig(sdkCfg)}
}

func (basics BucketBasics) UploadFileWithPresignedURL(ctx context.Context, objectKey string, file io.Reader) error {
	// Создаем PreSigned URL для загрузки файла
	presignClient := s3.NewPresignClient(basics.S3Client)
	presignResult, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(basics.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get presigned URL for upload. Here's why: %v\n", err)
		return err
	}

	// Загружаем файл с использованием PreSigned URL
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Couldn't read file content. Here's why: %v\n", err)
		return err
	}

	req, err := http.NewRequest("PUT", presignResult.URL, bytes.NewReader(fileContent))
	if err != nil {
		log.Printf("Couldn't create HTTP request. Here's why: %v\n", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Couldn't upload file. Here's why: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to upload file. Status code: %v\n", resp.StatusCode)
		return fmt.Errorf("failed to upload file, status code: %v", resp.StatusCode)
	}

	log.Println("File uploaded successfully!")
	return nil
}

func (basics BucketBasics) DownloadFile(ctx context.Context, objectKey string) ([]byte, error) {
	result, err := basics.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(basics.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Printf("Can't get object %s from bucket %s. No such key exists.\n", objectKey, basics.BucketName)
			err = noKey
		} else {
			log.Printf("Couldn't get object %v:%v. Here's why: %v\n", basics.BucketName, objectKey, err)
		}
		return nil, err
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	return body, err
}

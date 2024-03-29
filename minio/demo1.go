package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	bucket := aws.String("test1")
	key := aws.String("testobject")

	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("Ix9YPnPXZcj5yEFI", "YnazF7FvdXwFdKqDxRLwOsdQdHREKXOh", ""),
		Endpoint:         aws.String("http://172.22.50.25:31087"),
		Region:           aws.String("us-east-1a"),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}

	//newSession := session.New(s3Config)
	newSession, _ := session.NewSession(s3Config)

	s3Client := s3.New(newSession)

	cparams := &s3.CreateBucketInput{
		Bucket: bucket, // Required
	}

	// Create a new bucket using the CreateBucket call.
	_, err := s3Client.CreateBucket(cparams)
	if err != nil {
		// Message from an error.
		fmt.Println("创建桶报错:", err.Error())
		return
	}

	// Upload a new object "testobject" with the string "Hello World!" to our "newbucket".
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader("Hello from MinIO!!"),
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n", *bucket, *key, err.Error())
		return
	}
	fmt.Printf("Successfully created bucket %s and uploaded data with key %s\n", *bucket, *key)

	// Retrieve our "testobject" from our "newbucket" and store it locally in "testobject_local".
	file, err := os.Create("testobject_local")
	if err != nil {
		fmt.Println("Failed to create file", err)
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(newSession)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	if err != nil {
		fmt.Println("Failed to download file", err)
		return
	}
	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")
}

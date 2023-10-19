package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func main() {
	ctx := context.Background()
	endpoint := "172.22.50.25:31088"
	accessKeyID := "KKnTKEF2PA7qOVse"
	secretAccessKey := "uGQSnAdQ6xfu1GMqyBg5jS4cOS86WlNs"
	bucketName := "neolink"
	useSSL := false // 设置为 true 如果需要使用 SSL/TLS

	// 创建一个 MinIO 客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 检查指定的存储桶是否存在

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		log.Fatalf("Bucket %s does not exist", bucketName)
	}

	// 上传文件到存储桶
	objectName := "go.mod2"
	filePath := "/data/work/go/GoDevEach/minio/go.mod"
	_, err = client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("File uploaded successfully")

	objectsCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    "",   // 可选：指定文件名前缀
		Recursive: true, // 是否递归查找子目录
	})

	for object := range objectsCh {
		if object.Err != nil {
			log.Println(object.Err)
			return
		}
		fmt.Println(object.Key) // 打印文件名
	}

}

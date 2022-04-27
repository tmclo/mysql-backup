package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Backup service running...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	for range time.Tick(time.Minute * 30) {
		go func() {
			fmt.Println("Running backup...")
			runBackup()
		}()
	}

}

func runBackup() {
	currentTime := time.Now()
	ct := currentTime.Format("01-02-2006-15-04")
	fileName := "sql-backup-" + ct + ".sql"
	//ctx := context.Background()
	host := os.Getenv("MYSQL_DATABASE_HOST")
	password := os.Getenv("MYSQL_PASSWORD")

	str := "/usr/bin/mysqldump -h " + host + " -u root -p" + password + " --all-databases > /tmp/" + fileName
	cmd, err := exec.Command("bash", "-c", str).Output()

	if err != nil {
		fmt.Println(err)
	} else {
		copyFile("/tmp/" + fileName)
		fmt.Println("Backup completed...")
		time.Sleep(time.Second * 10)
		cmd := exec.Command("rm", "-rf", "/tmp/sql-backup*")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(string(cmd))
}

func copyFile(src string) {
	bucket := aws.String(os.Getenv("BUCKET"))
	key := aws.String(src)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("KEY_ID"), os.Getenv("SECRET_KEY"), ""),
		Endpoint:         aws.String("https://s3.us-west-000.backblazeb2.com"),
		Region:           aws.String("us-west-000"),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader("S3 Compatible API"),
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to upload object %s/%s, %s\n", *bucket, *key, err.Error())
	} else {
		fmt.Printf("Successfully uploaded key %s\n", *key)
	}
}

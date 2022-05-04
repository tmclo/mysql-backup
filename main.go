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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorCyan = "\033[36m"

func main() {

	cTime := time.Now()
	curTime := cTime.Format("01-02-2006-15-04")

	fmt.Println(string(colorCyan), curTime+" Backup service running...", string(colorReset))

	err := godotenv.Load()
	if err != nil {
		log.Fatal(string(colorRed), curTime+" Error loading .env file", string(colorReset))
	}

	for range time.Tick(time.Minute * 30) {
		go func() {
			fmt.Println(string(colorGreen), curTime+" Running backup...", string(colorReset))
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
		fmt.Println(string(colorGreen), ct+" Backup completed...", string(colorReset))
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

	cTime := time.Now()
	curTime := cTime.Format("01-02-2006-15-04")

	bucket := aws.String(os.Getenv("BUCKET"))
	key := aws.String(src)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("KEY_ID"), os.Getenv("SECRET_KEY"), ""),
		Endpoint:         aws.String("https://s3.us-west-000.backblazeb2.com"),
		Region:           aws.String("us-west-000"),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    key,
		Body:   strings.NewReader("S3 Compatible API"),
	})
	if err != nil {
		fmt.Printf(curTime+" Failed to upload object %s/%s, %s\n", *bucket, *key, err.Error())
	} else {
		fmt.Printf(curTime+" Successfully uploaded backup %s\n", *key)
	}
}

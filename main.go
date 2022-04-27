package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
	"github.com/kurin/blazer/b2"
)

func main() {
	fmt.Println("Backup service running...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	// b2_authorize_account
	b2, err := b2.NewClient(ctx, os.Getenv("KEY_ID"), os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatalln(err)
	}

	b2.ListBuckets(ctx)

	for range time.Tick(time.Hour * 6) {
		go func() {
			fmt.Println("Running backup...")
			runBackup()
		}()
	}

}

func runBackup() {
	ctx := context.Background()
	host := os.Getenv("MYSQL_DATABASE_HOST")
	password := os.Getenv("MYSQL_PASSWORD")

	cmd := exec.Command("mysqldump -h " + host + " -u root -p" + password + " --all-databases > /tmp/sql-backup.sql")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	} else {
		copyFile(ctx, &b2.Bucket{}, "/tmp/sql-backup.sql", "./mysql-backup-docker/")
		fmt.Println("Backup completed...")
		cmd := exec.Command("rm -rf /tmp/sql-backup.sql")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func copyFile(ctx context.Context, bucket *b2.Bucket, src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	obj := bucket.Object(dst)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, f); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

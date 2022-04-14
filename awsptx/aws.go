package awsptx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	AWSAccessKeyID string = "AWS_ACCESS_KEY_ID"
	AWSSecretKey   string = "AWS_SECRET_KEY"
	AWSRegion      string = "AWS_REGION"
	AWSBucketName  string = "AWS_BUCKET_NAME"
)

var (
	AccessKeyID     string
	SecretAccessKey string
	region          string
	awsAccessKeyID  = os.Getenv(AWSAccessKeyID)
	awsSecretKey    = os.Getenv(AWSSecretKey)
	awsRegion       = os.Getenv(AWSRegion)
	awsBucketName   = os.Getenv(AWSBucketName)
)

var instance *AWSS3

type AWSS3 struct {
	Sess *session.Session
}

func New() *AWSS3 {
	if instance == nil {
		instance = &AWSS3{
			Sess: connectAws(),
		}
	}
	return instance
}

func connectAws() *session.Session {
	if awsAccessKeyID == "" {
		panic("awsS3 access key must be set")
	}
	if awsSecretKey == "" {
		panic("awsS3 secret key must be set")
	}
	if awsRegion == "" {
		panic("awsS3 region must be set")
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(
				awsAccessKeyID,
				awsSecretKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func getS3() *s3.S3 {
	if instance == nil {
		New()
	}
	return s3.New(instance.Sess)
}

func (awsS3 *AWSS3) GetManagerUploader() *s3manager.Uploader {
	if awsS3.Sess == nil {
		New()
	}
	return s3manager.NewUploader(awsS3.Sess)
}

func (awsS3 *AWSS3) ListBuckets() {
	result, err := getS3().ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func (aswS3 *AWSS3) ListBucketItems() {
	if awsBucketName == "" {
		panic("bucket name must be set")
	}
	svc := getS3()
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(awsBucketName)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", awsBucketName, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}
}

func (awsS3 *AWSS3) UploadFile(filename string, file io.Reader) bool {
	if awsBucketName == "" {
		panic("bucket name must be set")
	}
	uploader := awsS3.GetManagerUploader()
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, awsBucketName, err)
		return false
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, awsBucketName)
	return true
}

func (awsS3 *AWSS3) DownloadFile(item string) *os.File {
	New()
	downloader := s3manager.NewDownloader(awsS3.Sess)
	file, err := os.Create(item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", item, err)
		return nil
	}
	defer file.Close()
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(awsBucketName),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
		return nil
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return file
}

func GetFile(key string) bytes.Buffer {
	svc := getS3()
	params := &s3.GetObjectInput{
		Bucket: aws.String(awsBucketName), // Required
		Key:    aws.String(key),           // Required
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Fatal(err.Error())
	}

	size := int(*resp.ContentLength)

	buffer := make([]byte, size)
	defer resp.Body.Close()
	var bbuffer bytes.Buffer
	for true {
		num, rerr := resp.Body.Read(buffer)
		if num > 0 {
			bbuffer.Write(buffer[:num])
		} else if rerr == io.EOF || rerr != nil {
			break
		}
	}
	return bbuffer
}

func (awsS3 *AWSS3) DeleteFile(filename string) bool {
	svc := getS3()
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(awsBucketName), Key: aws.String(filename)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", filename, awsBucketName, err)
		return false
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return false
	}
	return true
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

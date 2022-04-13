package awsptx

import (
	"fmt"
	"io"
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

var instance *session.Session

type AWSS3 struct {
	Sess *session.Session
}

func GetInstance() *AWSS3 {
	awsS3 := &AWSS3{}
	if instance == nil {
		awsS3.Sess = ConnectAws()
	}
	return awsS3
}

func ConnectAws() *session.Session {
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

func (awsS3 *AWSS3) GetClient() *s3.S3 {
	if awsS3.Sess == nil {
		GetInstance()
	}
	return s3.New(awsS3.Sess)
}

func (awsS3 *AWSS3) GetManagerUploader() *s3manager.Uploader {
	if awsS3.Sess == nil {
		GetInstance()
	}
	return s3manager.NewUploader(awsS3.Sess)
}

func (awsS3 *AWSS3) ListBuckets() {
	result, err := awsS3.GetClient().ListBuckets(nil)
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
	svc := GetInstance().GetClient()
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

func (awsS3 *AWSS3) UploadFile(filename string, file io.Reader) {
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
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, awsBucketName)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

package awsptx

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWSAccessKeyID string = "AWS_ACCESS_KEY_ID"
	AWSSecretKey   string = "AWS_SECRET_KEY"
	AWSRegion      string = "AWS_REGION"
)

var AccessKeyID string
var SecretAccessKey string
var region string

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
	AccessKeyID = os.Getenv(AWSAccessKeyID)
	if AccessKeyID == "" {
		panic("awsS3 access key must be set")
	}
	SecretAccessKey = os.Getenv(AWSSecretKey)
	if SecretAccessKey == "" {
		panic("awsS3 secret key must be set")
	}
	region = os.Getenv(AWSRegion)
	if region == "" {
		panic("awsS3 region must be set")
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

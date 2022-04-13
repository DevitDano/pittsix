package awsptx_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/devitdano/pittsix/awsptx"
	"github.com/stretchr/testify/assert"
)

func TestListBucket(t *testing.T) {
	ass := assert.New(t)
	awsptx.GetInstance().ListBuckets()
	ass.True(true)
}

func TestListBucketItems(t *testing.T) {
	ass := assert.New(t)
	awsptx.GetInstance().ListBucketItems()
	ass.True(true)
}

func TestListUploadFile(t *testing.T) {
	ass := assert.New(t)
	file, err := os.Open("../dano.jpg")
	if err != nil {
		panic(err)
	}
	awsptx.GetInstance().UploadFile("testfile", bufio.NewReader(file))
	ass.True(true)
}

func TestListDownloadFile(t *testing.T) {
	ass := assert.New(t)
	awsptx.GetInstance().DownloadFile("testfile")
	ass.True(true)
}

func TestListDeleteFile(t *testing.T) {
	ass := assert.New(t)
	awsptx.GetInstance().DeleteFile("testfile")
	ass.True(true)
}

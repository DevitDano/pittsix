package awsptx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awsptx "github.com/devitdano/pittsix/awsptx"
)

func TestListBucket(t *testing.T) {
	ass := assert.New(t)
	awsptx.GetInstance().ListBuckets()
	ass.True(true)
}

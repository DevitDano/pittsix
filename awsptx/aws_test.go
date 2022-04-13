package awsptx_test

import (
	"testing"

	"github.com/devitdano/pittsix/awsptx"
	"github.com/stretchr/testify/assert"
)

func TestListBucket(t *testing.T) {
	ass := assert.New(t)
	awsptx.GetInstance().ListBuckets()
	ass.True(true)
}

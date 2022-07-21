package filesystem

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"strings"
	"testing"
)

type mockedS3Api struct {
	getObjectOutput   *s3.GetObjectOutput
	listObjectsOutput *s3.ListObjectsV2Output
}

func (s3Api mockedS3Api) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return s3Api.getObjectOutput, nil
}

func (s3Api mockedS3Api) ListObjects(ctx context.Context, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return s3Api.listObjectsOutput, nil
}

func stringToReadCloser(s string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(s))
}

func TestS3Filesystem_GetFile(t *testing.T) {
	mockedS3Api := mockedS3Api{
		getObjectOutput: &s3.GetObjectOutput{
			Body: stringToReadCloser("body"),
		},
		listObjectsOutput: nil,
	}
	fs := newS3FilesystemFromApi("bucket", "no-region", mockedS3Api)
	result, _ := fs.GetFile("hasbody.txt")
	if result.GetData()["content"].(string) != "body" {
		t.Fail()
	}
}

func TestS3Filesystem_GetFolder(t *testing.T) {
	key1 := "key-1"
	key2 := "key-2"
	mockedS3Api := mockedS3Api{
		getObjectOutput: nil,
		listObjectsOutput: &s3.ListObjectsV2Output{
			Contents: []types.Object{
				types.Object{Key: &key1, Size: 101},
				types.Object{Key: &key2, Size: 5985},
			},
		},
	}
	fs := newS3FilesystemFromApi("bucket", "no-region", mockedS3Api)
	result, _ := fs.GetFolder("my-folder/")
	children := result.GetData()["children"].([]ChildResult)
	if len(children) != 2 {
		t.Fail()
	}
	if children[0].Key != key1 && children[1].Key != key2 {
		t.Fail()
	}
}

package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"strings"
)

type S3Filesystem struct {
	bucketName   string
	bucketRegion string
	config       aws.Config
	context      context.Context
}

func NewS3Filesystem(name string, region string) S3Filesystem {
	// Loads AWS credentials from environment
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}
	return S3Filesystem{
		bucketName:   name,
		bucketRegion: region,
		config:       cfg,
		context:      context.TODO(),
	}
}

func (fs S3Filesystem) GetFile(path string) (Result, error) {
	s3Client := s3.NewFromConfig(fs.config)

	// s3 doesn't have a "root" path like a filesystem would, so we trim the leading slash
	path = strings.TrimLeft(path, "/")

	input := s3.GetObjectInput{
		Bucket: &fs.bucketName,
		Key:    &path,
	}
	result, err := s3Client.GetObject(fs.context, &input)
	if err != nil {
		var nsk *types.NoSuchKey
		// NoSuchKey is an expected error
		if errors.As(err, &nsk) {
			return nil, fmt.Errorf("not found")
		} else { // Other error codes are unexpected errors
			// TODO: pattern match(?) on expected (not found) and unexpected errors
			return nil, err
		}
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, err
	}
	content := buf.String()
	return FileResult{Path: path, Content: content}, nil
}

func (fs S3Filesystem) GetFolder(path string) (Result, error) {
	s3Client := s3.NewFromConfig(fs.config)

	// s3 doesn't have true "folders", just path prefix, so a leading slash is equivalent to empty string
	// TODO: since both s3 file and folder query trim leading slash from URI, is it worth adding logic
	// 	to the Filesystem interface, i.e. NormalizePath(path string), that will be called in filesystem.GetPath()
	path = strings.TrimLeft(path, "/")

	prefixDelimiter := folderSuffixCharacter
	listInput := s3.ListObjectsV2Input{
		Bucket: &fs.bucketName,
		Prefix: &path,
		Delimiter: &prefixDelimiter,
	}
	listObjectsResult, err := s3Client.ListObjectsV2(fs.context, &listInput)
	if err != nil {
		fmt.Println("Error listing objects with prefix", path, err)
		return nil, err
	}
	// check if there are any child objects, if not, return empty list
	if listObjectsResult.Contents != nil {
		var children []ChildResult
		for _, obj := range listObjectsResult.Contents {
			relativeKey := strings.TrimLeft(*obj.Key, path)
			// the "folder object" is returned as a "child", so we skip it
			if relativeKey != "" {
				children = append(children, ChildResult{
					Key:  relativeKey,
					Size: obj.Size,
				})
			}
		}
		return FolderResult{Path: path, Children: children}, nil
	} else {
		return FolderResult{Path: path, Children: make([]ChildResult, 0)}, nil
	}

}

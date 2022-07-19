package filesystem

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Filesystem struct {
	bucketName   string
	bucketRegion string
	config       aws.Config
	context      context.Context
}

func NewS3Filesystem(name string, region string) S3Filesystem {
	// Loads AWS credentials from environment
	cfg, err := config.LoadDefaultConfig(context.TODO())
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

func (fs S3Filesystem) GetPath(path string) (Result, error) {
	s3Client := s3.NewFromConfig(fs.config)
	// First try to get the key
	result, err := fs.getFile(s3Client, path)
	if err != nil {
		// If there is no object at the key, treat it like a "folder"
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			// TODO: should return "error: not found" when no files at this folder path prefix?
			return fs.getFolder(s3Client, &fs.bucketName, &path)
		} else { // Other error codes are fatal errors
			return nil, err
		}
	}
	// Object was found at the key
	return result, nil
}

func (fs S3Filesystem) getFile(s3Client *s3.Client, path string) (Result, error) {
	input := s3.GetObjectInput{
		Bucket: &fs.bucketName,
		Key:    &path,
	}
	result, err := s3Client.GetObject(fs.context, &input)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, err
	}
	content := buf.String()
	return FileResult{Path: path, Content: content}, nil
}

func (fs S3Filesystem) getFolder(s3Client *s3.Client, bucketName *string, path *string) (Result, error) {
	listInput := s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: path,
	}
	listObjectsResult, err := s3Client.ListObjectsV2(fs.context, &listInput)
	if err != nil {
		return nil, err
	}
	var children []ChildResult
	for _, obj := range listObjectsResult.Contents {
		children = append(children, ChildResult{
			Key:  *obj.Key,
			Size: obj.Size,
		})
	}
	return FolderResult{Path: *path, Children: children}, nil
}

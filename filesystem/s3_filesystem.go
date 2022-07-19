package filesystem

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Filesystem struct {
	bucketName   string
	bucketRegion string
	session      *session.Session
}

func NewS3Filesystem(name string, region string) S3Filesystem {
	// Loads AWS credentials from environment
	awsSession := session.Must(session.NewSession())
	return S3Filesystem{
		bucketName:   name,
		bucketRegion: region,
		session:      awsSession,
	}
}

func (fs S3Filesystem) GetPath(path string) (Result, error) {
	s3Client := s3.New(fs.session)
	// First try to get the key
	result, err := fs.getFile(s3Client, path)
	if err != nil {
		// If there is no object at the key, treat it like a "folder"
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			// TODO: return "error: not found" when no files at this folder path?
			return getFolder(s3Client, &fs.bucketName, &path)
		} else { // Other error codes are fatal errors
			return nil, fmt.Errorf(aerr.Message())
		}
	}
	// Object was found at the key
	return result, nil
}

func (fs S3Filesystem) getFile(s3Client *s3.S3, path string) (Result, error) {
	input := s3.GetObjectInput{
		Bucket: &fs.bucketName,
		Key:    &path,
	}
	result, err := s3Client.GetObject(&input)
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

func getFolder(s3Client *s3.S3, bucketName *string, path *string) (Result, error) {
	listInput := s3.ListObjectsV2Input{
		Bucket: bucketName,
		Prefix: path,
	}
	listObjectsResult, err := s3Client.ListObjectsV2(&listInput)
	if err != nil {
		return nil, err
	}
	var children []ChildResult
	for _, obj := range listObjectsResult.Contents {
		children = append(children, ChildResult{
			Key:  *obj.Key,
			Size: *obj.Size,
		})
	}
	return FolderResult{Path: *path, Children: children}, nil
}

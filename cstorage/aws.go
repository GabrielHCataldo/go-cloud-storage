package cstorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type awsS3Client struct {
	*s3.S3
}

type iAwsS3Client interface {
	CreateBucket(ctx context.Context, input CreateBucketInput) error
	PutObject(ctx context.Context, input PutObjectInput) error
	GetObjectByKey(ctx context.Context, bucketName, key string) (*Object, error)
	GetObjectUrl(bucketName, key string) string
	ListObjects(ctx context.Context, bucketName string, opts ...OptsListObjects) ([]ObjectSummary, error)
	DeleteObject(ctx context.Context, input DeleteObjectInput) error
	DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error
	DeleteBucket(ctx context.Context, bucketName string) error
}

func newAwsS3StorageClient(cfgs ...*aws.Config) (iAwsS3Client, error) {
	nSession, err := session.NewSession(cfgs...)
	return awsS3Client{
		S3: s3.New(nSession),
	}, err
}

func (a awsS3Client) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	_, err := a.S3.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(input.Bucket),
	})
	return err
}

func (a awsS3Client) PutObject(ctx context.Context, input PutObjectInput) error {
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNotNil(err) {
		return err
	}
	_, err = a.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(bytesContent),
		Bucket: aws.String(input.Bucket),
		//ContentLength: aws.Int64(int64(len(bytesContent))),
		ContentType: aws.String(input.MimeType.String()),
		Key:         aws.String(input.Key),
	})
	return err
}

func (a awsS3Client) GetObjectByKey(ctx context.Context, bucketName, key string) (*Object, error) {
	obj, err := a.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if helper.IsNotNil(err) {
		return nil, err
	}
	objResult := parseAwsS3StorageObject(obj)
	objResult.Key = key
	objResult.Url = a.GetObjectUrl(bucketName, key)
	return &objResult, nil
}

func (a awsS3Client) GetObjectUrl(bucketName, key string) string {
	url := "https://%s.amazonaws.com/%s/%s"
	return fmt.Sprintf(url, a.S3.ResolvedRegion, bucketName, key)
}

func (a awsS3Client) ListObjects(ctx context.Context, bucketName string, opts ...OptsListObjects) ([]ObjectSummary, error) {
	opt := GetOptListObjectsByParams(opts)
	objs, err := a.S3.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket:    aws.String(bucketName),
		Delimiter: aws.String(opt.Delimiter),
		Prefix:    aws.String(opt.Prefix),
	})
	if helper.IsNotNil(err) {
		return nil, err
	}
	var result []ObjectSummary
	for _, obj := range objs.Contents {
		objResult := parseAwsS3StorageObjectSummary(obj)
		objResult.Url = a.GetObjectUrl(bucketName, objResult.Key)
		result = append(result, objResult)
	}
	return result, nil
}

func (a awsS3Client) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	_, err := a.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	return err
}

func (a awsS3Client) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	objs, err := a.S3.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(input.Bucket),
		Prefix: aws.String(input.Prefix),
	})
	if helper.IsNotNil(err) {
		return err
	}
	for _, obj := range objs.Contents {
		err = a.DeleteObject(ctx, DeleteObjectInput{
			Bucket: input.Bucket,
			Key:    helper.ConvertPointerToValue(obj.Key),
		})
		if helper.IsNotNil(err) {
			return err
		}
	}
	return nil
}

func (a awsS3Client) DeleteBucket(ctx context.Context, bucketName string) error {
	_, err := a.S3.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

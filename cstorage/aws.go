package cstorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type awsS3Client struct {
	config aws.Config
	*s3.Client
}

type iAwsS3Client interface {
	CreateBucket(ctx context.Context, input CreateBucketInput) error
	PutObject(ctx context.Context, input PutObjectInput) error
	GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error)
	GetObjectUrl(bucket, key string) string
	ListObjects(ctx context.Context, bucket string, opts ...OptsListObjects) ([]ObjectSummary, error)
	DeleteObject(ctx context.Context, input DeleteObjectInput) error
	DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error
	DeleteBucket(ctx context.Context, bucket string) error
}

func newAwsS3StorageClient(cfg aws.Config) iAwsS3Client {
	return awsS3Client{
		Client: s3.NewFromConfig(cfg),
		config: cfg,
	}
}

func (a awsS3Client) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	region := a.config.Region
	if helper.IsNotEmpty(input.Location) {
		region = input.Location
	}
	_, err := a.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(input.Bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	return err
}

func (a awsS3Client) PutObject(ctx context.Context, input PutObjectInput) error {
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNotNil(err) {
		return err
	}
	_, err = a.Client.PutObject(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(bytesContent),
		Bucket: aws.String(input.Bucket),
		//ContentLength: aws.Int64(int64(len(bytesContent))),
		ContentType: aws.String(input.MimeType.String()),
		Key:         aws.String(input.Key),
	})
	return err
}

func (a awsS3Client) GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error) {
	obj, err := a.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if helper.IsNotNil(err) {
		return nil, err
	}
	objResult := parseAwsS3StorageObject(obj)
	objResult.Key = key
	objResult.Url = a.GetObjectUrl(bucket, key)
	return &objResult, nil
}

func (a awsS3Client) GetObjectUrl(bucket, key string) string {
	url := "https://%s.amazonaws.com/%s/%s"
	return fmt.Sprintf(url, a.config.Region, bucket, key)
}

func (a awsS3Client) ListObjects(ctx context.Context, bucket string, opts ...OptsListObjects) ([]ObjectSummary, error) {
	opt := GetOptListObjectsByParams(opts)
	objs, err := a.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String(opt.Delimiter),
		Prefix:    aws.String(opt.Prefix),
	})
	if helper.IsNotNil(err) {
		return nil, err
	}
	var result []ObjectSummary
	for _, obj := range objs.Contents {
		objResult := parseAwsS3StorageObjectSummary(obj)
		objResult.Url = a.GetObjectUrl(bucket, objResult.Key)
		result = append(result, objResult)
	}
	return result, nil
}

func (a awsS3Client) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	_, err := a.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	return err
}

func (a awsS3Client) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	objs, err := a.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
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

func (a awsS3Client) DeleteBucket(ctx context.Context, bucket string) error {
	_, err := a.Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}

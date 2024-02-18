package cstorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
)

type awsS3Client struct {
	config aws.Config
	client *s3.Client
}

// NewAwsS3Storage new instance of connection with AWS S3 storage, to close it just use Disconnect() or SimpleDisconnect()
func NewAwsS3Storage(cfg aws.Config) CStorage {
	return &awsS3Client{
		client: s3.NewFromConfig(cfg),
		config: cfg,
	}
}

func (a *awsS3Client) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	region := a.config.Region
	if helper.IsNotEmpty(input.Location) {
		region = input.Location
	}
	_, err := a.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(input.Bucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	return err
}

func (a *awsS3Client) PutObject(ctx context.Context, input PutObjectInput) error {
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNil(err) {
		_, err = a.client.PutObject(ctx, &s3.PutObjectInput{
			Body:        bytes.NewReader(bytesContent),
			Bucket:      aws.String(input.Bucket),
			ContentType: aws.String(input.MimeType.String()),
			Key:         aws.String(input.Key),
		})
	}
	return err
}

func (a *awsS3Client) PutObjects(ctx context.Context, inputs ...PutObjectInput) []PutObjectOutput {
	var result []PutObjectOutput
	for _, input := range inputs {
		err := a.PutObject(ctx, input)
		result = append(result, PutObjectOutput{
			Bucket: input.Bucket,
			Key:    input.Key,
			Err:    err,
		})
	}
	return result
}

func (a *awsS3Client) GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error) {
	obj, err := a.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if helper.IsNotNil(err) {
		return nil, err
	}
	bs, _ := io.ReadAll(obj.Body)
	objResult := parseAwsS3StorageObject(obj)
	objResult.Key = key
	objResult.Url = a.GetObjectUrl(bucket, key)
	objResult.Content = bs
	return &objResult, nil
}

func (a *awsS3Client) GetObjectUrl(bucket, key string) string {
	url := "https://%s.amazonaws.com/%s/%s"
	return fmt.Sprintf(url, a.config.Region, bucket, key)
}

func (a *awsS3Client) ListObjects(ctx context.Context, bucket string, opts ...*OptsListObjects) ([]ObjectSummary, error) {
	opt := MergeOptsListObjectsByParams(opts)
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String(opt.Delimiter),
		Prefix:    aws.String(opt.Prefix),
	}
	objs, err := a.client.ListObjectsV2(ctx, input)
	var result []ObjectSummary
	if helper.IsNotNil(objs) {
		for _, obj := range objs.Contents {
			objResult := parseAwsS3StorageObjectSummary(obj)
			objResult.Url = a.GetObjectUrl(bucket, objResult.Key)
			result = append(result, objResult)
		}
	}
	return result, err
}

func (a *awsS3Client) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	return err
}

func (a *awsS3Client) DeleteObjects(ctx context.Context, inputs ...DeleteObjectInput) []DeleteObjectsOutput {
	var result []DeleteObjectsOutput
	for _, input := range inputs {
		err := a.DeleteObject(ctx, input)
		result = append(result, DeleteObjectsOutput{
			Bucket: input.Bucket,
			Key:    input.Key,
			Err:    err,
		})
	}
	return result
}

func (a *awsS3Client) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	objs, err := a.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(input.Bucket),
		Prefix: aws.String(input.Prefix),
	})
	if helper.IsNotNil(objs) {
		for _, obj := range objs.Contents {
			err = a.DeleteObject(ctx, DeleteObjectInput{
				Bucket: input.Bucket,
				Key:    helper.ConvertPointerToValue(obj.Key),
			})
		}
	}
	return err
}

func (a *awsS3Client) DeleteObjectsByPrefixes(ctx context.Context, inputs ...DeletePrefixInput) []DeletePrefixOutput {
	var result []DeletePrefixOutput
	for _, input := range inputs {
		err := a.DeleteObjectsByPrefix(ctx, input)
		result = append(result, DeletePrefixOutput{
			Bucket: input.Bucket,
			Prefix: input.Prefix,
			Err:    err,
		})
	}
	return result
}

func (a *awsS3Client) DeleteBucket(ctx context.Context, bucket string) error {
	_, err := a.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}

func (a *awsS3Client) DeleteBuckets(ctx context.Context, buckets ...string) []DeleteBucketsOutput {
	var result []DeleteBucketsOutput
	for _, bucket := range buckets {
		err := a.DeleteBucket(ctx, bucket)
		result = append(result, DeleteBucketsOutput{
			Bucket: bucket,
			Err:    err,
		})
	}
	return result
}

func (a *awsS3Client) Disconnect() error {
	return nil
}

func (a *awsS3Client) SimpleDisconnect() {
}

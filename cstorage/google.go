package cstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
)

type googleStorageClient struct {
	client *storage.Client
}

// NewGoogleStorage new instance of connection with Google storage, to close it just use CStorage.Disconnect() or CStorage.SimpleDisconnect()
func NewGoogleStorage(ctx context.Context, opts ...option.ClientOption) (i CStorage, err error) {
	client, err := storage.NewClient(ctx, opts...)
	if helper.IsNil(err) {
		i = &googleStorageClient{
			client: client,
		}
	}
	return i, err
}

func (g googleStorageClient) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	return g.client.Bucket(input.Bucket).Create(ctx, input.ProjectId, &storage.BucketAttrs{Location: input.Location})
}

func (g googleStorageClient) PutObject(ctx context.Context, input PutObjectInput) error {
	obj := g.client.Bucket(input.Bucket).Object(input.Key)
	fw := obj.NewWriter(ctx)
	fw.ContentType = input.MimeType.String()
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNil(err) {
		_, err = fw.Write(bytesContent)
	}
	_ = fw.Close()
	return err
}

func (g googleStorageClient) PutObjects(ctx context.Context, inputs ...PutObjectInput) []PutObjectOutput {
	var result []PutObjectOutput
	for _, input := range inputs {
		err := g.PutObject(ctx, input)
		result = append(result, PutObjectOutput{
			Bucket: input.Bucket,
			Key:    input.Key,
			Err:    err,
		})
	}
	return result
}

func (g googleStorageClient) GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error) {
	obj := g.client.Bucket(bucket).Object(key)
	attrs, err := obj.Attrs(ctx)
	if helper.IsNotNil(err) {
		return nil, err
	}
	var bs []byte
	reader, _ := obj.NewReader(ctx)
	bs, err = io.ReadAll(reader)
	objResult := parseGoogleStorageObject(attrs)
	objResult.Url = g.GetObjectUrl(bucket, key)
	objResult.Content = bs
	return &objResult, err
}

func (g googleStorageClient) GetObjectUrl(bucket, key string) string {
	url := "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(url, bucket, key)
}

func (g googleStorageClient) ListObjects(ctx context.Context, bucket string, opts ...*OptsListObjects) ([]ObjectSummary,
	error) {
	opt := MergeOptsListObjectsByParams(opts)
	bkt := g.client.Bucket(bucket)
	objs := bkt.Objects(ctx, &storage.Query{
		Delimiter: opt.Delimiter,
		Prefix:    opt.Prefix,
	})
	var result []ObjectSummary
	var rErr error
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(err) {
			rErr = err
			break
		} else {
			objResult := parseGoogleStorageObjectSummary(obj)
			objResult.Url = g.GetObjectUrl(bucket, obj.Name)
			result = append(result, objResult)
		}
	}
	return result, rErr
}

func (g googleStorageClient) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	bkt := g.client.Bucket(input.Bucket)
	return bkt.Object(input.Key).Delete(ctx)
}

func (g googleStorageClient) DeleteObjects(ctx context.Context, inputs ...DeleteObjectInput) []DeleteObjectsOutput {
	var result []DeleteObjectsOutput
	for _, input := range inputs {
		err := g.DeleteObject(ctx, input)
		result = append(result, DeleteObjectsOutput{
			Bucket: input.Bucket,
			Key:    input.Key,
			Err:    err,
		})
	}
	return result
}

func (g googleStorageClient) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	bkt := g.client.Bucket(input.Bucket)
	objs := bkt.Objects(ctx, &storage.Query{Prefix: input.Prefix})
	var rErr error
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(err) {
			rErr = err
			break
		} else {
			rErr = bkt.Object(obj.Name).Delete(ctx)
		}
	}
	return rErr
}

func (g googleStorageClient) DeleteObjectsByPrefixes(ctx context.Context, inputs ...DeletePrefixInput) []DeletePrefixOutput {
	var result []DeletePrefixOutput
	for _, input := range inputs {
		err := g.DeleteObjectsByPrefix(ctx, input)
		result = append(result, DeletePrefixOutput{
			Bucket: input.Bucket,
			Prefix: input.Prefix,
			Err:    err,
		})
	}
	return result
}

func (g googleStorageClient) DeleteBucket(ctx context.Context, bucket string) error {
	return g.client.Bucket(bucket).Delete(ctx)
}

func (g googleStorageClient) DeleteBuckets(ctx context.Context, buckets ...string) []DeleteBucketsOutput {
	var result []DeleteBucketsOutput
	for _, bucket := range buckets {
		err := g.DeleteBucket(ctx, bucket)
		result = append(result, DeleteBucketsOutput{
			Bucket: bucket,
			Err:    err,
		})
	}
	return result
}

func (g googleStorageClient) Disconnect() error {
	return g.client.Close()
}

func (g googleStorageClient) SimpleDisconnect() {
	_ = g.client.Close()
	logger.InfoSkipCaller(3, "Connection to google storage closed.")
}

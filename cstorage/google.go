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
	*storage.Client
}

type iGoogleStorageClient interface {
	CreateBucket(ctx context.Context, input CreateBucketInput) error
	PutObject(ctx context.Context, input PutObjectInput) error
	GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error)
	GetObjectUrl(bucket, key string) string
	ListObjects(ctx context.Context, bucket string, opts ...*OptsListObjects) ([]ObjectSummary, error)
	DeleteObject(ctx context.Context, input DeleteObjectInput) error
	DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error
	DeleteBucket(ctx context.Context, bucket string) error
	Disconnect() error
	SimpleDisconnect()
}

func newGoogleStorageClient(ctx context.Context, opts ...option.ClientOption) (iGoogleStorageClient, error) {
	client, err := storage.NewClient(ctx, opts...)
	return googleStorageClient{
		Client: client,
	}, errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	err := g.Client.Bucket(input.Bucket).Create(ctx, input.ProjectId, &storage.BucketAttrs{Location: input.Location})
	return errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) PutObject(ctx context.Context, input PutObjectInput) error {
	obj := g.Client.Bucket(input.Bucket).Object(input.Key)
	fw := obj.NewWriter(ctx)
	fw.ContentType = input.MimeType.String()
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNil(err) {
		_, err = fw.Write(bytesContent)
	}
	_ = fw.Close()
	return errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error) {
	obj := g.Client.Bucket(bucket).Object(key)
	attrs, err := obj.Attrs(ctx)
	if helper.IsNotNil(err) {
		return nil, errors.NewSkipCaller(3, err)
	}
	var bs []byte
	reader, _ := obj.NewReader(ctx)
	bs, err = io.ReadAll(reader)
	objResult := parseGoogleStorageObject(attrs)
	objResult.Url = g.GetObjectUrl(bucket, key)
	objResult.Content = bs
	return &objResult, errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) GetObjectUrl(bucket, key string) string {
	url := "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(url, bucket, key)
}

func (g googleStorageClient) ListObjects(ctx context.Context, bucket string, opts ...*OptsListObjects) ([]ObjectSummary,
	error) {
	opt := MergeOptsListObjectsByParams(opts)
	bkt := g.Client.Bucket(bucket)
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
	return result, errors.NewSkipCaller(3, rErr)
}

func (g googleStorageClient) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	bkt := g.Client.Bucket(input.Bucket)
	err := bkt.Object(input.Key).Delete(ctx)
	return errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	bkt := g.Client.Bucket(input.Bucket)
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
	return errors.NewSkipCaller(3, rErr)
}

func (g googleStorageClient) DeleteBucket(ctx context.Context, bucket string) error {
	err := g.Client.Bucket(bucket).Delete(ctx)
	return errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) Disconnect() error {
	err := g.Client.Close()
	return errors.NewSkipCaller(3, err)
}

func (g googleStorageClient) SimpleDisconnect() {
	_ = g.Client.Close()
	logger.InfoSkipCaller(3, "Connection to google storage closed.")
}

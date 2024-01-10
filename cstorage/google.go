package cstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
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
	ListObjects(ctx context.Context, bucket string, opts ...OptsListObjects) ([]ObjectSummary, error)
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
	}, err
}

func (g googleStorageClient) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	return g.Client.Bucket(input.Bucket).Create(ctx, input.ProjectId, &storage.BucketAttrs{
		Location: input.Location,
	})
}

func (g googleStorageClient) PutObject(ctx context.Context, input PutObjectInput) error {
	obj := g.Client.Bucket(input.Bucket).Object(input.Key)
	fw := obj.NewWriter(ctx)
	fw.ContentType = input.MimeType.String()
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNotNil(err) {
		return err
	}
	_, err = fw.Write(bytesContent)
	_ = fw.Close()
	return err
}

func (g googleStorageClient) GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error) {
	obj := g.Client.Bucket(bucket).Object(key)
	attrs, err := obj.Attrs(ctx)
	if helper.IsNotNil(err) {
		return nil, err
	}
	var bs []byte
	reader, err := obj.NewReader(ctx)
	if helper.IsNotNil(reader) {
		bs, err = io.ReadAll(reader)
	}
	objResult := parseGoogleStorageObject(attrs)
	objResult.Url = g.GetObjectUrl(bucket, key)
	objResult.Content = bs
	return &objResult, err
}

func (g googleStorageClient) GetObjectUrl(bucket, key string) string {
	url := "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(url, bucket, key)
}

func (g googleStorageClient) ListObjects(ctx context.Context, bucket string, opts ...OptsListObjects) (
	[]ObjectSummary, error) {
	opt := GetOptListObjectsByParams(opts)
	bkt := g.Client.Bucket(bucket)
	objs := bkt.Objects(ctx, &storage.Query{
		Delimiter: opt.Delimiter,
		Prefix:    opt.Prefix,
	})
	var result []ObjectSummary
	var err error
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(obj) {
			objResult := parseGoogleStorageObjectSummary(obj)
			objResult.Url = g.GetObjectUrl(bucket, obj.Name)
			result = append(result, objResult)
		}
	}
	return result, err
}

func (g googleStorageClient) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	bkt := g.Client.Bucket(input.Bucket)
	return bkt.Object(input.Key).Delete(ctx)
}

func (g googleStorageClient) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	bkt := g.Client.Bucket(input.Bucket)
	objs := bkt.Objects(ctx, &storage.Query{Prefix: input.Prefix})
	var err error
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(err) {
			return ErrPrefixNotExists
		}
		err = bkt.Object(obj.Name).Delete(ctx)
	}
	return err
}

func (g googleStorageClient) DeleteBucket(ctx context.Context, bucket string) error {
	return g.Client.Bucket(bucket).Delete(ctx)
}

func (g googleStorageClient) Disconnect() error {
	return g.Client.Close()
}

func (g googleStorageClient) SimpleDisconnect() {
	_ = g.Disconnect()
	logger.Info("connection to google storage closed.")
}

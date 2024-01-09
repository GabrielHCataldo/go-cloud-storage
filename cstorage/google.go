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
)

type googleStorageClient struct {
	*storage.Client
}

type iGoogleStorageClient interface {
	CreateBucket(ctx context.Context, input CreateBucketInput) error
	PutObject(ctx context.Context, input PutObjectInput) error
	GetObjectByKey(ctx context.Context, bucketName, key string) (*Object, error)
	GetObjectUrl(bucketName, key string) string
	ListObjects(ctx context.Context, bucketName string, opts ...OptsListObjects) ([]ObjectSummary, error)
	DeleteObject(ctx context.Context, input DeleteObjectInput) error
	DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error
	DeleteBucket(ctx context.Context, bucketName string) error
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
	f := input.Key + input.MimeType.Extension()
	obj := g.Client.Bucket(input.Bucket).Object(f)
	fw := obj.NewWriter(ctx)
	fw.ContentType = input.MimeType.String()
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNotNil(err) {
		return err
	}
	_, err = fw.Write(bytesContent)
	if helper.IsNotNil(err) {
		return err
	}
	return fw.Close()
}

func (g googleStorageClient) GetObjectByKey(ctx context.Context, bucketName, key string) (*Object, error) {
	obj, err := g.Client.Bucket(bucketName).Object(key).Attrs(ctx)
	if helper.IsNotNil(err) {
		return nil, err
	}
	objResult := parseGoogleStorageObject(obj)
	objResult.Url = g.GetObjectUrl(bucketName, key)
	return &objResult, nil
}

func (g googleStorageClient) GetObjectUrl(bucketName, key string) string {
	url := "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(url, bucketName, key)
}

func (g googleStorageClient) ListObjects(ctx context.Context, bucketName string, opts ...OptsListObjects) (
	[]ObjectSummary, error) {
	opt := GetOptListObjectsByParams(opts)
	bkt := g.Client.Bucket(bucketName)
	objs := bkt.Objects(ctx, &storage.Query{
		Delimiter: opt.Delimiter,
		Prefix:    opt.Prefix,
	})
	var result []ObjectSummary
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(err) {
			return nil, err
		}
		objResult := parseGoogleStorageObjectSummary(obj)
		objResult.Url = g.GetObjectUrl(bucketName, obj.Name)
		result = append(result, objResult)
	}
	return result, nil
}

func (g googleStorageClient) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	bkt := g.Client.Bucket(input.Bucket)
	return bkt.Object(input.Key).Delete(ctx)
}

func (g googleStorageClient) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	bkt := g.Client.Bucket(input.Bucket)
	objs := bkt.Objects(ctx, &storage.Query{Prefix: input.Prefix})
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(err) {
			return err
		}
		err = bkt.Object(obj.Name).Delete(ctx)
		if helper.IsNotNil(err) {
			return err
		}
	}
	return nil
}

func (g googleStorageClient) DeleteBucket(ctx context.Context, bucketName string) error {
	return g.Client.Bucket(bucketName).Delete(ctx)
}

func (g googleStorageClient) Disconnect() error {
	return g.Client.Close()
}

func (g googleStorageClient) SimpleDisconnect() {
	err := g.Disconnect()
	if err != nil {
		logger.Error("error disconnect google storage:", err)
		return
	}
	logger.Info("connection to google storage closed.")
}

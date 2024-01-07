package cstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type CreateGoogleBucketInput struct {
	BucketName string `validate:"required"`
	ProjectId  string `validate:"required"`
	Attrs      *storage.BucketAttrs
}

type UploadFileInput struct {
	BucketName string       `validate:"required"`
	FileName   string       `validate:"required"`
	MimeType   FileMimeType `validate:"required,enum"`
	Content    any          `validate:"required"`
}

type UploadFilesOutput struct {
	Attrs *storage.ObjectAttrs
	Err   error
}

type DeleteFileInput struct {
	BucketName string `validate:"required"`
	FileName   string `validate:"required"`
}

type DeleteFolderInput struct {
	BucketName string `validate:"required"`
	FolderName string `validate:"required"`
}

type DeleteFileOutput struct {
	BucketName string
	FileName   string
	Err        error
}

type DeleteFolderOutput struct {
	BucketName string
	FolderName string
	Err        error
}

type DeleteBucketOutput struct {
	BucketName string
	Err        error
}

type googleStorage struct {
	client *storage.Client
}

type GoogleStorage interface {
	CreateBucket(ctx context.Context, input CreateGoogleBucketInput) error
	GetBucket(ctx context.Context, bucketName string) (*storage.BucketAttrs, error)
	UploadFile(ctx context.Context, input UploadFileInput) (*storage.ObjectAttrs, error)
	UploadFiles(ctx context.Context, inputs ...UploadFileInput) []UploadFilesOutput
	FindOneByName(ctx context.Context, bucketName, fileName string) (*storage.ObjectAttrs, error)
	Find(ctx context.Context, bucketName string, opts ...OptionFind) ([]storage.ObjectAttrs, error)
	DeleteFile(ctx context.Context, input DeleteFileInput) error
	DeleteFiles(ctx context.Context, inputs ...DeleteFileInput) []DeleteFileOutput
	DeleteFolder(ctx context.Context, input DeleteFolderInput) error
	DeleteFolders(ctx context.Context, input ...DeleteFolderInput) []DeleteFolderOutput
	DeleteBucket(ctx context.Context, bucketName string) error
	DeleteBuckets(ctx context.Context, bucketName ...string) []DeleteBucketOutput
	Disconnect() error
	SimpleDisconnect()
}

func NewGoogleStorage(ctx context.Context, opts ...option.ClientOption) (GoogleStorage, error) {
	client, err := storage.NewClient(ctx, opts...)
	if helper.IsNotNil(err != nil) {
		return nil, err
	}
	return googleStorage{
		client: client,
	}, nil
}

func (g googleStorage) CreateBucket(ctx context.Context, input CreateGoogleBucketInput) error {
	if err := helper.Validate().Struct(input); helper.IsNotNil(err) {
		return err
	}
	return g.client.Bucket(input.BucketName).Create(ctx, input.ProjectId, input.Attrs)
}

func (g googleStorage) GetBucket(ctx context.Context, bucketName string) (*storage.BucketAttrs, error) {
	return g.client.Bucket(bucketName).Attrs(ctx)
}

func (g googleStorage) UploadFile(ctx context.Context, input UploadFileInput) (*storage.ObjectAttrs, error) {
	if err := helper.Validate().Struct(input); helper.IsNotNil(err) {
		return nil, err
	}
	f := input.FileName + input.MimeType.Extension()
	obj := g.client.Bucket(input.BucketName).Object(f)
	fw := obj.NewWriter(ctx)
	fw.ContentType = input.MimeType.String()
	bytesContent, err := helper.ConvertToBytes(input.Content)
	if helper.IsNotNil(err) {
		return nil, err
	}
	_, err = fw.Write(bytesContent)
	if helper.IsNotNil(err) {
		return nil, err
	}
	err = fw.Close()
	if helper.IsNotNil(err) {
		return nil, err
	}
	return obj.Attrs(ctx)
}

func (g googleStorage) UploadFiles(ctx context.Context, inputs ...UploadFileInput) []UploadFilesOutput {
	var result []UploadFilesOutput
	for _, input := range inputs {
		attrs, err := g.UploadFile(ctx, input)
		result = append(result, UploadFilesOutput{
			Attrs: attrs,
			Err:   err,
		})
	}
	return result
}

func (g googleStorage) FindOneByName(ctx context.Context, bucketName, fileName string) (*storage.ObjectAttrs, error) {
	bkt := g.client.Bucket(bucketName)
	return bkt.Object(fileName).Attrs(ctx)
}

func (g googleStorage) Find(ctx context.Context, bucketName string, opts ...OptionFind) ([]storage.ObjectAttrs, error) {
	opt := GetOptionFindByParams(opts)
	bkt := g.client.Bucket(bucketName)
	objs := bkt.Objects(ctx, &storage.Query{
		Delimiter:                opt.Delimiter,
		Prefix:                   opt.Prefix,
		Versions:                 opt.Versions,
		StartOffset:              opt.StartOffset,
		EndOffset:                opt.EndOffset,
		Projection:               storage.Projection(opt.Projection.Int()),
		IncludeTrailingDelimiter: opt.IncludeTrailingDelimiter,
		MatchGlob:                opt.MatchGlob,
	})
	var result []storage.ObjectAttrs
	for {
		obj, err := objs.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if helper.IsNotNil(err) {
			return nil, err
		}
		result = append(result, *obj)
	}
	return result, nil
}

func (g googleStorage) DeleteFile(ctx context.Context, input DeleteFileInput) error {
	if err := helper.Validate().Struct(input); helper.IsNotNil(err) {
		return err
	}
	bkt := g.client.Bucket(input.BucketName)
	return bkt.Object(input.FileName).Delete(ctx)
}

func (g googleStorage) DeleteFiles(ctx context.Context, inputs ...DeleteFileInput) []DeleteFileOutput {
	var result []DeleteFileOutput
	for _, input := range inputs {
		err := g.DeleteFile(ctx, input)
		result = append(result, DeleteFileOutput{
			BucketName: input.BucketName,
			FileName:   input.FileName,
			Err:        err,
		})
	}
	return result
}

func (g googleStorage) DeleteFolder(ctx context.Context, input DeleteFolderInput) error {
	if err := helper.Validate().Struct(input); helper.IsNotNil(err) {
		return err
	}
	bkt := g.client.Bucket(input.BucketName)
	objs := bkt.Objects(ctx, &storage.Query{Prefix: input.FolderName})
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

func (g googleStorage) DeleteFolders(ctx context.Context, inputs ...DeleteFolderInput) []DeleteFolderOutput {
	var result []DeleteFolderOutput
	for _, input := range inputs {
		err := g.DeleteFolder(ctx, input)
		result = append(result, DeleteFolderOutput{
			BucketName: input.BucketName,
			FolderName: input.FolderName,
			Err:        err,
		})
	}
	return result
}

func (g googleStorage) DeleteBucket(ctx context.Context, bucketName string) error {
	return g.client.Bucket(bucketName).Delete(ctx)
}

func (g googleStorage) DeleteBuckets(ctx context.Context, bucketName ...string) []DeleteBucketOutput {
	var result []DeleteBucketOutput
	for _, bktName := range bucketName {
		err := g.DeleteBucket(ctx, bktName)
		result = append(result, DeleteBucketOutput{
			BucketName: bktName,
			Err:        err,
		})
	}
	return result
}

func (g googleStorage) Disconnect() error {
	return g.client.Close()
}

func (g googleStorage) SimpleDisconnect() {
	err := g.Disconnect()
	if err != nil {
		logger.Error("error disconnect google storage:", err)
		return
	}
	logger.Info("connection to google storage closed.")
}

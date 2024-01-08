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

// CreateGoogleBucketInput input for creating a Google storage bucket
type CreateGoogleBucketInput struct {
	// BucketName name of the bucket to be created (required)
	BucketName string `validate:"required"`
	// ProjectId project id where the bucket will be created (required)
	ProjectId string `validate:"required"`
	// Attrs metadata options to customize bucket creation
	Attrs *storage.BucketAttrs
}

// UploadFileInput input for creating/updating a new object in the bucket
type UploadFileInput struct {
	// BucketName name of the bucket where the object will be created (required)
	BucketName string `validate:"required"`
	// FileName name of the object that will be created (required)
	FileName string `validate:"required"`
	// MimeType type of content of the object that will be created (required)
	MimeType FileMimeType `validate:"required,enum"`
	// Content of the object that will be created (required)
	Content any `validate:"required"`
}

// UploadFilesOutput output for creating/updating multiple objects in the bucket
type UploadFilesOutput struct {
	// Attrs metadata for the object created/updated
	Attrs *storage.ObjectAttrs
	// Err error occurred when creating/updating the object
	Err error
}

// DeleteFileInput input for removing an object from the bucket
type DeleteFileInput struct {
	// BucketName name of the bucket where the object will be deleted (required)
	BucketName string `validate:"required"`
	// FileName name of the object to be deleted (required)
	FileName string `validate:"required"`
}

// DeleteFolderInput input to remove a folder (prefix) of objects from the bucket
type DeleteFolderInput struct {
	// BucketName name of the bucket where the objects will be deleted (required)
	BucketName string `validate:"required"`
	// FolderName folder name (prefix) where the objects will be deleted (required)
	FolderName string `validate:"required"`
}

// DeleteFileOutput output of removing several objects from the bucket
type DeleteFileOutput struct {
	// BucketName name of the bucket where the object was deleted
	BucketName string
	// FileName name of the object that was deleted
	FileName string
	// Err error occurred when deleting the object
	Err error
}

// DeleteFolderOutput output of removing multiple object folders from bucket
type DeleteFolderOutput struct {
	// BucketName name of the bucket where the objects were deleted
	BucketName string
	// FolderName name of the folder (prefix) that was deleted
	FolderName string
	// Err an error occurred while deleting objects from the folder
	Err error
}

// DeleteBucketOutput output of removing multiple buckets
type DeleteBucketOutput struct {
	// BucketName deleted bucket name
	BucketName string
	// Err an error occurred while deleting the bucket
	Err error
}

type googleStorage struct {
	client *storage.Client
}

type GoogleStorage interface {
	// CreateBucket creates the Bucket in the project
	// If attrs is nil the API defaults will be used
	CreateBucket(ctx context.Context, input CreateGoogleBucketInput) error
	// GetBucket returns the metadata for the bucket
	GetBucket(ctx context.Context, bucketName string) (*storage.BucketAttrs, error)
	// UploadFile set the value passed in the indicated bucket
	UploadFile(ctx context.Context, input UploadFileInput) (*storage.ObjectAttrs, error)
	// UploadFiles set multiple values passed in the indicated bucket
	UploadFiles(ctx context.Context, inputs ...UploadFileInput) []UploadFilesOutput
	// FindOneByName returns the metadata for the object by name
	FindOneByName(ctx context.Context, bucketName, fileName string) (*storage.ObjectAttrs, error)
	// Find returns a list of object metadata, customize the search using the opts parameter (OptionFind)
	Find(ctx context.Context, bucketName string, opts ...OptionFind) ([]storage.ObjectAttrs, error)
	// DeleteFile deletes the single specified object
	DeleteFile(ctx context.Context, input DeleteFileInput) error
	// DeleteFiles deletes multiple objects specified in the input
	DeleteFiles(ctx context.Context, inputs ...DeleteFileInput) []DeleteFileOutput
	// DeleteFolder deletes all objects from a folder (prefix)
	DeleteFolder(ctx context.Context, input DeleteFolderInput) error
	// DeleteFolders deletes all objects from all folders (prefix) mentioned in the input
	DeleteFolders(ctx context.Context, input ...DeleteFolderInput) []DeleteFolderOutput
	// DeleteBucket deletes the Bucket
	DeleteBucket(ctx context.Context, bucketName string) error
	// DeleteBuckets deletes multiple buckets mentioned in the input
	DeleteBuckets(ctx context.Context, bucketName ...string) []DeleteBucketOutput
	// Disconnect close connect to google storage
	Disconnect() error
	// SimpleDisconnect close connect to google storage, without error
	SimpleDisconnect()
}

// NewGoogleStorage new instance of connection with google storage, to close it just use Disconnect() or SimpleDisconnect()
func NewGoogleStorage(ctx context.Context, opts ...option.ClientOption) (GoogleStorage, error) {
	client, err := storage.NewClient(ctx, opts...)
	return googleStorage{
		client: client,
	}, err
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

package cstorage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"google.golang.org/api/option"
)

// CreateBucketInput input for creating a bucket
type CreateBucketInput struct {
	// Bucket name of the bucket to be created (required)
	Bucket string
	// ProjectId project id where the bucket will be created (required only google storage)
	ProjectId string
	// Location bucket, if empty using default region
	Location string
}

// PutObjectInput input for creating/updating an object in the bucket
type PutObjectInput struct {
	// Bucket name of the bucket where the object will be created (required)
	Bucket string
	// Key of the object that will be created (required)
	Key string
	// MimeType type of content of the object that will be created (required)
	MimeType MimeType
	// Content of the object that will be created (required)
	Content any
}

// DeletePrefixInput input to remove a folder (prefix) of objects from the bucket
type DeletePrefixInput struct {
	// Bucket name of the bucket where the objects will be deleted (required)
	Bucket string
	// Prefix name where the objects will be deleted (required)
	Prefix string
}

// PutObjectOutput output for creating/updating multiple objects in the bucket
type PutObjectOutput struct {
	// Bucket name of the bucket where the object will be created
	Bucket string
	// Key of the object that will be created
	Key string
	// Err error occurred when putting the object
	Err error
}

// DeleteObjectInput input for removing an object from the bucket
type DeleteObjectInput struct {
	// Bucket name of the bucket where the object will be deleted (required)
	Bucket string
	// Key of the object to be deleted (required)
	Key string
}

// DeleteObjectsOutput output of removing several objects from the bucket
type DeleteObjectsOutput struct {
	// Bucket name of the bucket where the object was deleted
	Bucket string
	// Key of the object that was deleted
	Key string
	// Err error occurred when deleting the object
	Err error
}

// DeletePrefixOutput output of removing multiple object folders from bucket
type DeletePrefixOutput struct {
	// Bucket name of the bucket where the objects were deleted
	Bucket string
	// Prefix that was deleted
	Prefix string
	// Err an error occurred while deleting objects from the folder
	Err error
}

// DeleteBucketsOutput output of removing multiple buckets
type DeleteBucketsOutput struct {
	// Bucket deleted bucket name
	Bucket string
	// Err an error occurred while deleting the bucket
	Err error
}

type cStorage struct {
	storageSelected     storageSelected
	googleStorageClient iGoogleStorageClient
	awsS3               iAwsS3Client
}

// CStorage interface de integração com armazenamento em nuvem como Google Storage e AWS S3
type CStorage interface {
	// CreateBucket creates the Bucket in the project.
	CreateBucket(ctx context.Context, input CreateBucketInput) error
	// PutObject set the value passed in the indicated bucket
	PutObject(ctx context.Context, input PutObjectInput) error
	// PutObjects set multiple values passed in the indicated bucket
	PutObjects(ctx context.Context, inputs ...PutObjectInput) []PutObjectOutput
	// GetObjectByKey returns the data for the object by name
	GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error)
	// ListObjects return list objects by bucket, custom query using opts param (OptsListObjects)
	ListObjects(ctx context.Context, bucket string, opts ...OptsListObjects) ([]ObjectSummary, error)
	// DeleteObject deletes the single specified object
	DeleteObject(ctx context.Context, input DeleteObjectInput) error
	// DeleteObjects deletes multiple objects specified in the input
	DeleteObjects(ctx context.Context, inputs ...DeleteObjectInput) []DeleteObjectsOutput
	// DeleteObjectsByPrefix deletes all objects from a folder (prefix)
	DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error
	// DeleteObjectsByPrefixes deletes all objects from all folders (prefix) mentioned in the input
	DeleteObjectsByPrefixes(ctx context.Context, input ...DeletePrefixInput) []DeletePrefixOutput
	// DeleteBucket deletes the Bucket
	DeleteBucket(ctx context.Context, bucket string) error
	// DeleteBuckets deletes multiple buckets mentioned in the input
	DeleteBuckets(ctx context.Context, buckets ...string) []DeleteBucketsOutput
	// Disconnect close connect to google storage
	Disconnect() error
	// SimpleDisconnect close connect to google storage, without error
	SimpleDisconnect()
}

// NewGoogleStorage new instance of connection with Google storage, to close it just use Disconnect() or SimpleDisconnect()
func NewGoogleStorage(ctx context.Context, opts ...option.ClientOption) (CStorage, error) {
	client, err := newGoogleStorageClient(ctx, opts...)
	return cStorage{
		storageSelected:     googleStorage,
		googleStorageClient: client,
		awsS3:               nil,
	}, err
}

// NewAwsS3Storage new instance of connection with AWS S3 storage, to close it just use Disconnect() or SimpleDisconnect()
func NewAwsS3Storage(cfg aws.Config) (CStorage, error) {
	return cStorage{
		storageSelected:     awsStorage,
		googleStorageClient: nil,
		awsS3:               newAwsS3StorageClient(cfg),
	}, nil
}

func (g cStorage) CreateBucket(ctx context.Context, input CreateBucketInput) error {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.CreateBucket(ctx, input)
	case awsStorage:
		return g.awsS3.CreateBucket(ctx, input)
	default:
		return ErrWithoutStorageClient
	}
}

func (g cStorage) PutObject(ctx context.Context, input PutObjectInput) error {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.PutObject(ctx, input)
	case awsStorage:
		return g.awsS3.PutObject(ctx, input)
	default:
		return ErrWithoutStorageClient
	}
}

func (g cStorage) PutObjects(ctx context.Context, inputs ...PutObjectInput) []PutObjectOutput {
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

func (g cStorage) GetObjectByKey(ctx context.Context, bucket, key string) (*Object, error) {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.GetObjectByKey(ctx, bucket, key)
	case awsStorage:
		return g.awsS3.GetObjectByKey(ctx, bucket, key)
	default:
		return nil, ErrWithoutStorageClient
	}
}

func (g cStorage) ListObjects(ctx context.Context, bucket string, opts ...OptsListObjects) ([]ObjectSummary, error) {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.ListObjects(ctx, bucket, opts...)
	case awsStorage:
		return g.awsS3.ListObjects(ctx, bucket, opts...)
	default:
		return nil, ErrWithoutStorageClient
	}
}

func (g cStorage) DeleteObject(ctx context.Context, input DeleteObjectInput) error {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.DeleteObject(ctx, input)
	case awsStorage:
		return g.awsS3.DeleteObject(ctx, input)
	default:
		return ErrWithoutStorageClient
	}
}

func (g cStorage) DeleteObjects(ctx context.Context, inputs ...DeleteObjectInput) []DeleteObjectsOutput {
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

func (g cStorage) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.DeleteObjectsByPrefix(ctx, input)
	case awsStorage:
		return g.awsS3.DeleteObjectsByPrefix(ctx, input)
	default:
		return ErrWithoutStorageClient
	}
}

func (g cStorage) DeleteObjectsByPrefixes(ctx context.Context, inputs ...DeletePrefixInput) []DeletePrefixOutput {
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

func (g cStorage) DeleteBucket(ctx context.Context, bucket string) error {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.DeleteBucket(ctx, bucket)
	case awsStorage:
		return g.awsS3.DeleteBucket(ctx, bucket)
	default:
		return ErrWithoutStorageClient
	}
}

func (g cStorage) DeleteBuckets(ctx context.Context, buckets ...string) []DeleteBucketsOutput {
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

func (g cStorage) Disconnect() error {
	switch g.storageSelected {
	case googleStorage:
		return g.googleStorageClient.Disconnect()
	case awsStorage:
		return nil
	default:
		return ErrWithoutStorageClient
	}
}

func (g cStorage) SimpleDisconnect() {
	switch g.storageSelected {
	case googleStorage:
		g.googleStorageClient.SimpleDisconnect()
	}
}

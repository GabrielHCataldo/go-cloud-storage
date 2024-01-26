package cstorage

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
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

type CStorage struct {
	storageSelected     storageSelected
	googleStorageClient iGoogleStorageClient
	awsS3               iAwsS3Client
}

// NewGoogleStorage new instance of connection with Google storage, to close it just use Disconnect() or SimpleDisconnect()
func NewGoogleStorage(ctx context.Context, opts ...option.ClientOption) (*CStorage, error) {
	client, err := newGoogleStorageClient(ctx, opts...)
	return &CStorage{
		storageSelected:     googleStorage,
		googleStorageClient: client,
		awsS3:               nil,
	}, errors.NewSkipCaller(2, err)
}

// NewAwsS3Storage new instance of connection with AWS S3 storage, to close it just use Disconnect() or SimpleDisconnect()
func NewAwsS3Storage(cfg aws.Config) (*CStorage, error) {
	return &CStorage{
		storageSelected:     awsStorage,
		googleStorageClient: nil,
		awsS3:               newAwsS3StorageClient(cfg),
	}, nil
}

// CreateBucket creates the Bucket in the project.
func (c *CStorage) CreateBucket(ctx context.Context, input CreateBucketInput) (err error) {
	switch c.storageSelected {
	case googleStorage:
		err = c.googleStorageClient.CreateBucket(ctx, input)
		break
	case awsStorage:
		err = c.awsS3.CreateBucket(ctx, input)
		break
	default:
		return errWithoutStorageClient(2)
	}
	return errors.NewSkipCaller(2, err)
}

// PutObject set the value passed in the indicated bucket
func (c *CStorage) PutObject(ctx context.Context, input PutObjectInput) (err error) {
	return c.putObject(ctx, input)
}

// PutObjects set multiple values passed in the indicated bucket
func (c *CStorage) PutObjects(ctx context.Context, inputs ...PutObjectInput) []PutObjectOutput {
	var result []PutObjectOutput
	for _, input := range inputs {
		err := c.putObject(ctx, input)
		result = append(result, PutObjectOutput{
			Bucket: input.Bucket,
			Key:    input.Key,
			Err:    err,
		})
	}
	return result
}

// GetObjectByKey returns the data for the object by name
func (c *CStorage) GetObjectByKey(ctx context.Context, bucket, key string) (obj *Object, err error) {
	switch c.storageSelected {
	case googleStorage:
		obj, err = c.googleStorageClient.GetObjectByKey(ctx, bucket, key)
		break
	case awsStorage:
		obj, err = c.awsS3.GetObjectByKey(ctx, bucket, key)
		break
	default:
		return nil, errWithoutStorageClient(2)
	}
	return obj, errors.NewSkipCaller(2, err)
}

// GetObjectUrl returns the object public url
func (c *CStorage) GetObjectUrl(bucket, key string) string {
	switch c.storageSelected {
	case googleStorage:
		return c.googleStorageClient.GetObjectUrl(bucket, key)
	case awsStorage:
		return c.awsS3.GetObjectUrl(bucket, key)
	default:
		return ""
	}
}

// ListObjects return list objects by bucket, custom query using opts param (OptsListObjects)
func (c *CStorage) ListObjects(ctx context.Context, bucket string, opts ...*OptsListObjects) (
	sliceObjSummary []ObjectSummary, err error) {
	switch c.storageSelected {
	case googleStorage:
		sliceObjSummary, err = c.googleStorageClient.ListObjects(ctx, bucket, opts...)
		break
	case awsStorage:
		sliceObjSummary, err = c.awsS3.ListObjects(ctx, bucket, opts...)
		break
	default:
		return nil, errWithoutStorageClient(2)
	}
	return sliceObjSummary, errors.NewSkipCaller(2, err)
}

// DeleteObject deletes the single specified object
func (c *CStorage) DeleteObject(ctx context.Context, input DeleteObjectInput) (err error) {
	return c.deleteObject(ctx, input)
}

// DeleteObjects deletes multiple objects specified in the input
func (c *CStorage) DeleteObjects(ctx context.Context, inputs ...DeleteObjectInput) []DeleteObjectsOutput {
	var result []DeleteObjectsOutput
	for _, input := range inputs {
		err := c.deleteObject(ctx, input)
		result = append(result, DeleteObjectsOutput{
			Bucket: input.Bucket,
			Key:    input.Key,
			Err:    err,
		})
	}
	return result
}

// DeleteObjectsByPrefix deletes all objects from a folder (prefix)
func (c *CStorage) DeleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) error {
	return c.deleteObjectsByPrefix(ctx, input)
}

// DeleteObjectsByPrefixes deletes all objects from all folders (prefix) mentioned in the input
func (c *CStorage) DeleteObjectsByPrefixes(ctx context.Context, inputs ...DeletePrefixInput) []DeletePrefixOutput {
	var result []DeletePrefixOutput
	for _, input := range inputs {
		err := c.deleteObjectsByPrefix(ctx, input)
		result = append(result, DeletePrefixOutput{
			Bucket: input.Bucket,
			Prefix: input.Prefix,
			Err:    err,
		})
	}
	return result
}

// DeleteBucket deletes the Bucket
func (c *CStorage) DeleteBucket(ctx context.Context, bucket string) error {
	return c.deleteBucket(ctx, bucket)
}

// DeleteBuckets deletes multiple buckets mentioned in the input
func (c *CStorage) DeleteBuckets(ctx context.Context, buckets ...string) []DeleteBucketsOutput {
	var result []DeleteBucketsOutput
	for _, bucket := range buckets {
		err := c.deleteBucket(ctx, bucket)
		result = append(result, DeleteBucketsOutput{
			Bucket: bucket,
			Err:    err,
		})
	}
	return result
}

// Disconnect close connect to google storage
func (c *CStorage) Disconnect() error {
	switch c.storageSelected {
	case googleStorage:
		return c.googleStorageClient.Disconnect()
	case awsStorage:
		return nil
	default:
		return errWithoutStorageClient(2)
	}
}

// SimpleDisconnect close connect to google storage, without error
func (c *CStorage) SimpleDisconnect() {
	switch c.storageSelected {
	case googleStorage:
		c.googleStorageClient.SimpleDisconnect()
	}
}

func (c *CStorage) putObject(ctx context.Context, input PutObjectInput) (err error) {
	switch c.storageSelected {
	case googleStorage:
		err = c.googleStorageClient.PutObject(ctx, input)
		break
	case awsStorage:
		err = c.awsS3.PutObject(ctx, input)
		break
	default:
		return errWithoutStorageClient(3)
	}
	return errors.NewSkipCaller(3, err)
}

func (c *CStorage) deleteObject(ctx context.Context, input DeleteObjectInput) (err error) {
	switch c.storageSelected {
	case googleStorage:
		err = c.googleStorageClient.DeleteObject(ctx, input)
		break
	case awsStorage:
		err = c.awsS3.DeleteObject(ctx, input)
		break
	default:
		return errWithoutStorageClient(3)
	}
	return errors.NewSkipCaller(3, err)
}

func (c *CStorage) deleteObjectsByPrefix(ctx context.Context, input DeletePrefixInput) (err error) {
	switch c.storageSelected {
	case googleStorage:
		err = c.googleStorageClient.DeleteObjectsByPrefix(ctx, input)
		break
	case awsStorage:
		err = c.awsS3.DeleteObjectsByPrefix(ctx, input)
		break
	default:
		return errWithoutStorageClient(3)
	}
	return errors.NewSkipCaller(3, err)
}

func (c *CStorage) deleteBucket(ctx context.Context, bucket string) (err error) {
	switch c.storageSelected {
	case googleStorage:
		err = c.googleStorageClient.DeleteBucket(ctx, bucket)
		break
	case awsStorage:
		err = c.awsS3.DeleteBucket(ctx, bucket)
		break
	default:
		return errWithoutStorageClient(3)
	}
	return errors.NewSkipCaller(3, err)
}

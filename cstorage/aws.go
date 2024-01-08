package cstorage

import "cloud.google.com/go/storage"

type awsS3Storage struct {
	client *storage.Client
}

type AwsS3Storage interface {
}

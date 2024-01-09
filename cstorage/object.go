package cstorage

import (
	"cloud.google.com/go/storage"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"time"
)

type ObjectPage struct {
}

type Object struct {
	Key            string
	Url            string
	MimeType       MimeType
	Content        io.Reader
	Size           int64
	VersionId      string
	LastModifiedAt time.Time
}

type ObjectSummary struct {
	Key            string
	Url            string
	Size           int64
	LastModifiedAt time.Time
}

func (o Object) ParseContent(dest any) error {
	return helper.ConvertToDest(o.Content, dest)
}

func parseAwsS3StorageObject(obj *s3.GetObjectOutput) Object {
	return Object{
		MimeType:       MimeType(helper.ConvertPointerToValue(obj.ContentType)),
		Content:        obj.Body,
		Size:           helper.ConvertPointerToValue(obj.ContentLength),
		VersionId:      helper.ConvertPointerToValue(obj.VersionId),
		LastModifiedAt: helper.ConvertPointerToValue(obj.LastModified),
	}
}

func parseAwsS3StorageObjectSummary(obj types.Object) ObjectSummary {
	return ObjectSummary{
		Key:            helper.ConvertPointerToValue(obj.Key),
		Size:           helper.ConvertPointerToValue(obj.Size),
		LastModifiedAt: helper.ConvertPointerToValue(obj.LastModified),
	}
}

func parseGoogleStorageObject(obj *storage.ObjectAttrs) Object {
	lastModifiedAt := obj.Created
	if helper.IsAfter(obj.Updated, lastModifiedAt) {
		lastModifiedAt = obj.Updated
	}
	if helper.IsAfter(obj.Deleted, lastModifiedAt) {
		lastModifiedAt = obj.Deleted
	}
	return Object{
		Key:            obj.Name,
		MimeType:       MimeType(obj.ContentType),
		Content:        helper.SimpleConvertToReader(obj.ContentEncoding),
		Size:           obj.Size,
		LastModifiedAt: lastModifiedAt,
	}
}

func parseGoogleStorageObjectSummary(obj *storage.ObjectAttrs) ObjectSummary {
	lastModifiedAt := obj.Created
	if helper.IsAfter(obj.Updated, lastModifiedAt) {
		lastModifiedAt = obj.Updated
	}
	if helper.IsAfter(obj.Deleted, lastModifiedAt) {
		lastModifiedAt = obj.Deleted
	}
	return ObjectSummary{
		Key:            obj.Name,
		Size:           obj.Size,
		LastModifiedAt: lastModifiedAt,
	}
}

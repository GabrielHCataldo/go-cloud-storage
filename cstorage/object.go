package cstorage

import (
	"cloud.google.com/go/storage"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"time"
)

type ObjectPage struct {
}

type Object struct {
	Key            string
	Url            string
	MimeType       MimeType
	Content        []byte
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
	return Object{
		Key:            obj.Name,
		MimeType:       MimeType(obj.ContentType),
		Size:           obj.Size,
		LastModifiedAt: obj.Updated,
	}
}

func parseGoogleStorageObjectSummary(obj *storage.ObjectAttrs) ObjectSummary {
	return ObjectSummary{
		Key:            obj.Name,
		Size:           obj.Size,
		LastModifiedAt: obj.Updated,
	}
}

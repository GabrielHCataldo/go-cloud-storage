package main

import (
	"context"
	"github.com/GabrielHCataldo/go-cloud-storage/cstorage"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/aws/aws-sdk-go-v2/config"
	"google.golang.org/api/option"
	"os"
	"time"
)

type testStruct struct {
	Name      string    `json:"name,omitempty"`
	BirthDate time.Time `json:"birthDate,omitempty"`
	Balance   float64   `json:"balance"`
	Emails    []string  `json:"emails,omitempty"`
}

const bucketNameDefault = "go-cloud-storage-example"

func initTestStruct() testStruct {
	return testStruct{
		Name:      "Foo Bar",
		BirthDate: time.Now(),
		Balance:   203.12,
		Emails:    []string{"foobar@gmail.com", "foobar2@gmail.com"},
	}
}

func main() {
	createBucket()
	deleteBucket()
	putObject()
	getObjectByKey()
	listObjects()
	deleteObject()
	deletePrefix()
}

func newInstanceAwsS3Storage() (*cstorage.CStorage, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return cstorage.NewAwsS3Storage(cfg)
}

func newInstanceGoogleStorage() (*cstorage.CStorage, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	return cstorage.NewGoogleStorage(ctx, option.WithCredentialsFile("firebase-admin-sdk.json"))
}

func createBucket() {
	cs, err := newInstanceAwsS3Storage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = cs.CreateBucket(ctx, cstorage.CreateBucketInput{
		// name of the bucket to be created (required)
		Bucket: bucketNameDefault,
		// project id where the bucket will be created (required only google storage)
		ProjectId: os.Getenv("GOOGLE_STORAGE_PROJECT_ID"),
		// bucket, if empty using default region
		Location: "",
	})
	if helper.IsNotNil(err) {
		logger.Error("error create bucket:", err)
	} else {
		logger.Info("bucket", bucketNameDefault, "created successfully!")
	}
}

func deleteBucket() {
	cs, err := newInstanceAwsS3Storage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = cs.DeleteBucket(ctx, bucketNameDefault)
	if helper.IsNotNil(err) {
		logger.Error("error delete bucket:", err)
	} else {
		logger.Info("bucket", bucketNameDefault, "deleted successfully!")
	}
}

func putObject() {
	cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	keyObject := "examples/json-example"
	err = cs.PutObject(ctx, cstorage.PutObjectInput{
		// name of the bucket where the object will be created (required)
		Bucket: "go-cloud-storage",
		// key of the object that will be created (required)
		Key: keyObject,
		// type of content of the object that will be created (required)
		MimeType: cstorage.MimeTypeJson,
		// content of the object that will be created (required)
		Content: initTestStruct(),
	})
	if helper.IsNotNil(err) {
		logger.Error("error put object on bucket:", err)
	} else {
		logger.Info("object", keyObject, "putted successfully!")
	}
}

func getObjectByKey() {
	cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	keyObject := "examples/json-example"
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	obj, err := cs.GetObjectByKey(ctx, "go-cloud-storage", keyObject)
	if helper.IsNotNil(err) {
		logger.Error("error get object by key on bucket:", err)
	} else {
		var dest testStruct
		err = obj.ParseContent(&dest)
		if helper.IsNotNil(err) {
			logger.Error("error parse object content:", err)
		}
		logger.Info("object", keyObject, "obtained successfully! obj:", obj, "content parsed:", dest)
	}
}

func listObjects() {
	cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := cstorage.NewOptsListObjects().SetPrefix("").SetDelimiter("")
	objs, err := cs.ListObjects(ctx, "go-cloud-storage", opt)
	if helper.IsNotNil(err) {
		logger.Error("error list objects on bucket:", err)
	} else {
		logger.Info("list objects obtained successfully! objs:", objs)
	}
}

func deleteObject() {
	cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	keyObject := "examples/json-example"
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = cs.DeleteObject(ctx, cstorage.DeleteObjectInput{
		Bucket: "go-cloud-storage",
		Key:    keyObject,
	})
	if helper.IsNotNil(err) {
		logger.Error("error delete object by key on bucket:", err)
	} else {
		logger.Info("object", keyObject, "deleted successfully!")
	}
}

func deletePrefix() {
	cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
	if helper.IsNotNil(err) {
		logger.Error("error create new instance cloud storage:", err)
		return
	}
	defer cs.SimpleDisconnect()
	prefix := "examples/"
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = cs.DeleteObjectsByPrefix(ctx, cstorage.DeletePrefixInput{
		Bucket: "go-cloud-storage",
		Prefix: prefix,
	})
	if helper.IsNotNil(err) {
		logger.Error("error delete objects by prefix on bucket:", err)
	} else {
		logger.Info("prefix", prefix, "deleted successfully!")
	}
}

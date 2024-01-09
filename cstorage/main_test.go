package cstorage

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/aws/aws-sdk-go-v2/config"
	"google.golang.org/api/option"
	"os"
	"testing"
	"time"
)

const googleStorageProjectId = "GOOGLE_STORAGE_PROJECT_ID"
const bucketNameDefault = "go-cloud-storage"
const objectKeyDefault = "object-test"

type testStruct struct {
	Name      string    `json:"name,omitempty"`
	BirthDate time.Time `json:"birthDate,omitempty"`
	Balance   float64   `json:"balance"`
	Emails    []string  `json:"emails,omitempty"`
}

type testCreateBucket struct {
	name     string
	input    CreateBucketInput
	cstorage CStorage
	wantErr  bool
}

type testPutObject struct {
	name     string
	input    PutObjectInput
	cstorage CStorage
	wantErr  bool
}

type testGetObjectByKey struct {
	name     string
	key      string
	cstorage CStorage
	wantErr  bool
}

func TestMain(m *testing.M) {
	m.Run()
}

func initGoogleStorage() CStorage {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	cs, err := NewGoogleStorage(ctx, option.WithCredentialsFile("../firebase-admin-sdk.json"))
	if helper.IsNotNil(err) {
		logger.Error("error start google storage:", err)
		return nil
	}
	return cs
}

func initAwsS3Storage() CStorage {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx)
	if helper.IsNotNil(err) {
		logger.Error("error get aws config default:", err)
		return nil
	}
	cs, err := NewAwsS3Storage(cfg)
	if helper.IsNotNil(err) {
		logger.Error("error start aws s3 storage:", err)
		return nil
	}
	_ = cs.CreateBucket(ctx, CreateBucketInput{
		Bucket:   bucketNameDefault,
		Location: "sa-east-1",
	})
	return cs
}

func initObject(cs CStorage) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err := cs.PutObject(ctx, PutObjectInput{
		Bucket:   bucketNameDefault,
		Key:      objectKeyDefault,
		MimeType: MimeTypeJson,
		Content:  initTestStruct(),
	})
	if helper.IsNotNil(err) {
		logger.Error("error init object on storage:", err)
	} else {
		logger.Info("init object on storage completed successfully!")
	}
}

func initListTestCreateBucket() []testCreateBucket {
	return []testCreateBucket{
		{
			name:     "success google",
			input:    initTestCreateBucketInput(""),
			cstorage: initGoogleStorage(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			input:    initTestCreateBucketInput("sa-east-1"),
			cstorage: initAwsS3Storage(),
			wantErr:  false,
		},
		{
			name:     "failed google",
			cstorage: initGoogleStorage(),
			wantErr:  true,
		},
		{
			name:     "failed aws",
			cstorage: initAwsS3Storage(),
			wantErr:  true,
		},
		{
			name:     "failed instance",
			cstorage: cStorage{},
			wantErr:  true,
		},
	}
}

func initListTestPutObject() []testPutObject {
	return []testPutObject{
		{
			name:     "success google",
			input:    initTestPutObjectInput(),
			cstorage: initGoogleStorage(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			input:    initTestPutObjectInput(),
			cstorage: initAwsS3Storage(),
			wantErr:  false,
		},
		{
			name:     "failed google",
			cstorage: initGoogleStorage(),
			wantErr:  true,
		},
		{
			name:     "failed aws",
			cstorage: initAwsS3Storage(),
			wantErr:  true,
		},
		{
			name:     "failed instance",
			cstorage: cStorage{},
			wantErr:  true,
		},
	}
}

func initListTestGetObjectByKey() []testGetObjectByKey {
	return []testGetObjectByKey{
		{
			name:     "success google",
			key:      objectKeyDefault,
			cstorage: initGoogleStorage(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			key:      objectKeyDefault,
			cstorage: initAwsS3Storage(),
			wantErr:  false,
		},
		{
			name:     "failed google",
			cstorage: initGoogleStorage(),
			wantErr:  true,
		},
		{
			name:     "failed aws",
			cstorage: initAwsS3Storage(),
			wantErr:  true,
		},
		{
			name:     "failed instance",
			cstorage: cStorage{},
			wantErr:  true,
		},
	}
}

func initTestStruct() *testStruct {
	return &testStruct{
		Name:      "Foo Bar",
		BirthDate: time.Now(),
		Balance:   203.12,
		Emails:    []string{"foobar@gmail.com", "foobar2@gmail.com"},
	}
}

func initTestCreateBucketInput(location string) CreateBucketInput {
	return CreateBucketInput{
		Bucket:    bucketNameDefault + "-" + helper.SimpleConvertToString(time.Now().UnixMilli()),
		ProjectId: os.Getenv(googleStorageProjectId),
		Location:  location,
	}
}

func initTestPutObjectInput() PutObjectInput {
	return PutObjectInput{
		Bucket:   bucketNameDefault,
		Key:      helper.SimpleConvertToString(time.Now().UnixMilli()),
		MimeType: MimeTypeJson,
		Content:  initTestStruct(),
	}
}

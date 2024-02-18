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
const bucketNameToDeleteDefault = "go-cloud-storage-to-delete"
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

type testListObjects struct {
	name     string
	cstorage CStorage
	bucket   string
	opts     *OptsListObjects
	wantErr  bool
}

type testDeleteObject struct {
	name     string
	cstorage CStorage
	input    DeleteObjectInput
	wantErr  bool
}

type testDeleteObjectsByPrefix struct {
	name     string
	cstorage CStorage
	input    DeletePrefixInput
	wantErr  bool
}

type testDeleteBucket struct {
	name     string
	cstorage CStorage
	bucket   string
	wantErr  bool
}

type testDisconnect struct {
	name     string
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
	cs := NewAwsS3Storage(cfg)
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

func initBucket(cs CStorage) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err := cs.CreateBucket(ctx, CreateBucketInput{
		Bucket:    bucketNameToDeleteDefault,
		ProjectId: os.Getenv(googleStorageProjectId),
	})
	if helper.IsNotNil(err) {
		logger.Error("error init bucket to delete on storage:", err)
	} else {
		logger.Info("init bucket to delete on storage completed successfully!")
	}
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
	googleStorage := initGoogleStorage()
	return []testCreateBucket{
		{
			name:     "success google",
			input:    initTestCreateBucketInput(""),
			cstorage: googleStorage,
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
	}
}

func initListTestListObjects() []testListObjects {
	return []testListObjects{
		{
			name:     "success google",
			cstorage: initGoogleStorage(),
			bucket:   bucketNameDefault,
			wantErr:  false,
		},
		{
			name:     "success empty google",
			cstorage: initGoogleStorage(),
			bucket:   bucketNameDefault,
			opts:     initTestOptsListObjects(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			cstorage: initAwsS3Storage(),
			bucket:   bucketNameDefault,
			wantErr:  false,
		},
		{
			name:     "success empty aws",
			cstorage: initAwsS3Storage(),
			bucket:   bucketNameDefault,
			opts:     initTestOptsListObjects(),
			wantErr:  false,
		},
		{
			name:     "failed google",
			cstorage: initGoogleStorage(),
			bucket:   "bucket-not-exists",
			wantErr:  true,
		},
		{
			name:     "failed aws",
			cstorage: initAwsS3Storage(),
			bucket:   "bucket-not-exists",
			wantErr:  true,
		},
	}
}

func initListTestDeleteObject() []testDeleteObject {
	return []testDeleteObject{
		{
			name:     "success google",
			cstorage: initGoogleStorage(),
			input:    initTestDeleteObjectInput(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			input:    initTestDeleteObjectInput(),
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
	}
}

func initListTestDeleteObjectsByPrefix() []testDeleteObjectsByPrefix {
	return []testDeleteObjectsByPrefix{
		{
			name:     "success google",
			cstorage: initGoogleStorage(),
			input:    initTestDeletePrefixInput(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			input:    initTestDeletePrefixInput(),
			cstorage: initAwsS3Storage(),
			wantErr:  false,
		},
		{
			name:     "failed google",
			cstorage: initGoogleStorage(),
			input: DeletePrefixInput{
				Bucket: "not-exists",
				Prefix: "",
			},
			wantErr: true,
		},
		{
			name:     "failed aws",
			cstorage: initAwsS3Storage(),
			input: DeletePrefixInput{
				Bucket: "not-exists",
			},
			wantErr: true,
		},
	}
}

func initListTestDeleteBucket() []testDeleteBucket {
	return []testDeleteBucket{
		{
			name:     "success google",
			cstorage: initGoogleStorage(),
			bucket:   bucketNameToDeleteDefault,
			wantErr:  false,
		},
		{
			name:     "success aws",
			cstorage: initAwsS3Storage(),
			bucket:   bucketNameToDeleteDefault,
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
	}
}

func initListTestDisconnect() []testDisconnect {
	return []testDisconnect{
		{
			name:     "success google",
			cstorage: initGoogleStorage(),
			wantErr:  false,
		},
		{
			name:     "success aws",
			cstorage: initAwsS3Storage(),
			wantErr:  false,
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

func initTestOptsListObjects() *OptsListObjects {
	return NewOptsListObjects().SetPrefix("test").SetDelimiter("test")
}

func initTestDeleteObjectInput() DeleteObjectInput {
	return DeleteObjectInput{
		Bucket: bucketNameDefault,
		Key:    objectKeyDefault,
	}
}

func initTestDeletePrefixInput() DeletePrefixInput {
	return DeletePrefixInput{
		Bucket: bucketNameDefault,
		Prefix: "",
	}
}

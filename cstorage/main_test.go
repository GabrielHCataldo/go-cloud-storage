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

var cstorage CStorage

const googleStorageProjectId = "GOOGLE_STORAGE_PROJECT_ID"
const bucketNameDefault = "go-cloud-storage"

type testCreateBucket struct {
	name    string
	input   CreateBucketInput
	wantErr bool
}

func TestMain(m *testing.M) {
	initAwsS3Storage()
	m.Run()
}

func initGoogleStorage() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	cs, err := NewGoogleStorage(ctx, option.WithCredentialsFile("../firebase-admin-sdk.json"))
	if helper.IsNotNil(err) {
		logger.Error("error start google storage:", err)
		return
	}
	cstorage = cs
}

func initAwsS3Storage() {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx)
	if helper.IsNotNil(err) {
		logger.Error("error get aws config default:", err)
		return
	}
	cs, err := NewAwsS3Storage(cfg)
	if helper.IsNotNil(err) {
		logger.Error("error start aws s3 storage:", err)
		return
	}
	cstorage = cs
}

func initListTestCreateBucket() []testCreateBucket {
	return []testCreateBucket{
		{
			name: "success",
			input: CreateBucketInput{
				Bucket:    bucketNameDefault + "-" + helper.SimpleConvertToString(time.Now().UnixMilli()),
				ProjectId: os.Getenv(googleStorageProjectId),
				Location:  "sa-east-1",
			},
			wantErr: false,
		},
		{
			name:    "failed",
			wantErr: true,
		},
	}
}

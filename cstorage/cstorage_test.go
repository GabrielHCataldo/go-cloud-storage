package cstorage

import (
	"context"
	"github.com/GabrielHCataldo/go-logger/logger"
	"testing"
	"time"
)

func TestCStorageCreateBucket(t *testing.T) {
	for _, tt := range initListTestCreateBucket() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := tt.cstorage.CreateBucket(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				logger.Errorf("CreateBucket() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStoragePutObject(t *testing.T) {
	for _, tt := range initListTestPutObject() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := tt.cstorage.PutObject(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				logger.Errorf("PutObject() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStoragePutObjects(t *testing.T) {
	for _, tt := range initListTestPutObject() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			output := tt.cstorage.PutObjects(ctx, tt.input)
			logger.Errorf("PutObjects() output = %v", output)
		})
	}
}

func TestCStorageGetObjectByKey(t *testing.T) {
	for _, tt := range initListTestGetObjectByKey() {
		t.Run(tt.name, func(t *testing.T) {
			initObject(tt.cstorage)
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result, err := tt.cstorage.GetObjectByKey(ctx, bucketNameDefault, tt.key)
			if (err != nil) != tt.wantErr {
				logger.Errorf("GetObjectByKey() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			}
			logger.Infof("GetObjectByKey() result = %v, err = %v", result, err)
		})
	}
}

func TestCStorage_ListObjects(t *testing.T) {

}

func TestCStorage_DeleteObject(t *testing.T) {

}

func TestCStorage_DeleteObjects(t *testing.T) {

}

func TestCStorage_DeleteObjectsByPrefix(t *testing.T) {

}

func TestCStorage_DeleteObjectsByPrefixes(t *testing.T) {

}

func TestCStorage_DeleteBucket(t *testing.T) {

}

func TestCStorage_DeleteBuckets(t *testing.T) {

}

func TestCStorage_Disconnect(t *testing.T) {

}

func TestCStorage_SimpleDisconnect(t *testing.T) {

}

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
			err := cstorage.CreateBucket(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				logger.Errorf("CreateBucket() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStorage_PutObject(t *testing.T) {

}

func TestCStorage_PutObjects(t *testing.T) {

}

func TestCStorage_GetObjectByKey(t *testing.T) {

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

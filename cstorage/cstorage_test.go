package cstorage

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
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
			logger.Infof("PutObjects() output = %v", output)
		})
	}
}

func TestCStorageGetObjectByKey(t *testing.T) {
	for _, tt := range initListTestGetObjectByKey() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initObject(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result, err := tt.cstorage.GetObjectByKey(ctx, bucketNameDefault, tt.key)
			if (err != nil) != tt.wantErr {
				logger.Errorf("GetObjectByKey() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			} else if helper.IsNotNil(result) {
				var destContent testStruct
				_ = result.ParseContent(&destContent)
				logger.Infof("GetObjectByKey() result = %v, err = %v", destContent, err)
			} else {
				logger.Infof("GetObjectByKey() result = %v, err = %v", result, err)
			}
		})
	}
}

func TestCStorageGetObjectUrl(t *testing.T) {
	for _, tt := range initListTestGetObjectByKey() {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.cstorage.GetObjectUrl(bucketNameDefault, tt.key)
			logger.Infof("GetObjectUrl() output = %v", output)
		})
	}
}

func TestCStorageListObjects(t *testing.T) {
	for _, tt := range initListTestListObjects() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initObject(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result, err := tt.cstorage.ListObjects(ctx, bucketNameDefault, tt.opts)
			if (err != nil) != tt.wantErr {
				logger.Errorf("ListObjects() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			}
			logger.Infof("ListObjects() result = %v, err = %v", result, err)
		})
	}
}

func TestCStorageDeleteObject(t *testing.T) {
	for _, tt := range initListTestDeleteObject() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initObject(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := tt.cstorage.DeleteObject(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				logger.Errorf("DeleteObject() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStorageDeleteObjects(t *testing.T) {
	for _, tt := range initListTestDeleteObject() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initObject(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			output := tt.cstorage.DeleteObjects(ctx, tt.input)
			logger.Infof("DeleteObjects() output = %v", output)
		})
	}
}

func TestCStorageDeleteObjectsByPrefix(t *testing.T) {
	for _, tt := range initListTestDeleteObjectsByPrefix() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initObject(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := tt.cstorage.DeleteObjectsByPrefix(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				logger.Errorf("DeleteObjectsByPrefix() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStorageDeleteObjectsByPrefixes(t *testing.T) {
	for _, tt := range initListTestDeleteObjectsByPrefix() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initObject(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			output := tt.cstorage.DeleteObjectsByPrefixes(ctx, tt.input)
			logger.Infof("DeleteObjectsByPrefixes() output = %v", output)
		})
	}
}

func TestCStorageDeleteBucket(t *testing.T) {
	for _, tt := range initListTestDeleteBucket() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initBucket(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := tt.cstorage.DeleteBucket(ctx, tt.bucket)
			if (err != nil) != tt.wantErr {
				logger.Errorf("DeleteBucket() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStorageDeleteBuckets(t *testing.T) {
	for _, tt := range initListTestDeleteBucket() {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				initBucket(tt.cstorage)
			}
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			output := tt.cstorage.DeleteBuckets(ctx, tt.bucket)
			logger.Infof("DeleteBuckets() output = %v", output)
		})
	}
}

func TestCStorageDisconnect(t *testing.T) {
	for _, tt := range initListTestDisconnect() {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				tt.cstorage.SimpleDisconnect()
			}
			err := tt.cstorage.Disconnect()
			if (err != nil) != tt.wantErr {
				logger.Errorf("Disconnect() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestCStorageSimpleDisconnect(t *testing.T) {
	for _, tt := range initListTestDisconnect() {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				tt.cstorage.SimpleDisconnect()
			}
			tt.cstorage.SimpleDisconnect()
		})
	}
}

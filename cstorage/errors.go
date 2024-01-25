package cstorage

import "github.com/GabrielHCataldo/go-errors/errors"

var MsgErrWithoutStorageClient = "error cstorage: select storage client"
var ErrWithoutStorageClient = errors.New(MsgErrWithoutStorageClient)

func errWithoutStorageClient(skip int) error {
	ErrWithoutStorageClient = errors.NewSkipCaller(skip+1, MsgErrWithoutStorageClient)
	return ErrWithoutStorageClient
}

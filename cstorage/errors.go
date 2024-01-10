package cstorage

import "errors"

var ErrWithoutStorageClient = errors.New("error cstorage: select storage client")
var ErrPrefixNotExists = errors.New("error cstorage: prefix not found")

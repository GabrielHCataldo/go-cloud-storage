package cstorage

import (
	"github.com/GabrielHCataldo/go-helper/helper"
)

// OptsListObjects bucket object search options
type OptsListObjects struct {
	// Delimiter returns results in a directory-like fashion.
	// Results will contain only objects whose names, aside from the
	// prefix, do not contain delimiter. Objects whose names,
	// aside from the prefix, contain delimiter will have their name,
	// truncated after the delimiter, returned in prefixes.
	// Duplicate prefixes are omitted.
	// Must be set to / when used with the MatchGlob parameter to filter results
	// in a directory-like mode.
	// Optional.
	Delimiter string
	// Prefix is the prefix filter to query objects
	// whose names begin with this prefix.
	// Optional.
	Prefix string
}

// NewOptsListObjects creates a new OptsListObjects instance
func NewOptsListObjects() *OptsListObjects {
	return &OptsListObjects{}
}

// SetDelimiter sets value for the Delimiter field
func (o *OptsListObjects) SetDelimiter(s string) *OptsListObjects {
	o.Delimiter = s
	return o
}

// SetPrefix sets value for the Prefix field
func (o *OptsListObjects) SetPrefix(s string) *OptsListObjects {
	o.Prefix = s
	return o
}

// GetOptListObjectsByParams assembles the OptsGoogleFind object from optional parameters.
func GetOptListObjectsByParams(opts []*OptsListObjects) *OptsListObjects {
	result := &OptsListObjects{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if helper.IsNotEmpty(opt.Delimiter) {
			result.Delimiter = opt.Delimiter
		}
		if helper.IsNotEmpty(opt.Prefix) {
			result.Prefix = opt.Prefix
		}
	}
	return result
}

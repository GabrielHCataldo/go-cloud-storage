package cstorage

import (
	"github.com/GabrielHCataldo/go-helper/helper"
)

type OptionFind struct {
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
	// Versions indicates whether multiple versions of the same
	// object will be included in the results.
	Versions bool
	// attrSelection is used to select only specific fields to be returned by
	// the query. It is set by the user calling SetAttrSelection. These
	// are used by toFieldMask and toFieldSelection for gRPC and HTTP/JSON
	// clients respectively.
	attrSelection []string
	// StartOffset is used to filter results to objects whose names are
	// lexicographically equal to or after startOffset. If endOffset is also set,
	// the objects listed will have names between startOffset (inclusive) and
	// endOffset (exclusive).
	StartOffset string
	// EndOffset is used to filter results to objects whose names are
	// lexicographically before endOffset. If startOffset is also set, the objects
	// listed will have names between startOffset (inclusive) and endOffset (exclusive).
	EndOffset string
	// Projection defines the set of properties to return. It will default to ProjectionFull,
	// which returns all properties. Passing ProjectionNoACL will omit Owner and ACL,
	// which may improve performance when listing many objects.
	Projection Projection
	// IncludeTrailingDelimiter controls how objects which end in a single
	// instance of Delimiter (for example, if Query.Delimiter = "/" and the
	// object name is "foo/bar/") are included in the results. By default, these
	// objects only show up as prefixes. If IncludeTrailingDelimiter is set to
	// true, they will also be included as objects and their metadata will be
	// populated in the returned ObjectAttrs.
	IncludeTrailingDelimiter bool
	// MatchGlob is a glob pattern used to filter results (for example, foo*bar). See
	// https://cloud.google.com/storage/docs/json_api/v1/objects/list#list-object-glob
	// for syntax details. When Delimiter is set in conjunction with MatchGlob,
	// it must be set to /.
	MatchGlob string
}

func NewFind() OptionFind {
	return OptionFind{}
}

func (f OptionFind) SetDelimiter(s string) OptionFind {
	f.Delimiter = s
	return f
}

func (f OptionFind) SetPrefix(s string) OptionFind {
	f.Prefix = s
	return f
}

func (f OptionFind) SetVersions(b bool) OptionFind {
	f.Versions = b
	return f
}

func (f OptionFind) SetAttrSelection(ss []string) OptionFind {
	f.attrSelection = ss
	return f
}

func (f OptionFind) SetStartOffset(s string) OptionFind {
	f.StartOffset = s
	return f
}

func (f OptionFind) SetEndOffset(s string) OptionFind {
	f.EndOffset = s
	return f
}

func (f OptionFind) SetProjection(p Projection) OptionFind {
	f.Projection = p
	return f
}

func (f OptionFind) SetIncludeTrailingDelimiter(b bool) OptionFind {
	f.IncludeTrailingDelimiter = b
	return f
}

func (f OptionFind) SetMatchGlob(s string) OptionFind {
	f.MatchGlob = s
	return f
}

func GetOptionFindByParams(opts []OptionFind) OptionFind {
	result := OptionFind{}
	for _, opt := range opts {
		if helper.IsNotEmpty(opt.Delimiter) {
			result.Delimiter = opt.Delimiter
		}
		if helper.IsNotEmpty(opt.Prefix) {
			result.Prefix = opt.Prefix
		}
		if opt.Versions {
			result.Versions = opt.Versions
		}
		if helper.IsNotEmpty(opt.attrSelection) {
			result.attrSelection = opt.attrSelection
		}
		if helper.IsNotEmpty(opt.StartOffset) {
			result.StartOffset = opt.StartOffset
		}
		if helper.IsNotEmpty(opt.EndOffset) {
			result.EndOffset = opt.EndOffset
		}
		if opt.Projection.IsEnumValid() {
			result.Projection = opt.Projection
		}
		if opt.IncludeTrailingDelimiter {
			result.IncludeTrailingDelimiter = opt.IncludeTrailingDelimiter
		}
		if helper.IsNotEmpty(opt.MatchGlob) {
			result.MatchGlob = opt.MatchGlob
		}
	}
	return result
}

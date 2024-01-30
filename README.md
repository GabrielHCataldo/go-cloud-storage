Go Cloud Storage
=================
<!--suppress ALL -->
<img align="right" src="gopher-cloud.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.1.3-vividgreen.svg)](https://github.com/GabrielHCataldo/go-cloud-storage/releases/tag/v1.1.3)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/go-cloud-storage)](https://goreportcard.com/report/github.com/GabrielHCataldo/go-cloud-storage)
[![Coverage Status](https://coveralls.io/repos/GabrielHCataldo/go-cloud-storage/badge.svg?branch=main&service=github)](https://coveralls.io/github/GabrielHCataldo/go-cloud-storage?branch=main)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/go-cloud-storage/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/go-cloud-storage)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/go-cloud-storage?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/go-cloud-storage/cstorage)
![License](https://img.shields.io/dub/l/vibe-d.svg)

[//]: # ([![build workflow]&#40;https://github.com/GabrielHCataldo/go-cloud-storage/actions/workflows/go.yml/badge.svg&#41;]&#40;https://github.com/GabrielHCataldo/go-cloud-storage/actions&#41;)

[//]: # ([![Source graph]&#40;https://sourcegraph.com/github.com/go-cloud-storage/cstorage/-/badge.svg&#41;]&#40;https://sourcegraph.com/github.com/go-cloud-storage/cstorage?badge&#41;)

[//]: # ([![TODOs]&#40;https://badgen.net/https/api.tickgit.com/badgen/github.com/GabrielHCataldo/go-cloud-storage/cstorage&#41;]&#40;https://www.tickgit.com/browse?repo=github.com/GabrielHCataldo/go-cloud-storage&#41;)

The go-cloud-storage project came to make the use of Cloud Storage easier and more flexible, regardless of the provider, just use a simple and intuitive library interface. Below we list some implemented features:

- Simple bucket creation and deletion regardless of provider.
- Simple object insertion/update without worrying about conversions or pointers.
- Ease of obtaining the object with automatic conversion to the type you want.
- Object listing.
- Removal of object, multiple objects and prefixes.

Implemented providers:

- Google Cloud Storage
- AWS S3

Installation
------------

Use go get.

	go get github.com/GabrielHCataldo/go-cloud-storage

Then import the go-cloud-storage package into your own code.

```go
import "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
```

Usability and documentation
------------
Let's start by instantiating the go-cloud-storage Interface indicating which provider we will use, see:

**NOTE**: Each instance follows the standard of the provider's documentation, just follow it.

**IMPORTANT**: Always check the documentation in the structures and functions fields.
For more details on the examples, visit [All examples link](https://github/GabrielHCataldo/go-cloud-storage/blob/main/_example).

- Google Storage

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "time"
)

func main() {
    cs, err := newInstanceGoogleStorage()
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
    } else {
        logger.Info("cloud storage instance created successfully!")
        cs.SimpleDisconnect()
    }
}

func newInstanceGoogleStorage() (*cstorage.CStorage, error) {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    return cstorage.NewGoogleStorage(ctx, option.WithCredentialsFile("firebase-admin-sdk.json"))
}
```

Output:

    [INFO 2024/01/12 08:17:17] main.go:18: cloud storage instance created successfully!

- AWS S3

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "github.com/aws/aws-sdk-go-v2/config"
    "time"
)

func main() {
    cs, err := newInstanceAwsS3Storage()
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
    } else {
        logger.Info("cloud storage instance created successfully!")
        cs.SimpleDisconnect()
    }
}

func newInstanceAwsS3Storage() (*cstorage.CStorage, error) {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    cfg, err := config.LoadDefaultConfig(ctx)
    if helper.IsNotNil(err) {
        return nil, err
    }
    return cstorage.NewAwsS3Storage(cfg)
}
```

Output: 

    [INFO 2024/01/12 08:21:44] main.go:18: cloud storage instance created successfully!

With the instance created, we will continue with basic examples:

#### Create Bucket
Creating a bucket is very simple, see the example below:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "os"
    "time"
)

const bucketNameDefault = "go-cloud-storage-example"

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    err = cs.CreateBucket(ctx, cstorage.CreateBucketInput{
        // name of the bucket to be created (required)
        Bucket: bucketNameDefault,
        // project id where the bucket will be created (required only google storage)
        ProjectId: os.Getenv("GOOGLE_STORAGE_PROJECT_ID"),
        // bucket, if empty using default region
        Location: "",
    })
    if helper.IsNotNil(err) {
        logger.Error("error create bucket:", err)
    } else {
        logger.Info("bucket", bucketNameDefault, "created successfully!")
    }
}
```

Output:
    
    [INFO 2024/01/12 09:18:05] main.go:51: bucket go-cloud-storage-example created successfully!

For more bucket examples, such as multiple creation,
access the [link](https://github/GabrielHCataldo/go-cloud-storage/blob/main/_example/main).

#### Remove Bucket
Removing a bucket is very simple, just enter the name of the bucket you want to delete, see the example below:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "os"
    "time"
)

const bucketNameDefault = "go-cloud-storage-example"

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    err = cs.DeleteBucket(ctx, bucketNameDefault)
    if helper.IsNotNil(err) {
        logger.Error("error delete bucket:", err)
    } else {
        logger.Info("bucket", bucketNameDefault, "deleted successfully!")
    }
}
```

Output:

    [INFO 2024/01/12 09:22:47] main.go:29: bucket go-cloud-storage-example deleted successfully!

For more bucket examples, such as multiple deletion,
access the [link](https://github/GabrielHCataldo/go-cloud-storage/blob/main/_example/main).

#### Put Object
You can put any type of content in the bucket, see the example:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "os"
    "time"
)

type testStruct struct {
    Name      string    `json:"name,omitempty"`
    BirthDate time.Time `json:"birthDate,omitempty"`
    Balance   float64   `json:"balance"`
    Emails    []string  `json:"emails,omitempty"`
}

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    keyObject := "examples/json-example"
    err = cs.PutObject(ctx, cstorage.PutObjectInput{
        // name of the bucket where the object will be created (required)
        Bucket: "go-cloud-storage",
        // key of the object that will be created (required)
        Key: keyObject,
        // type of content of the object that will be created (required)
        MimeType: cstorage.MimeTypeJson,
        // content of the object that will be created (required)
        Content: initTestStruct(),
    })
    if helper.IsNotNil(err) {
        logger.Error("error put object on bucket:", err)
    } else {
        logger.Info("object", keyObject, "putted successfully!")
    }
}

func initTestStruct() testStruct {
    return testStruct{
        Name:      "Foo Bar",
        BirthDate: time.Now(),
        Balance:   203.12,
        Emails:    []string{"foobar@gmail.com", "foobar2@gmail.com"},
    }
}
```

Output:
    
    [INFO 2024/01/12 09:38:36] main.go:44: object examples/json-example putted successfully!

For more object examples, such as multiple creation,
access [link](https://github/GabrielHCataldo/go-cloud-storage/blob/main/_example/main).

#### Get Object By Key
To obtain a single object, simply pass the name and key of the bucket,
to parse the content to the desired type, simply use the **ParseContent**
function of the returned object, see:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "os"
    "time"
)

type testStruct struct {
    Name      string    `json:"name,omitempty"`
    BirthDate time.Time `json:"birthDate,omitempty"`
    Balance   float64   `json:"balance"`
    Emails    []string  `json:"emails,omitempty"`
}

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    keyObject := "examples/json-example"
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    obj, err := cs.GetObjectByKey(ctx, "go-cloud-storage", keyObject)
    if helper.IsNotNil(err) {
        logger.Error("error get object by key on bucket:", err)
    } else {
        var dest testStruct
        err = obj.ParseContent(&dest)
        if helper.IsNotNil(err) {
            logger.Error("error parse object content:", err)
        }
        logger.Info("object", keyObject, "obtained successfully! obj:", obj, "content parsed:", dest)
    }
}

func initTestStruct() testStruct {
    return testStruct{
        Name:      "Foo Bar",
        BirthDate: time.Now(),
        Balance:   203.12,
        Emails:    []string{"foobar@gmail.com", "foobar2@gmail.com"},
    }
}
```

Output:

    [INFO 2024/01/12 09:49:29] main.go:40: object examples/json-example obtained successfully! obj: {"Key":"examples/json-example","Url":"https://storage.googleapis.com/go-cloud-storage/examples/json-example","MimeType":"application/json","Content":[123,34,110,97,109,101,34,58,34,70,111,111,32,66,97,114,34,44,34,98,105,114,116,104,68,97,116,101,34,58,34,50,48,50,52,45,48,49,45,49,50,84,48,57,58,51,56,58,51,54,46,50,51,50,51,52,50,45,48,51,58,48,48,34,44,34,98,97,108,97,110,99,101,34,58,50,48,51,46,49,50,44,34,101,109,97,105,108,115,34,58,91,34,102,111,111,98,97,114,64,103,109,97,105,108,46,99,111,109,34,44,34,102,111,111,98,97,114,50,64,103,109,97,105,108,46,99,111,109,34,93,125],"Size":132,"VersionId":"","LastModifiedAt":"2024-01-12T12:38:37Z"} content parsed: {"name":"Foo Bar","birthDate":"2024-01-12T09:38:36-03:00","balance":203.12,"emails":["foobar@gmail.com","foobar2@gmail.com"]}

#### List Objects
To list, simply enter which bucket you want, you can customize your searches using the opts parameter, see:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "os"
    "time"
)

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    opt := cstorage.NewOptsListObjects().SetPrefix("").SetDelimiter("")
    objs, err := cs.ListObjects(ctx, "go-cloud-storage", opt)
    if helper.IsNotNil(err) {
        logger.Error("error list objects on bucket:", err)
    } else {
        logger.Info("list objects obtained successfully! objs:", objs)
    }
}

```

Output

    [INFO 2024/01/12 09:59:19] main.go:28: list objects obtained successfully! objs: [{"Key":"examples/json-example","LastModifiedAt":"2024-01-12T12:38:37Z","Size":132,"Url":"https://storage.googleapis.com/go-cloud-storage/examples/json-example"}]

#### Delete object
To remove a specific object, simply enter the bucket and key, see:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "google.golang.org/api/option"
    "os"
    "time"
)

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    keyObject := "examples/json-example"
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    err = cs.DeleteObject(ctx, cstorage.DeleteObjectInput{
        Bucket: "go-cloud-storage",
        Key:    keyObject,
    })
    if helper.IsNotNil(err) {
        logger.Error("error delete object by key on bucket:", err)
    } else {
        logger.Info("object", keyObject, "deleted successfully!")
    }
}
```

Output:

    [INFO 2024/01/12 10:47:09] main.go:31: object examples/json-example deleted successfully!

For more examples of object deletion, 
visit [link](https://github/GabrielHCataldo/go-cloud-storage/blob/main/_example/main).

#### Delete objects by prefix
To remove all objects from a prefix, see:

```go
package main

import (
    "context"
    "github.com/GabrielHCataldo/go-cloud-storage/cstorage"
    "github.com/GabrielHCataldo/go-helper/helper"
    "github.com/GabrielHCataldo/go-logger/logger"
    "github.com/aws/aws-sdk-go-v2/config"
    "google.golang.org/api/option"
    "os"
    "time"
)

func main() {
    cs, err := newInstanceGoogleStorage() // or newInstanceAwsS3Storage
    if helper.IsNotNil(err) {
        logger.Error("error create new instance cloud storage:", err)
        return
    }
    defer cs.SimpleDisconnect()
    prefix := "examples/"
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    err = cs.DeleteObjectsByPrefix(ctx, cstorage.DeletePrefixInput{
        Bucket: "go-cloud-storage",
        Prefix: prefix,
    })
    if helper.IsNotNil(err) {
        logger.Error("error delete objects by prefix on bucket:", err)
    } else {
        logger.Info("prefix", prefix, "deleted successfully!")
    }
}
```

Output:

    [INFO 2024/01/12 10:52:58] main.go:32: prefix examples/ deleted successfully!

For more examples of prefix deletion, 
visit [link](https://github/GabrielHCataldo/go-cloud-storage/blob/main/_example/main).

Used go drives
------
- https://github.com/googleapis/google-cloud-go
- https://github.com/aws/aws-sdk-go

How to contribute
------
Make a pull request, or if you find a bug, open it
an Issues.

License
-------
Distributed under MIT license, see the license file within the code for more details.   
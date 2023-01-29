package backend

import (
    "context"
    "fmt"
    "io"
// 对文件做io
    // "around/constants"
    "around/util"

    "cloud.google.com/go/storage"
)

var (
    GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
    client *storage.Client
    bucket string
}

func InitGCSBackend(config *util.GCSInfo) {

    client, err := storage.NewClient(context.Background())
    if err != nil {
        panic(err)
    }

    GCSBackend = &GoogleCloudStorageBackend{
        client: client,
        bucket: config.Bucket,
    }
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
    ctx := context.Background()
    object := backend.client.Bucket(backend.bucket).Object(objectName)
    wc := object.NewWriter(ctx)
    if _, err := io.Copy(wc, r); err != nil {
        return "", err
    }
// 把文件流从前端读一下，存进去
    if err := wc.Close(); err != nil {
        return "", err
    }

    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        return "", err
		//object access control level所有user都可读权限
    }

    attrs, err := object.Attrs(ctx)
	// attributes
    if err != nil {
        return "", err
    }

    fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
	// 存好的url
}

package gcs

import (
	"context"
	"fmt"
	"io"

	consts "../constant"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func SaveToGCS(r io.Reader, bucketName, objectName string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(consts.CREDENTIAL_PATH))
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)
	if _, err := bucket.Attrs(ctx); err != nil {
		return nil, err
	}
	object := bucket.Object(objectName)
	w := object.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		fmt.Printf("can't close the gcs client %s\n", err)
		return nil, err
	}
	if err = object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, err
	}
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Image is saved to GCS %s\n", attrs.MediaLink)
	return attrs, nil

}

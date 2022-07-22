package utils

import (
	"fmt"

	"cloud.google.com/go/storage"
)

func ObjectURL(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}

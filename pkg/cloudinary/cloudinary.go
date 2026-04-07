package cloudinary

import (
	"context"
	"mime/multipart"
	"os"

	cld "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadImage uploads a multipart file to Cloudinary under the given folder
// and returns the secure HTTPS URL.
func UploadImage(ctx context.Context, file multipart.File, folder string) (string, error) {
	client, err := cld.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", err
	}

	result, err := client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

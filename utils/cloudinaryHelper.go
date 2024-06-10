package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func ConnectToCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Cloudinary: %w", err)
	}

	return cld, nil
}

func UploadToCloudinary(cld *cloudinary.Cloudinary, filePath string) (*uploader.UploadResult, error) {
	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{})
	if err != nil {
		return nil, fmt.Errorf("cannot upload image: %w", err)
	}

	return resp, nil
}

func DeleteFromCloudinary(cld *cloudinary.Cloudinary, publicID string) error {
	var ctx = context.Background()
	_, err := cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		PublicIDs: []string{publicID},
	})
	if err != nil {
		return fmt.Errorf("cannot delete image: %w", err)
	}
	return nil
}

package services

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tus/tusd/v2/pkg/azurestore"
	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/memorylocker"
	"log"
	"os"
)

var stdout = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
var stderr = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

func CreateComposer() (*handler.StoreComposer, error) {
	tusdAzureConfig, err := ReadAzureConfig()
	if err != nil {
		return nil, err
	}

	stdout.Printf("Using Azure endpoint %s.\n", tusdAzureConfig.Endpoint)

	Composer := handler.NewStoreComposer()

	azConfig := &azurestore.AzConfig{
		AccountName:         tusdAzureConfig.AccountName,
		AccountKey:          tusdAzureConfig.AccountKey,
		ContainerName:       tusdAzureConfig.ContainerName,
		ContainerAccessType: tusdAzureConfig.ContainerAccessType,
		BlobAccessTier:      tusdAzureConfig.BlobAccessTier,
		Endpoint:            tusdAzureConfig.Endpoint,
	}

	azService, err := azurestore.NewAzureService(azConfig)
	if err != nil {
		stderr.Fatalf(err.Error())
	}

	store := azurestore.New(azService)
	store.ObjectPrefix = tusdAzureConfig.ObjectPrefix
	store.Container = tusdAzureConfig.ContainerName
	store.UseIn(Composer)

	locker := memorylocker.New()
	locker.UseIn(Composer)

	return Composer, nil
}

func ReadAzureConfig() (*TusdAzureConfig, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("unable to load .env file: %v", err)
	}

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	if accountName == "" {
		return nil, fmt.Errorf("no service account name for Azure BlockBlob Storage using the AZURE_STORAGE_ACCOUNT environment variable")
	}

	accountKey := os.Getenv("AZURE_STORAGE_KEY")
	if accountKey == "" {
		return nil, fmt.Errorf("no service account key for Azure BlockBlob Storage using the AZURE_STORAGE_KEY environment variable")
	}

	azureStorage := os.Getenv("AZURE_STORAGE_CONTAINER")
	if azureStorage == "" {
		return nil, fmt.Errorf("no service account key for Azure BlockBlob Storage using the AZURE_STORAGE_CONTAINER environment variable")
	}

	azureEndpoint := os.Getenv("AZURE_ENDPOINT")
	if azureEndpoint == "" {
		azureEndpoint = fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	}

	objectPrefix := os.Getenv("AZURE_OBJECT_PREFIX")
	containerAccessType := os.Getenv("AZURE_CONTAINER_ACCESS_TYPE")
	blobAccessTier := os.Getenv("AZURE_BLOB_ACCESS_TIER")

	config := &TusdAzureConfig{
		AccountName:         accountName,
		AccountKey:          accountKey,
		ContainerName:       azureStorage,
		Endpoint:            azureEndpoint,
		ObjectPrefix:        objectPrefix,
		ContainerAccessType: containerAccessType,
		BlobAccessTier:      blobAccessTier,
	}

	return config, nil
}

type TusdAzureConfig struct {
	AccountName         string
	AccountKey          string
	ContainerName       string
	Endpoint            string
	ObjectPrefix        string
	ContainerAccessType string
	BlobAccessTier      string
}

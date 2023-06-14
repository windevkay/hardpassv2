package azure

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	// _ "github.com/joho/godotenv/autoload"
)

type AzureClient struct {
	Keys    *azkeys.Client
	Secrets *azsecrets.Client
}

func SetupClient() (*AzureClient, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	keysClient, err := azkeys.NewClient(os.Getenv("AZURE_KEYVAULT_URL"), credential, nil)
	if err != nil {
		return nil, err
	}

	secretsClient, err := azsecrets.NewClient(os.Getenv("AZURE_KEYVAULT_URL"), credential, nil)
	if err != nil {
		return nil, err
	}

	return &AzureClient{
		Keys:    keysClient,
		Secrets: secretsClient,
	}, nil
}

package azure

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"

	_ "github.com/joho/godotenv/autoload"
)

func SetupClient() (*azkeys.Client, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	client, err := azkeys.NewClient(os.Getenv("AZURE_KEYVAULT_URL"), credential, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

package azure

import (
	"testing"
)

// var client *azkeys.Client

// func setupClient(t *testing.T) {
// 	t.Setenv("AZURE_TENANT_ID", os.Getenv("AZURE_TENANT_ID"))
// 	t.Setenv("AZURE_CLIENT_ID", os.Getenv("AZURE_CLIENT_ID"))
// 	t.Setenv("AZURE_CLIENT_SECRET", os.Getenv("AZURE_CLIENT_SECRET"))
// 	t.Setenv("AZURE_KEYVAULT_URL", os.Getenv("AZURE_KEYVAULT_URL"))

// 	fmt.Printf("value ---------> %v\n", os.Getenv("AZURE_KEYVAULT_URL"))

// 	credential, err := azidentity.NewDefaultAzureCredential(nil)

// 	if err != nil {
// 		fmt.Print(err)
// 		t.Error(err)
// 	}

// 	cl, err := azkeys.NewClient(os.Getenv("AZURE_KEYVAULT_URL"), credential, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	client = cl
// }

func TestGenerateSecureKey(t *testing.T) {
	client, err := SetupClient()
	if err != nil {
		t.Error(err)
	}

	_, err = generateSecureKey(client, "test")
	if err != nil {
		t.Error(err)
	}
}

func TestGenPassword(t *testing.T) {
	client, err := SetupClient()
	if err != nil {
		t.Error(err)
	}

	_, err = GenPassword(client, "test", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecryptPassword(t *testing.T) {
	client, err := SetupClient()
	if err != nil {
		t.Error(err)
	}

	originalString := "random-password-string"
	password, err := GenPassword(client, "test", originalString)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := DecryptPassword(client, password.KeyIdentifier, password.KeyVersion, password.Text)
	if err != nil {
		t.Error(err)
	}

	if decrypted != originalString {
		t.Errorf("Decrypted password %v does not match original string %v", decrypted, originalString)
	}
}

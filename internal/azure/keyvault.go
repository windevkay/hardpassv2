package azure

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`"

type Password struct {
	Text          string
	KeyIdentifier string
}

func GenPassword(client *AzureClient, keyIdentifier string, password any) (*Password, error) {
	// generate secure key
	keystring, err := generateSecureKey(client.Keys, keyIdentifier)
	if err != nil {
		return nil, err
	}
	// extract version from keystring
	split := strings.Split(string(*keystring), "/")
	versionIndex := len(split) - 1
	keyVersion := split[versionIndex]

	// generate secure string or use provided password
	var text string
	if password != nil {
		text = password.(string)
	} else {
		text, err = genSecureString()
		if err != nil {
			return nil, err
		}
	}

	// encrypt secure string
	resp, err := client.Keys.Encrypt(context.TODO(), keyIdentifier, keyVersion, azkeys.KeyOperationsParameters{
		Algorithm: to.Ptr(azkeys.JSONWebKeyEncryptionAlgorithmRSAOAEP256),
		Value:     []byte(text),
	}, nil)

	if err != nil {
		return nil, err
	}

	// save key version as secret
	_, err = client.Secrets.SetSecret(context.TODO(), keyIdentifier, azsecrets.SetSecretParameters{
		Value: &keyVersion,
	}, nil)

	if err != nil {
		return nil, err
	}

	return &Password{Text: hex.EncodeToString(resp.Result), KeyIdentifier: keyIdentifier}, nil
}

func DecryptPassword(client *AzureClient, keyIdentifier string, text string) (string, error) {
	decodedString, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	// get key version from secret
	version, err := client.Secrets.GetSecret(context.TODO(), keyIdentifier, "", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Keys.Decrypt(context.TODO(), keyIdentifier, *version.Value, azkeys.KeyOperationsParameters{
		Algorithm: to.Ptr(azkeys.JSONWebKeyEncryptionAlgorithmRSAOAEP256),
		Value:     decodedString,
	}, nil)

	if err != nil {
		return "", err
	}

	return string(resp.Result), nil
}

// GenerateSecureKey generates a secure key using the Azure Key Vault SDK
func generateSecureKey(client *azkeys.Client, keyIdentifier string) (*azkeys.ID, error) {
	params := azkeys.CreateKeyParameters{
		Kty: to.Ptr(azkeys.JSONWebKeyTypeRSA),
	}
	// if a key with the same name already exists, a new version of that key is created
	resp, err := client.CreateKey(context.TODO(), keyIdentifier, params, nil)

	if err != nil {
		return nil, err
	}
	// *azkeys.ID is a struct with a field KID string
	return resp.Key.KID, nil
}

func genSecureString() (string, error) {
	charBytes := make([]byte, 64)
	_, err := rand.Read(charBytes)
	if err != nil {
		return "", err
	}

	for i, b := range charBytes {
		random := int(b) % len(charset)
		charBytes[i] = charset[random]
	}

	return string(charBytes), nil
}

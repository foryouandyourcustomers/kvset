// https://blog.abhi.host/blog/2019/08/17/fetch-certificates-from-keyvault-in-go/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/prometheus/common/log"
)

var (
	keyvaultAuthorizer autorest.Authorizer
	kvClient           azureKvClient

	vaultName   string
	secretName  string
	secretValue string
)

type azureKvClient struct {
	ctx   context.Context
	vault string

	client keyvault.BaseClient

	authenticated bool
	vaultBaseURL  string
}

func (akv *azureKvClient) authenticate(v string) {
	akv.vault = v

	// lets try to get authorizer first from cli
	// then from environment
	// this allows us to run the cli locally without
	// any additional work while also giving us the
	// possibility to use it in a CI with env vars set

	a, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		// looking at the newauthorizerfromenvrionment funciton it
		// seems that thing never returns an error whatsoever!
		log.Debug("Unable to create authorizer from az cli. Lets load the authorizer from the environment and hope for the best!")
		a, _ = auth.NewAuthorizerFromEnvironment()

	}

	akv.client.Authorizer = a
	akv.authenticated = true

	akv.vaultBaseURL = fmt.Sprintf("https://%s.%s", akv.vault, azure.PublicCloud.KeyVaultDNSSuffix)
}

func (akv *azureKvClient) setSecret(s string, a string) error {
	p := &keyvault.SecretSetParameters{Value: &a}
	_, err := akv.client.SetSecret(akv.ctx, akv.vaultBaseURL, s, *p)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	flag.StringVar(&vaultName, "v", "", "Name of the keyvault (env var: VAULT)")
	flag.StringVar(&secretName, "s", "", "Name of the secret to set (env var: SECRET)")
	flag.StringVar(&secretValue, "a", "", "Value of the secret (env var: VALUE)")
	flag.Parse()

	if os.Getenv("SECRET") != "" {
		secretName = os.Getenv("SECRET")
	}
	if os.Getenv("VAULT") != "" {
		vaultName = os.Getenv("VAULT")
	}
	if os.Getenv("secretValue") != "" {
		secretValue = os.Getenv("secretValue")
	}

	if (vaultName == "") || (secretName == "") || (secretValue == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {

	kvClient.ctx = context.Background()
	kvClient.authenticate(vaultName)

	err := kvClient.setSecret(secretName, secretValue)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("secret '%s' created or updated", secretName)
}

package vault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

func Test_getSecret(t *testing.T) {
	err := os.Setenv("VAULT_ADDR", "https://vault.example.com")
	if err != nil {
		return
	}
	err = os.Setenv("VAULT_TOKEN", "my-test-token")
	if err != nil {
		return
	}

	// Create a test server that will simulate the Vault API
	// and provide a canned response to the getSecret function
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected request method to be GET, got '%s'", r.Method)
		}

		// Assert that the request URL is correct
		if r.URL.String() != "/v1/secret/data/my-secret" {
			t.Errorf("Expected request URL to be '/v1/secret/data/my-secret', got '%s'", r.URL.String())
		}

		// Assert that the request has the correct authorization header
		if r.Header.Get("X-Vault-Token") != "my-test-token" {
			t.Errorf("Expected request to have authorization header 'my-test-token', got '%s'", r.Header.Get("X-Vault-Token"))
		}

		secret := vault.KVSecret{
			Data: map[string]interface{}{
				"username": "my-username",
				"password": "my-password",
			},
		}
		// Write the secret data to the response with error handling
		_, err := fmt.Fprintln(w, secret)
		if err != nil {
			t.Errorf("Error writing secret to test server: %v", err)
			return
		}
	}))
	defer ts.Close()

	// os.Setenv("VAULT_ADDR", ts.URL)
	// ctx, client := createClientContext()
	// secret := getSecret(ctx, client, "my-secret")

	// if secret["username"] != "my-username" {
	//  t.Errorf("Expected secret username to be 'my-username', got '%s'", secret["username"])
	// }
	// if secret["password"] != "my-password" {
	//  t.Errorf("Expected secret password to be 'my-password', got '%s'", secret["password"])
	// }

}

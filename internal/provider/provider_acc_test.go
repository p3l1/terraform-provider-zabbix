// ABOUTME: Shared test setup for acceptance tests.
// ABOUTME: Provides provider factories and pre-check functions.

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "zabbix" {}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"zabbix": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	if os.Getenv("ZABBIX_URL") == "" {
		os.Setenv("ZABBIX_URL", "http://127.0.0.1:8080/api_jsonrpc.php")
	}

	if os.Getenv("ZABBIX_API_TOKEN") == "" {
		os.Setenv("ZABBIX_API_TOKEN", "071fb9d2e8f72cf9c40128f0f5aab3def1bab0893413314b083fdcb4551eb01a")
	}
}

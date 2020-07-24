package limelight

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"limelight": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("error validating provider: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	var requiredVariables = []string{"LLNW_API_USERNAME", "LLNW_API_KEY", "LLNW_TEST_SHORTNAME"}
	for _, element := range requiredVariables {
		if v := os.Getenv(element); v == "" {
			str := fmt.Sprintf("%s must be set for acceptance tests", element)
			t.Fatal(str)
		}
	}
}

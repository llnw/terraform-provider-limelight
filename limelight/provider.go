package limelight

import (
	"fmt"

	"github.com/llnw/llnw-sdk-go/configuration"
	"github.com/llnw/llnw-sdk-go/edgefunctions"
	"github.com/llnw/terraform-provider-limelight/version"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"config_api_base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LLNW_CONFIG_API_URL", nil),
				Description: "The base URL for the Limelight Networks Configuration API (trailing / should be omitted)",
			},
			"edgefunctions_api_base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LLNW_EDGEFUNCTIONS_API_URL", nil),
				Description: "The base URL for the Limelight Networks EdgeFunctions API (trailing / should be omitted)",
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("LLNW_API_USERNAME", nil),
				Description:  "The username to be used for authenticating with the Limelight Networks Configuration API",
				ValidateFunc: validation.NoZeroValues,
			},
			"api_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("LLNW_API_KEY", nil),
				Description:  "The API key to be used for authenticating with the Limelight Networks Configuration API",
				ValidateFunc: validation.NoZeroValues,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"limelight_ip_ranges": dataSourceLimelightIPRanges(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"limelight_delivery":           resourceLimelightDelivery(),
			"limelight_edgefunction":       resourceLimelightEdgeFunction(),
			"limelight_edgefunction_alias": resourceLimelightEdgeFunctionAlias(),
			// TODO: enable in RTS v2
			//"limelight_realtime_streaming_slot": resourceLimelightRealtimeStreamingSlot(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	configBaseURL := d.Get("config_api_base_url").(string)
	edgefunctionsBaseURL := d.Get("edgefunctions_api_base_url").(string)
	username := d.Get("username").(string)
	apiKey := d.Get("api_key").(string)

	m := make(map[string]interface{})

	var configurationClient *configuration.ConfigurationClient
	var edgeFunctionsClient *edgefunctions.EdgeFunctionsClient

	if len(configBaseURL) == 0 {
		configurationClient = configuration.NewClient(username, apiKey)
	} else {
		configurationClient = configuration.NewClientOverrideBaseUrl(username, apiKey, configBaseURL)
	}

	if len(edgefunctionsBaseURL) == 0 {
		edgeFunctionsClient = edgefunctions.NewClient(username, apiKey)
	} else {
		edgeFunctionsClient = edgefunctions.NewClientOverrideBaseUrl(username, apiKey, edgefunctionsBaseURL)
	}

	terraformUserAgent := httpclient.TerraformUserAgent(terraformVersion)
	providerUserAgent := fmt.Sprintf("terraform-provider-limelight/%s", version.ProviderVersion)
	userAgent := fmt.Sprintf("%s %s", terraformUserAgent, providerUserAgent)
	configurationClient.SetUserAgent(userAgent)
	edgeFunctionsClient.SetUserAgent(userAgent)

	m["config"] = configurationClient
	m["edgefunctions"] = edgeFunctionsClient

	return m, nil
}

func getConfigurationClient(m interface{}) *configuration.ConfigurationClient {
	clients := m.(map[string]interface{})
	return clients["config"].(*configuration.ConfigurationClient)
}

func getEdgeFunctionsClient(m interface{}) *edgefunctions.EdgeFunctionsClient {
	clients := m.(map[string]interface{})
	return clients["edgefunctions"].(*edgefunctions.EdgeFunctionsClient)
}

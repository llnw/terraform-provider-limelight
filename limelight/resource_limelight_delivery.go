package limelight

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/llnw/llnw-sdk-go/configuration"
)

func resourceLimelightDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceLimelightDeliveryCreate,
		Read:   resourceLimelightDeliveryRead,
		Update: resourceLimelightDeliveryUpdate,
		Delete: resourceLimelightDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"shortname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"service_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LLNW-Generic",
			},
			"protocol_set": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"published_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
							}, false),
						},
						"source_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http",
								"https",
							}, false),
						},
						"source_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"option": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"parameters": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"published_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"published_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version_number": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceLimelightDeliveryCreate(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	shortname := d.Get("shortname").(string)
	serviceProfile := d.Get("service_profile").(string)
	protocolSets := expandProtocolSets(d.Get("protocol_set").([]interface{}), c, shortname, serviceProfile)
	publishedHostname := d.Get("published_hostname").(string)
	publishedPath := d.Get("published_path").(string)
	sourceHostname := d.Get("source_hostname").(string)
	sourcePath := d.Get("source_path").(string)

	body := &configuration.DeliveryServiceInstanceBody{
		ServiceProfileName: serviceProfile,
		ProtocolSets:       protocolSets,
		PublishedHostname:  publishedHostname,
		PublishedURLPath:   publishedPath,
		SourceHostname:     sourceHostname,
		SourceURLPath:      sourcePath,
		ServiceKey: configuration.ServiceKey{
			Name: "delivery",
		},
	}

	log.Printf("[INFO] Creating delivery configuration for service profile: %s", serviceProfile)
	deliveryServiceInstance, _, err := c.CreateDeliveryServiceInstance(body, shortname)

	if err != nil {
		return fmt.Errorf("error creating delivery configuration: %s", err)
	}

	d.SetId(deliveryServiceInstance.UUID)

	return resourceLimelightDeliveryRead(d, m)
}

func resourceLimelightDeliveryRead(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	log.Printf("[INFO] Fetching delivery configuration: %s", d.Id())
	deliveryServiceInstance, resp, err := c.GetDeliveryServiceInstance(d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Delivery configuration %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading delivery configuration: %s", err)
	}

	d.Set("shortname", deliveryServiceInstance.Shortname)
	d.Set("service_profile", deliveryServiceInstance.Body.ServiceProfileName)
	d.Set("published_hostname", deliveryServiceInstance.Body.PublishedHostname)
	d.Set("published_path", deliveryServiceInstance.Body.PublishedURLPath)
	d.Set("source_hostname", deliveryServiceInstance.Body.SourceHostname)
	d.Set("source_path", deliveryServiceInstance.Body.SourceURLPath)
	d.Set("version_number", deliveryServiceInstance.Revision.VersionNumber)
	d.Set("protocol_set", flattenProtocolSets(deliveryServiceInstance.Body.ProtocolSets))

	return nil
}

func resourceLimelightDeliveryUpdate(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	shortname := d.Get("shortname").(string)
	serviceProfile := d.Get("service_profile").(string)
	protocolSets := expandProtocolSets(d.Get("protocol_set").([]interface{}), c, shortname, serviceProfile)
	publishedHostname := d.Get("published_hostname").(string)
	publishedPath := d.Get("published_path").(string)
	sourceHostname := d.Get("source_hostname").(string)
	sourcePath := d.Get("source_path").(string)

	body := &configuration.DeliveryServiceInstanceBody{
		ServiceProfileName: serviceProfile,
		ProtocolSets:       protocolSets,
		PublishedHostname:  publishedHostname,
		PublishedURLPath:   publishedPath,
		SourceHostname:     sourceHostname,
		SourceURLPath:      sourcePath,
		ServiceKey: configuration.ServiceKey{
			Name: "delivery",
		},
	}

	log.Printf("[INFO] Updating delivery configuration for: %s", d.Id())
	_, _, err := c.UpdateDeliveryServiceInstance(d.Id(), body, shortname)

	if err != nil {
		return fmt.Errorf("error updating delivery configuration: %s", err)
	}

	return resourceLimelightDeliveryRead(d, m)
}

func resourceLimelightDeliveryDelete(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	log.Printf("[INFO] Deleting delivery configuration: %s", d.Id())
	_, _, err := c.DeleteDeliveryServiceInstance(d.Id())

	if err != nil {
		return fmt.Errorf("error deleting delivery configuration: %s", err)
	}

	return nil
}

func flattenProtocolSets(expandedProtocolSets []configuration.ProtocolSet) []map[string]interface{} {
	flattenedProtocolSets := make([]map[string]interface{}, len(expandedProtocolSets), len(expandedProtocolSets))

	for i, v := range expandedProtocolSets {
		m := make(map[string]interface{})
		m["published_protocol"] = v.PublishedProtocol
		m["source_protocol"] = v.SourceProtocol
		m["source_port"] = v.SourcePort
		m["option"] = flattenOptions(v.Options)
		flattenedProtocolSets[i] = m
	}

	return flattenedProtocolSets
}

func flattenOptions(expandedOptions []configuration.Option) []map[string]interface{} {
	flattenedOptions := make([]map[string]interface{}, len(expandedOptions), len(expandedOptions))

	for i, v := range expandedOptions {
		m := make(map[string]interface{})
		m["name"] = v.Name
		params := make([]string, len(v.Parameters), len(v.Parameters))
		for j, p := range v.Parameters {
			params[j] = fmt.Sprintf("%v", p)
		}
		m["parameters"] = params
		flattenedOptions[i] = m
	}

	return flattenedOptions
}

func expandProtocolSets(flattenedProtocolSets []interface{}, c *configuration.ConfigurationClient, shortname string, serviceProfile string) []configuration.ProtocolSet {
	expandedProtocolSets := make([]configuration.ProtocolSet, len(flattenedProtocolSets), len(flattenedProtocolSets))

	for i, v := range flattenedProtocolSets {
		rawProtocolSet := v.(map[string]interface{})
		protocolSet := configuration.ProtocolSet{
			PublishedProtocol: rawProtocolSet["published_protocol"].(string),
			SourceProtocol:    rawProtocolSet["source_protocol"].(string),
			Options:           expandOptions(rawProtocolSet["option"].([]interface{}), c, shortname, serviceProfile),
		}

		sourcePort := rawProtocolSet["source_port"].(int)
		if sourcePort != 0 {
			protocolSet.SourcePort = &sourcePort
		}

		expandedProtocolSets[i] = protocolSet
	}

	return expandedProtocolSets
}

func expandOptions(flattenedOptions []interface{}, c *configuration.ConfigurationClient, shortname string, serviceProfile string) []configuration.Option {
	expandedOptions := make([]configuration.Option, len(flattenedOptions), len(flattenedOptions))

	for i, v := range flattenedOptions {
		rawOption := v.(map[string]interface{})
		opt := configuration.Option{
			Name:       rawOption["name"].(string),
			Parameters: expandOptionParameters(rawOption["parameters"].([]interface{}), rawOption["name"].(string), c, shortname, serviceProfile),
		}
		expandedOptions[i] = opt
	}

	return expandedOptions
}

func expandOptionParameters(flattenedOptionParams []interface{}, optionName string, c *configuration.ConfigurationClient, shortname string, serviceProfile string) []interface{} {
	expandedOptionParams := make([]interface{}, len(flattenedOptionParams), len(flattenedOptionParams))

	for i, v := range flattenedOptionParams {
		stringVal := v.(string)

		argInt, _ := c.IsOptionArgumentInteger(shortname, serviceProfile, optionName, i)
		if argInt {
			if intVal, err := strconv.Atoi(v.(string)); err == nil {
				expandedOptionParams[i] = intVal
			} else {
				expandedOptionParams[i] = stringVal
			}
		} else {
			expandedOptionParams[i] = stringVal
		}
	}

	return expandedOptionParams
}

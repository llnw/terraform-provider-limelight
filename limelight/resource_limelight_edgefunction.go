package limelight

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/llnw/llnw-sdk-go/edgefunctions"
)

func resourceLimelightEdgeFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceLimelightEdgeFunctionCreate,
		Read:   resourceLimelightEdgeFunctionRead,
		Update: resourceLimelightEdgeFunctionUpdate,
		Delete: resourceLimelightEdgeFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"shortname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"function_archive": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"handler": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"runtime": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  256,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000,
			},
			"can_debug": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"environment_variable": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     envVarsElem,
			},
			"function_sha256": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"revision_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"reserved_concurrency": &schema.Schema{
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceLimelightEdgeFunctionCreate(d *schema.ResourceData, m interface{}) error {
	c := getEdgeFunctionsClient(m)

	shortname := d.Get("shortname").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	functionArchive := d.Get("function_archive").(string)
	handler := d.Get("handler").(string)
	runtime := d.Get("runtime").(string)
	memory := d.Get("memory").(int)
	timeout := d.Get("timeout").(int)
	canDebug := d.Get("can_debug").(bool)
	environmentVariables := expandEnvVars(d.Get("environment_variable").(*schema.Set))
	concurrency := d.Get("reserved_concurrency").(int)

	zipFile, err := loadZipFile(functionArchive)

	if err != nil {
		return err
	}

	edgeFunction := &edgefunctions.EdgeFunction{
		Name:                 name,
		Description:          description,
		FunctionArchive:      zipFile,
		Handler:              handler,
		Runtime:              runtime,
		Memory:               memory,
		Timeout:              timeout,
		CanDebug:             canDebug,
		EnvironmentVariables: environmentVariables,
	}

	log.Printf("[INFO] Creating EdgeFunction: %s", name)
	edgeFunctionResponse, _, err := c.CreateEdgeFunction(shortname, edgeFunction)

	if err != nil {
		return fmt.Errorf("error creating EdgeFunction: %s", err)
	}

	if concurrency > 0 {
		log.Printf("[INFO] Setting EdgeFunction concurrency to: %v", concurrency)
		_, err := c.SetEdgeFunctionConcurrency(name, shortname, concurrency)

		if err != nil {
			log.Printf("[WARN] Failed to set EdgeFunction currency for %s, rolling back EdgeFunction creation", name)
			_, delErr := c.DeleteEdgeFunction(name, shortname)
			if delErr != nil {
				log.Printf("[ERROR] Failed to delete EdgeFunction %s due to: %s", name, delErr)
			}
			return fmt.Errorf("failed to set EdgeFunction currency for %s due to: %s", name, err)
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", shortname, edgeFunctionResponse.Name))

	return resourceLimelightEdgeFunctionRead(d, m)
}

func resourceLimelightEdgeFunctionRead(d *schema.ResourceData, m interface{}) error {
	c := getEdgeFunctionsClient(m)

	shortname, name, err := resourceLimelightEdgeFunctionSplitID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Fetching EdgeFunction: %s", name)
	edgeFunction, resp, err := c.GetEdgeFunction(name, shortname)

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] EdgeFunction %s was not found", name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading EdgeFunction: %s", err)
	}

	d.Set("shortname", shortname)
	d.Set("name", name)
	d.Set("description", edgeFunction.Description)
	d.Set("handler", edgeFunction.Handler)
	d.Set("runtime", edgeFunction.Runtime)
	d.Set("memory", edgeFunction.Memory)
	d.Set("timeout", edgeFunction.Timeout)
	d.Set("can_debug", edgeFunction.CanDebug)
	d.Set("function_sha256", edgeFunction.Sha256)
	d.Set("environment_variable", flattenEnvVars(edgeFunction.EnvironmentVariables))
	d.Set("revision_id", edgeFunction.RevisionID)
	d.Set("reserved_concurrency", edgeFunction.ReservedConcurrency)

	return nil
}

func resourceLimelightEdgeFunctionUpdate(d *schema.ResourceData, m interface{}) error {
	c := getEdgeFunctionsClient(m)

	shortname, name, err := resourceLimelightEdgeFunctionSplitID(d.Id())

	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("function_sha256") {
		functionArchive := d.Get("function_archive").(string)
		zipFile, zipErr := loadZipFile(functionArchive)

		if zipErr != nil {
			return zipErr
		}

		log.Printf("[INFO] Updating EdgeFunction code for: %s", name)
		_, _, err := c.UpdateEdgeFunctionCode(name, shortname, zipFile)

		if err != nil {
			return fmt.Errorf("error updating EdgeFunction code: %s", err)
		}
		d.SetPartial("function_archive")
		d.SetPartial("function_sha256")
	}

	if d.HasChange("description") || d.HasChange("handler") || d.HasChange("runtime") || d.HasChange("memory") || d.HasChange("timeout") || d.HasChange("can_debug") || d.HasChange("environment_variable") {
		description := d.Get("description").(string)
		handler := d.Get("handler").(string)
		runtime := d.Get("runtime").(string)
		memory := d.Get("memory").(int)
		timeout := d.Get("timeout").(int)
		canDebug := d.Get("can_debug").(bool)
		environmentVariables := expandEnvVars(d.Get("environment_variable").(*schema.Set))

		edgeFunction := &edgefunctions.EdgeFunction{
			Description:          description,
			Handler:              handler,
			Runtime:              runtime,
			Memory:               memory,
			Timeout:              timeout,
			CanDebug:             canDebug,
			EnvironmentVariables: environmentVariables,
		}

		log.Printf("[INFO] Updating EdgeFunction configuration for: %s", name)
		_, _, err := c.UpdateEdgeFunctionConfiguration(name, shortname, edgeFunction)

		if err != nil {
			return fmt.Errorf("error updating EdgeFunction configuration: %s", err)
		}
	}

	if d.HasChange("reserved_concurrency") {
		concurrency := d.Get("reserved_concurrency").(int)
		log.Printf("[INFO] Updating EdgeFunction reserved concurrency to: %d", concurrency)
		_, err := c.SetEdgeFunctionConcurrency(name, shortname, concurrency)
		if err != nil {
			return fmt.Errorf("error updating EdgeFunction concurrency due to: %s", err)
		}
	}

	d.Partial(false)

	return resourceLimelightEdgeFunctionRead(d, m)
}

func resourceLimelightEdgeFunctionDelete(d *schema.ResourceData, m interface{}) error {
	c := getEdgeFunctionsClient(m)

	shortname, name, err := resourceLimelightEdgeFunctionSplitID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting EdgeFunction: %s", name)
	_, err = c.DeleteEdgeFunction(name, shortname)

	if err != nil {
		return fmt.Errorf("error deleting EdgeFunction: %s", err)
	}

	return nil
}

func loadZipFile(path string) ([]byte, error) {
	zipFile, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func resourceLimelightEdgeFunctionSplitID(id string) (string, string, error) {
	shortName, fnName, err := splitSeparatedPair(id, ":")

	if err != nil {
		return "", "", fmt.Errorf("EdgeFunction ID in unexpected format (expected '<shortname>:<function name>'): %s", id)
	}

	return shortName, fnName, nil
}

var envVarsElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func flattenEnvVars(expandedEnvVars []edgefunctions.EnvironmentVariable) *schema.Set {
	flattenedEnvVars := make([]interface{}, len(expandedEnvVars), len(expandedEnvVars))

	for i, v := range expandedEnvVars {
		m := make(map[string]interface{})
		m["name"] = v.Name
		m["value"] = v.Value
		flattenedEnvVars[i] = m
	}

	return schema.NewSet(schema.HashResource(envVarsElem), flattenedEnvVars)
}

func expandEnvVars(flattenedEnvVars *schema.Set) []edgefunctions.EnvironmentVariable {
	expandedEnvVars := make([]edgefunctions.EnvironmentVariable, flattenedEnvVars.Len(), flattenedEnvVars.Len())

	for i, v := range flattenedEnvVars.List() {
		rawEnvVar := v.(map[string]interface{})
		opt := edgefunctions.EnvironmentVariable{
			Name:  rawEnvVar["name"].(string),
			Value: rawEnvVar["value"].(string),
		}
		expandedEnvVars[i] = opt
	}

	return expandedEnvVars
}

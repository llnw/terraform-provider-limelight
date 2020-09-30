package limelight

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/llnw/llnw-sdk-go/edgefunctions"
)

func resourceLimelightEdgeFunctionAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceLimelightEdgeFunctionAliasCreate,
		Read:   resourceLimelightEdgeFunctionAliasRead,
		Update: resourceLimelightEdgeFunctionAliasUpdate,
		Delete: resourceLimelightEdgeFunctionAliasDelete,
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
			"function_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"function_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"revision_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceLimelightEdgeFunctionAliasCreate(d *schema.ResourceData, m interface{}) error {
	client := getEdgeFunctionsClient(m)

	shortname := d.Get("shortname").(string)
	name := d.Get("name").(string)
	fnName := d.Get("function_name").(string)
	fnVersion := d.Get("function_version").(string)
	description := d.Get("description").(string)

	alias := &edgefunctions.EdgeFunctionAlias{
		Name:            name,
		Description:     description,
		FunctionVersion: fnVersion,
	}

	log.Printf("[INFO] Creating Alias %s for EdgeFunction %s", name, fnName)
	_, _, err := client.CreateEdgeFunctionAlias(fnName, shortname, alias)

	if err != nil {
		return fmt.Errorf("error creating EdgeFunction Alias %s: %v", name, err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", shortname, fnName, name))

	return resourceLimelightEdgeFunctionAliasRead(d, m)
}

func resourceLimelightEdgeFunctionAliasRead(d *schema.ResourceData, m interface{}) error {
	client := getEdgeFunctionsClient(m)

	shortname, fnName, aliasName, err := resourceLimelightEdgeFunctionAliasSplitID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Fetching Alias %s for EdgeFunction %s", aliasName, fnName)
	aliasResponse, resp, err := client.GetEdgeFunctionAlias(fnName, shortname, aliasName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Alias %s for EdgeFunction %s was not found", aliasName, fnName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading EdgeFunction Alias: %s", err)
	}

	d.Set("shortname", shortname)
	d.Set("name", aliasResponse.Name)
	d.Set("description", aliasResponse.Description)
	d.Set("function_name", aliasResponse.Function)
	d.Set("function_version", aliasResponse.FunctionVersion)
	d.Set("revision_id", aliasResponse.RevisionID)

	return nil
}

func resourceLimelightEdgeFunctionAliasUpdate(d *schema.ResourceData, m interface{}) error {
	client := getEdgeFunctionsClient(m)

	shortname := d.Get("shortname").(string)
	aliasName := d.Get("name").(string)
	fnName := d.Get("function_name").(string)
	fnVersion := d.Get("function_version").(string)
	description := d.Get("description").(string)
	revisionID := d.Get("revision_id").(int)

	alias := &edgefunctions.EdgeFunctionAlias{
		Description:     description,
		FunctionVersion: fnVersion,
		RevisionID:      revisionID,
	}

	log.Printf("[INFO] Updating Alias %s for EdgeFunction %s", aliasName, fnName)
	_, _, err := client.UpdateEdgeFunctionAlias(fnName, shortname, aliasName, alias)
	if err != nil {
		return fmt.Errorf("error updating Edge Function Alias %s: %v", aliasName, err)
	}

	return resourceLimelightEdgeFunctionAliasRead(d, m)
}

func resourceLimelightEdgeFunctionAliasDelete(d *schema.ResourceData, m interface{}) error {
	client := getEdgeFunctionsClient(m)

	shortname, fnName, aliasName, err := resourceLimelightEdgeFunctionAliasSplitID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Alias %s for EdgeFunction %s", aliasName, fnName)
	_, err = client.DeleteEdgeFunctionAlias(fnName, shortname, aliasName)
	if err != nil {
		return fmt.Errorf("failed to delete Edge Function Alias %s due to: %v", aliasName, err)
	}

	return nil
}

func resourceLimelightEdgeFunctionAliasSplitID(id string) (string, string, string, error) {
	shortName, fnName, aliasName, err := splitSeparatedTriple(id, ":")

	if err != nil {
		return "", "", "", fmt.Errorf("EdgeFunction Alias ID in unexpected format (expected '<shortname>:<function name>:<alias name>'): %s", id)
	}

	return shortName, fnName, aliasName, nil
}

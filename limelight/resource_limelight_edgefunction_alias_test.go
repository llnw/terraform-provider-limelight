package limelight

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceLimelightEdgeFunctionAlias_minimal(t *testing.T) {
	testResourceName := "limelight_edgefunction_alias.test_alias"
	fnName := "terraform_efa_test_minimal_fn"
	aliasName := "terraform_efa_minimal_alias"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightEdgeFunctionAliasCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightEdgeFunctionAliasMinimalTemplate(fnName, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightEdgeFunctionAliasExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", aliasName),
					resource.TestCheckResourceAttrSet(testResourceName, "function_name"),
					resource.TestCheckResourceAttr(testResourceName, "revision_id", "0"),
					resource.TestCheckResourceAttr(testResourceName, "function_version", "$LATEST"),
				),
			},
		},
	})
}

func TestAccResourceLimelightEdgeFunctionAlias_update(t *testing.T) {
	testResourceName := "limelight_edgefunction_alias.test_alias"
	fnName := "terraform_efa_test_update_fn"
	aliasName := "terraform_efa_update_alias"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightEdgeFunctionAliasCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightEdgeFunctionAliasBasicTemplate(fnName, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightEdgeFunctionAliasExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", aliasName),
					resource.TestCheckResourceAttrSet(testResourceName, "function_name"),
					resource.TestCheckResourceAttr(testResourceName, "revision_id", "0"),
					resource.TestCheckResourceAttr(testResourceName, "function_version", "$LATEST"),
					resource.TestCheckResourceAttr(testResourceName, "description", "Alias description"),
				),
			},
			{
				Config: testAccLimelightEdgeFunctionAliasBasicTemplate_update(fnName, aliasName),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightEdgeFunctionAliasExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", aliasName),
					resource.TestCheckResourceAttrSet(testResourceName, "function_name"),
					resource.TestCheckResourceAttr(testResourceName, "revision_id", "1"),
					resource.TestCheckResourceAttr(testResourceName, "function_version", "$LATEST"),
					resource.TestCheckResourceAttr(testResourceName, "description", "Alias description updated"),
				),
			},
		},
	})
}

func TestAccResourceLimelightEdgeFunctionAlias_importBasic(t *testing.T) {
	testResourceName := "limelight_edgefunction_alias.test_alias"
	fnName := "terraform_efa_test_import_fn"
	aliasName := "terraform_efa_import_alias"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightEdgeFunctionAliasCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightEdgeFunctionAliasMinimalTemplate(fnName, aliasName),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLimelightEdgeFunctionAliasCheckDestroy(state *terraform.State) error {
	client := getEdgeFunctionsClient(testAccProvider.Meta().(map[string]interface{}))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "limelight_edgefunction_alias" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		shortname, fnName, aliasName, err := splitSeparatedTriple(resourceID, ":")

		if err != nil {
			return err
		}

		_, resp, err := client.GetEdgeFunctionAlias(fnName, shortname, aliasName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("error retrieving EdgeFunction Alias with ID %s. Error: %v", resourceID, err)
		}

	}
	return fmt.Errorf("EdgeFunction Alias still exists")
}

func testAccLimelightEdgeFunctionAliasExists(testResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		client := getEdgeFunctionsClient(testAccProvider.Meta().(map[string]interface{}))

		rs, ok := state.RootModule().Resources[testResourceName]
		if !ok {
			return fmt.Errorf("EdgeFunction Alias not found in resources")
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("EdgeFunction Alias ID not set in resources ")
		}

		shortname, fnName, aliasName, err := splitSeparatedTriple(resourceID, ":")

		if err != nil {
			return err
		}

		edgeFunctionAlias, resp, err := client.GetEdgeFunctionAlias(fnName, shortname, aliasName)
		if err != nil {
			return fmt.Errorf("error retrieving EdgeFunction Alias: %s", resourceID)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error while checking if EdgeFunction Alias %s exists. HTTP return code: %d", resourceID, resp.StatusCode)
		}

		if edgeFunctionAlias.Name == aliasName {
			return nil
		}

		return fmt.Errorf("EdgeFunction Alias wasn't found")
	}
}

func testAccLimelightEdgeFunctionAliasMinimalTemplate(fnName, aliasName string) string {
	return fmt.Sprintf(`
locals {
	fn_sha = filesha256("testdata/edgefunc/py_function.zip")
}

resource "limelight_edgefunction" "test_ef" {
	shortname        = "%s"
	name             = "%s"
	function_archive = "testdata/edgefunc/py_function.zip"
	handler          = "hello_world.handler"
	runtime          = "python3"
	function_sha256  = local.fn_sha
}

resource "limelight_edgefunction_alias" "test_alias" {
    shortname        = "%s"
    name             = "%s"
    function_name    = limelight_edgefunction.test_ef.name
    function_version = "$LATEST"
}`, getShortname(), fnName, getShortname(), aliasName)
}

func testAccLimelightEdgeFunctionAliasBasicTemplate(fnName, aliasName string) string {
	return fmt.Sprintf(`
locals {
	fn_sha = filesha256("testdata/edgefunc/py_function.zip")
}

resource "limelight_edgefunction" "test_ef" {
	shortname        = "%s"
	name             = "%s"
	function_archive = "testdata/edgefunc/py_function.zip"
	handler          = "hello_world.handler"
	runtime          = "python3"
	function_sha256  = local.fn_sha
}

resource "limelight_edgefunction_alias" "test_alias" {
    shortname        = "%s"
    name             = "%s"
    function_name    = limelight_edgefunction.test_ef.name
	function_version = "$LATEST"
	description      = "Alias description"
}`, getShortname(), fnName, getShortname(), aliasName)
}

func testAccLimelightEdgeFunctionAliasBasicTemplate_update(fnName, aliasName string) string {
	return fmt.Sprintf(`
locals {
	fn_sha = filesha256("testdata/edgefunc/py_function.zip")
}

resource "limelight_edgefunction" "test_ef" {
	shortname        = "%s"
	name             = "%s"
	function_archive = "testdata/edgefunc/py_function.zip"
	handler          = "hello_world.handler"
	runtime          = "python3"
	function_sha256  = local.fn_sha
}

resource "limelight_edgefunction_alias" "test_alias" {
    shortname        = "%s"
    name             = "%s"
    function_name    = limelight_edgefunction.test_ef.name
	function_version = "$LATEST"
	description      = "Alias description updated"
}`, getShortname(), fnName, getShortname(), aliasName)
}

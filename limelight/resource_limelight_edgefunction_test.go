package limelight

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceLimelightEdgeFunction_minimal(t *testing.T) {
	testResourceName := "limelight_edgefunction.test_ef"
	fnName := "terraform_ef_test_minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightEdgeFunctionCheckDestroy(state, fnName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightEdgeFunctionMinimalTemplate(fnName),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightEdgeFunctionExists(testResourceName, fnName),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", fnName),
					resource.TestCheckResourceAttr(testResourceName, "function_archive", "testdata/edgefunc/py_function.zip"),
					resource.TestCheckResourceAttr(testResourceName, "handler", "hello_world.handler"),
					resource.TestCheckResourceAttr(testResourceName, "runtime", "python3"),
					resource.TestCheckResourceAttr(testResourceName, "memory", "256"),
					resource.TestCheckResourceAttr(testResourceName, "timeout", "5000"),
					resource.TestCheckResourceAttr(testResourceName, "can_debug", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "function_sha256"),
					resource.TestCheckResourceAttr(testResourceName, "revision_id", "0"),
					resource.TestCheckResourceAttr(testResourceName, "reserved_concurrency", "0"),
					resource.TestCheckResourceAttr(testResourceName, "environment_variable.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceLimelightEdgeFunction_update(t *testing.T) {
	testResourceName := "limelight_edgefunction.test_ef"
	fnName := "terraform_ef_test_update"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightEdgeFunctionCheckDestroy(state, fnName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightEdgeFunctionBasicTemplate(fnName),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightEdgeFunctionExists(testResourceName, fnName),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", fnName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Test function"),
					resource.TestCheckResourceAttr(testResourceName, "function_archive", "testdata/edgefunc/py_function.zip"),
					resource.TestCheckResourceAttr(testResourceName, "handler", "hello_world.handler"),
					resource.TestCheckResourceAttr(testResourceName, "runtime", "python3"),
					resource.TestCheckResourceAttr(testResourceName, "memory", "512"),
					resource.TestCheckResourceAttr(testResourceName, "timeout", "6000"),
					resource.TestCheckResourceAttr(testResourceName, "can_debug", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "function_sha256"),
					resource.TestCheckResourceAttr(testResourceName, "revision_id", "0"),
					resource.TestCheckResourceAttr(testResourceName, "reserved_concurrency", "3"),
					resource.TestCheckResourceAttr(testResourceName, "environment_variable.#", "1"),
				),
			},
			{
				Config: testAccLimelightEdgeFunctionBasicTemplate_update(fnName),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightEdgeFunctionExists(testResourceName, fnName),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", fnName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Test function updated"),
					resource.TestCheckResourceAttr(testResourceName, "function_archive", "testdata/edgefunc/py_function.zip"),
					resource.TestCheckResourceAttr(testResourceName, "handler", "hello_world.handler"),
					resource.TestCheckResourceAttr(testResourceName, "runtime", "python3"),
					resource.TestCheckResourceAttr(testResourceName, "memory", "1024"),
					resource.TestCheckResourceAttr(testResourceName, "timeout", "7000"),
					resource.TestCheckResourceAttr(testResourceName, "can_debug", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "function_sha256"),
					resource.TestCheckResourceAttr(testResourceName, "revision_id", "1"),
					resource.TestCheckResourceAttr(testResourceName, "reserved_concurrency", "5"),
					resource.TestCheckResourceAttr(testResourceName, "environment_variable.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceLimelightEdgeFunction_importBasic(t *testing.T) {
	testResourceName := "limelight_edgefunction.test_ef"
	fnName := "terraform_ef_test_import"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightEdgeFunctionCheckDestroy(state, fnName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightEdgeFunctionMinimalTemplate(fnName),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"function_archive"},
			},
		},
	})
}

func testAccLimelightEdgeFunctionCheckDestroy(state *terraform.State, fnName string) error {
	client := getEdgeFunctionsClient(testAccProvider.Meta().(map[string]interface{}))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "limelight_edgefunction" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		shortname, name, err := resourceLimelightEdgeFunctionSplitID(resourceID)

		if err != nil {
			return err
		}

		_, resp, err := client.GetEdgeFunction(name, shortname)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("error retrieving EdgeFunction with ID %s. Error: %v", resourceID, err)
		}

	}
	return fmt.Errorf("EdgeFunction %s still exists", fnName)
}

func testAccLimelightEdgeFunctionExists(testResourceName, fnName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		client := getEdgeFunctionsClient(testAccProvider.Meta().(map[string]interface{}))

		rs, ok := state.RootModule().Resources[testResourceName]
		if !ok {
			return fmt.Errorf("EdgeFunction %s not found in resources", fnName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("EdgeFunction ID not set in resource")
		}

		shortname, name, err := resourceLimelightEdgeFunctionSplitID(resourceID)
		if err != nil {
			return err
		}

		edgeFunction, resp, err := client.GetEdgeFunction(name, shortname)
		if err != nil {
			return fmt.Errorf("error retrieving EdgeFunction %s", resourceID)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error while checking if EdgeFunction %s exists. HTTP return code: %d", resourceID, resp.StatusCode)
		}

		if edgeFunction.Name == fnName {
			return nil
		}

		return fmt.Errorf("EdgeFunction %s wasn't found", fnName)
	}
}

func testAccLimelightEdgeFunctionMinimalTemplate(fnName string) string {
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
}`, getShortname(), fnName)
}

func testAccLimelightEdgeFunctionBasicTemplate(fnName string) string {
	return fmt.Sprintf(`
locals {
	fn_sha = filesha256("testdata/edgefunc/py_function.zip")
}

resource "limelight_edgefunction" "test_ef" {
	shortname            = "%s"
	description          = "Test function"
	name                 = "%s"
	function_archive     = "testdata/edgefunc/py_function.zip"
	handler              = "hello_world.handler"
	runtime              = "python3"
	function_sha256      = local.fn_sha
	memory               = 512
	reserved_concurrency = 3
	timeout              = 6000
	environment_variable {
		name  = "NAME"
		value = "World"
	}
}`, getShortname(), fnName)
}

func testAccLimelightEdgeFunctionBasicTemplate_update(fnName string) string {
	return fmt.Sprintf(`
locals {
	fn_sha = filesha256("testdata/edgefunc/py_function.zip")
}

resource "limelight_edgefunction" "test_ef" {
	shortname            = "%s"
	description          = "Test function updated"
	name                 = "%s"
	function_archive     = "testdata/edgefunc/py_function.zip"
	handler              = "hello_world.handler"
	runtime              = "python3"
	function_sha256      = local.fn_sha
	reserved_concurrency = 5
	memory               = 1024
	timeout              = 7000
	can_debug            = true
	environment_variable {
		name  = "NAME"
		value = "World"
	}
	environment_variable {
		name  = "MYKEY"
		value = "MyValue"
	}
}`, getShortname(), fnName)
}

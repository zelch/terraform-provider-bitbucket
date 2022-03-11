package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBitbucketDeploymentVariable_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	owner := os.Getenv("BITBUCKET_TEAM")
	resourceName := "bitbucket_deployment_variable.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDeploymentVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDeploymentVariableConfig(owner, rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeploymentVariableExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "deployment", "bitbucket_deployment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "key", "test"),
					resource.TestCheckResourceAttr(resourceName, "value", "test"),
					resource.TestCheckResourceAttr(resourceName, "secured", "false"),
				),
			},
		},
	})
}

func TestAccBitbucketDeploymentVariable_secure(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	owner := os.Getenv("BITBUCKET_TEAM")
	resourceName := "bitbucket_deployment_variable.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketDeploymentVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDeploymentVariableConfig(owner, rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketDeploymentVariableExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "deployment", "bitbucket_deployment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "key", "test"),
					resource.TestCheckResourceAttr(resourceName, "value", "test"),
					resource.TestCheckResourceAttr(resourceName, "secured", "true"),
				),
			},
		},
	})
}

func testAccCheckBitbucketDeploymentVariableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bitbucket_deployment_variable" {
			continue
		}

		repository, deployment := parseDeploymentId(rs.Primary.Attributes["deployment"])
		response, err := client.Get(fmt.Sprintf("2.0/repositories/%s/deployments_config/environments/%s/variables?pagelen=100", repository, deployment))

		if err == nil {
			return fmt.Errorf("The resource was found should have errored")
		}

		if response.StatusCode != 404 {
			return fmt.Errorf("Deployment Variable still exists")
		}
	}
	return nil
}

func testAccCheckBitbucketDeploymentVariableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found %s", n)
		}

		return nil
	}
}

func testAccBitbucketDeploymentVariableConfig(owner, rName string, secure bool) string {
	return fmt.Sprintf(`
resource "bitbucket_repository" "test" {
  owner = %[1]q
  name  = %[2]q
}

resource "bitbucket_deployment" "test" {
  name       = %[2]q
  stage      = "Test"
  repository = bitbucket_repository.test.id
}

resource "bitbucket_deployment_variable" "test" {
  key        = "test"
  value      = "test"
  deployment = bitbucket_deployment.test.id
  secured    = %[3]t
}
`, owner, rName, secure)
}

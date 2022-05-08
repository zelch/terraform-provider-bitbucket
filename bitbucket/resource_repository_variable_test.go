package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBitbucketRepositoryVariable_basic(t *testing.T) {

	owner := os.Getenv("BITBUCKET_TEAM")
	rName := acctest.RandomWithPrefix("tf-test")
	resourceName := "bitbucket_repository_variable.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketRepositoryVariableConfig(owner, rName, "test-val"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryVariableExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "key", "test"),
					resource.TestCheckResourceAttr(resourceName, "value", "test-val"),
					resource.TestCheckResourceAttr(resourceName, "secured", "false"),
				),
			},
			{
				Config: testAccBitbucketRepositoryVariableConfig(owner, rName, "test-val-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryVariableExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "key", "test"),
					resource.TestCheckResourceAttr(resourceName, "value", "test-val-2"),
					resource.TestCheckResourceAttr(resourceName, "secured", "false"),
				),
			},
		},
	})
}

func testAccCheckBitbucketRepositoryVariableDestroy(s *terraform.State) error {
	_, ok := s.RootModule().Resources["bitbucket_repository_variable.test"]
	if !ok {
		return fmt.Errorf("Not found %s", "bitbucket_repository_variable.test")
	}
	return nil
}

func testAccCheckBitbucketRepositoryVariableExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found %s", n)
		}

		return nil
	}
}

func testAccBitbucketRepositoryVariableConfig(team, rName, val string) string {
	return fmt.Sprintf(`
resource "bitbucket_repository" "test" {
  owner = %[1]q
  name  = %[2]q
}

resource "bitbucket_repository_variable" "test" {
  key        = "test"
  value      = %[3]q
  repository = bitbucket_repository.test.id
  secured = false
}
`, team, rName, val)
}

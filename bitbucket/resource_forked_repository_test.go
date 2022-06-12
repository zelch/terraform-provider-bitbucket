package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBitbucketForkedRepository_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	testUser := os.Getenv("BITBUCKET_TEAM")
	resourceName := "bitbucket_forked_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketForkedRepoConfig(testUser, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-fork", rName)),
					resource.TestCheckResourceAttr(resourceName, "owner", testUser),
					resource.TestCheckResourceAttr(resourceName, "scm", "git"),
					resource.TestCheckResourceAttr(resourceName, "has_wiki", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "fork_policy", "allow_forks"),
					resource.TestCheckResourceAttr(resourceName, "language", ""),
					resource.TestCheckResourceAttr(resourceName, "has_issues", "false"),
					resource.TestCheckResourceAttr(resourceName, "slug", fmt.Sprintf("%s-fork", rName)),
					resource.TestCheckResourceAttr(resourceName, "is_private", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "link.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "link.0.avatar.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "link.0.avatar.0.href"),
					resource.TestCheckResourceAttrSet(resourceName, "project_key"),
					resource.TestCheckResourceAttr(resourceName, "parent.%", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "parent.slug", "bitbucket_repository.test", "slug"),
					resource.TestCheckResourceAttrPair(resourceName, "parent.owner", "bitbucket_repository.test", "owner"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBitbucketForkedRepository_project(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-test")
	testUser := os.Getenv("BITBUCKET_TEAM")
	resourceName := "bitbucket_forked_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketForkedRepoProjectConfig(testUser, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-fork", rName)),
					resource.TestCheckResourceAttr(resourceName, "owner", testUser),
					resource.TestCheckResourceAttrPair(resourceName, "project_key", "bitbucket_project.test", "key"),
					resource.TestCheckResourceAttr(resourceName, "parent.%", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "parent.slug", "bitbucket_repository.test", "slug"),
					resource.TestCheckResourceAttrPair(resourceName, "parent.owner", "bitbucket_repository.test", "owner"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBitbucketForkedRepoConfig(testUser, rName string) string {
	return fmt.Sprintf(`
resource "bitbucket_repository" "test" {
  owner = %[1]q
  name  = %[2]q
}

resource "bitbucket_forked_repository" "test" {
  owner = bitbucket_repository.test.owner
  name  = "%[2]s-fork"

  parent = {
    slug  = bitbucket_repository.test.slug
	owner = bitbucket_repository.test.owner
  }
}
`, testUser, rName)
}

func testAccBitbucketForkedRepoProjectConfig(testUser, rName string) string {
	return fmt.Sprintf(`
resource "bitbucket_project" "test" {
  owner = %[1]q
  name  = %[2]q
  key   = "AAAAAAA"
}

resource "bitbucket_repository" "test" {
  owner = %[1]q
  name  = %[2]q
}
	
resource "bitbucket_forked_repository" "test" {
  owner       = bitbucket_repository.test.owner
  name        = "%[2]s-fork"
  project_key = bitbucket_project.test.key

  parent = {
    slug  = bitbucket_repository.test.slug
	owner = bitbucket_repository.test.owner
  }  
}
`, testUser, rName)
}

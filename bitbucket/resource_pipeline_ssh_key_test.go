package bitbucket

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBitbucketPipelineSshKey_basic(t *testing.T) {
	resourceName := "bitbucket_pipeline_ssh_key.test"

	rName := acctest.RandomWithPrefix("tf-test")
	owner := os.Getenv("BITBUCKET_TEAM")
	userEmail := os.Getenv("BITBUCKET_USERNAME")
	publicKey, privateKey, err := RandSSHKeyPairSize(2048, userEmail)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	publicKey2, privateKey2, err := RandSSHKeyPairSize(2048, userEmail)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketPipelineSshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketPipelineSshKeyConfig(owner, rName, publicKey, privateKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketPipelineSshKeyExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "workspace", "bitbucket_repository.test", "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttr(resourceName, "private_key", privateKey),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
			{
				Config: testAccBitbucketPipelineSshKeyConfig(owner, rName, publicKey2, privateKey2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketPipelineSshKeyExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "workspace", "bitbucket_repository.test", "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey2),
					resource.TestCheckResourceAttr(resourceName, "private_key", privateKey2),
				),
			},
		},
	})
}

func testAccCheckBitbucketPipelineSshKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(Clients).genClient
	pipeApi := client.ApiClient.PipelinesApi
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bitbucket_pipeline_ssh_key" {
			continue
		}

		workspace, repo, err := pipeSshKeyId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, res, err := pipeApi.GetRepositoryPipelineSshKeyPair(client.AuthContext, workspace, repo)

		if err == nil {
			return fmt.Errorf("The resource was found should have errored")
		}

		if res.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Pipeline Ssh Key Key still exists")
		}

	}
	return nil
}

func testAccCheckBitbucketPipelineSshKeyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pipeline Ssh Key Key ID is set")
		}
		return nil
	}
}

func testAccBitbucketPipelineSshKeyConfig(workspace, rName, pubKey, privKey string) string {
	return fmt.Sprintf(`
resource "bitbucket_repository" "test" {
  owner = %[1]q
  name  = %[2]q
}

resource "bitbucket_pipeline_ssh_key" "test" {
  workspace   = %[1]q
  repository  = bitbucket_repository.test.name
  public_key  = %[3]q
  private_key = %[4]q
}
`, workspace, rName, pubKey, privKey)
}

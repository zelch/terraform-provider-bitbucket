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

func TestAccBitbucketPipelineSshKnownHost_basic(t *testing.T) {
	resourceName := "bitbucket_pipeline_ssh_known_host.test"

	rName := acctest.RandomWithPrefix("tf-test")
	owner := os.Getenv("BITBUCKET_TEAM")
	publicKey, err := RandPlainSSHKeyPairSize(2048)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	publicKey2, err := RandPlainSSHKeyPairSize(2048)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketPipelineSshKnownHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketPipelineSshKnownHostConfig(owner, rName, publicKey, "[example.com]"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketPipelineSshKnownHostExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "workspace", "bitbucket_repository.test", "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "[example.com]"),
					resource.TestCheckResourceAttr(resourceName, "public_key.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_key.0.key_type", "ssh-rsa"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key.0.md5_fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key.0.sha256_fingerprint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBitbucketPipelineSshKnownHostConfig(owner, rName, publicKey2, "[example2.com]:22"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "workspace", "bitbucket_repository.test", "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "[example2.com]:22"),
					resource.TestCheckResourceAttr(resourceName, "public_key.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_key.0.key_type", "ssh-rsa"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key.0.md5_fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key.0.sha256_fingerprint"),
				),
			},
		},
	})
}

func testAccCheckBitbucketPipelineSshKnownHostDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(Clients).genClient
	pipeApi := client.ApiClient.PipelinesApi

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bitbucket_pipeline_ssh_known_host" {
			continue
		}

		workspace, repo, uuid, err := pipeSshKnownHostId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, res, err := pipeApi.GetRepositoryPipelineKnownHost(client.AuthContext, workspace, repo, uuid)

		if err == nil {
			return fmt.Errorf("The resource was found should have errored")
		}

		if res.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Pipeline Ssh Known Host still exists")
		}

	}
	return nil
}

func testAccCheckBitbucketPipelineSshKnownHostExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pipeline Ssh Key Known Host ID is set")
		}
		return nil
	}
}

func testAccBitbucketPipelineSshKnownHostConfig(workspace, rName, pubKey, host string) string {
	return fmt.Sprintf(`
resource "bitbucket_repository" "test" {
  owner = %[1]q
  name  = %[2]q
}

resource "bitbucket_pipeline_ssh_known_host" "test" {
  workspace  = %[1]q
  repository = bitbucket_repository.test.name
  hostname   = %[4]q

  public_key {
    key_type = "ssh-rsa" 
    key      = %[3]q
  }
}
`, workspace, rName, pubKey, host)
}

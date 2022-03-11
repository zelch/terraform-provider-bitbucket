package bitbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBitbucketPipelineSshKnownHost_basic(t *testing.T) {
	var pipelineSshKnownHost PiplineSshKnownHost
	resourceName := "bitbucket_pipeline_ssh_known_host.test"

	rName := acctest.RandomWithPrefix("tf-test")
	owner := os.Getenv("BITBUCKET_TEAM")
	userEmail := os.Getenv("BITBUCKET_USERNAME")
	publicKey, _, err := RandSSHKeyPairSize(2048, userEmail)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	publicKey2, _, err := RandSSHKeyPairSize(2048, userEmail)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketPipelineSshKnownHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketPipelineSshKnownHostConfig(owner, rName, publicKey, "example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketPipelineSshKnownHostExists(resourceName, &pipelineSshKnownHost),
					resource.TestCheckResourceAttrPair(resourceName, "workspace", "bitbucket_repository.test", "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "public_key.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_key.0.key_type", "RSA"),
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
				Config: testAccBitbucketPipelineSshKnownHostConfig(owner, rName, publicKey2, "example2.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "workspace", "bitbucket_repository.test", "owner"),
					resource.TestCheckResourceAttrPair(resourceName, "repository", "bitbucket_repository.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "example2.com"),
					resource.TestCheckResourceAttr(resourceName, "public_key.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_key.0.key_type", "RSA"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key.0.md5_fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key.0.sha256_fingerprint"),
				),
			},
		},
	})
}

func testAccCheckBitbucketPipelineSshKnownHostDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bitbucket_pipeline_ssh_known_host" {
			continue
		}

		workspace, repo, uuid, err := pipeSshKnownHostId(rs.Primary.ID)
		if err != nil {
			return err
		}

		response, err := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/known_hosts/%s",
			workspace, repo, uuid))

		if err == nil {
			return fmt.Errorf("The resource was found should have errored")
		}

		if response.StatusCode != 404 {
			return fmt.Errorf("Pipeline Ssh Known Host still exists")
		}

	}
	return nil
}

func testAccCheckBitbucketPipelineSshKnownHostExists(n string, pipelineSshKnownHost *PiplineSshKnownHost) resource.TestCheckFunc {
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
    key_type = "RSA" 
    key      = base64encode(%[3]q)
  }
}
`, workspace, rName, pubKey, host)
}

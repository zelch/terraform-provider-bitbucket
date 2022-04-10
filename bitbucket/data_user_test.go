package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUser_uuid(t *testing.T) {
	dataSourceName := "data.bitbucket_user.test"
	currUserDataSource := "data.bitbucket_current_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketUserUUIDConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", currUserDataSource, "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "nickname", currUserDataSource, "nickname"),
					resource.TestCheckResourceAttrPair(dataSourceName, "display_name", currUserDataSource, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "account_id", currUserDataSource, "account_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "account_status", currUserDataSource, "account_status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_staff", currUserDataSource, "is_staff"),
					// resource.TestCheckResourceAttrPair(dataSourceName, "username", currUserDataSource, "username"),
				),
			},
		},
	})
}

func TestAccUser_accountId(t *testing.T) {
	dataSourceName := "data.bitbucket_user.test"
	currUserDataSource := "data.bitbucket_current_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "uuid", currUserDataSource, "uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "nickname", currUserDataSource, "nickname"),
					resource.TestCheckResourceAttrPair(dataSourceName, "display_name", currUserDataSource, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "account_id", currUserDataSource, "account_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "account_status", currUserDataSource, "account_status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_staff", currUserDataSource, "is_staff"),
					// resource.TestCheckResourceAttrPair(dataSourceName, "username", currUserDataSource, "username"),
				),
			},
		},
	})
}

func testAccBitbucketUserConfig() string {
	return `
data "bitbucket_current_user" "test" {}

data "bitbucket_user" "test" {
  account_id = data.bitbucket_current_user.test.account_id
}
`
}

func testAccBitbucketUserUUIDConfig() string {
	return `
data "bitbucket_current_user" "test" {}

data "bitbucket_user" "test" {
  uuid = data.bitbucket_current_user.test.uuid
}
`
}

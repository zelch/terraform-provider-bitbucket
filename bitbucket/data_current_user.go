package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCurrentUser() *schema.Resource {
	return &schema.Resource{
		Read: dataReadCurrentUser,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_staff": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataReadCurrentUser(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	curUser, err := c.Get("2.0/user")
	if err != nil {
		return err
	}

	if curUser.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if curUser.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching user")
	}

	body, readerr := ioutil.ReadAll(curUser.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Current User Response JSON: %v", string(body))

	var u apiUser

	decodeerr := json.Unmarshal(body, &u)
	if decodeerr != nil {
		return decodeerr
	}

	log.Printf("[DEBUG] Current User Response Decoded: %#v", u)

	d.SetId(u.UUID)
	d.Set("uuid", u.UUID)
	d.Set("nickname", u.Nickname)
	d.Set("display_name", u.DisplayName)
	d.Set("account_id", u.AccountId)
	d.Set("account_status", u.AccountStatus)
	d.Set("is_staff", u.IsStaff)

	return nil
}

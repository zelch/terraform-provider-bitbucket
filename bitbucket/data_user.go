package bitbucket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataReadUser,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataReadUser(d *schema.ResourceData, m interface{}) error {
	c := m.(Clients).genClient
	usersApi := c.ApiClient.UsersApi
	var selectedUser string

	if v, ok := d.GetOk("uuid"); ok && v.(string) != "" {
		selectedUser = v.(string)
	}

	user, userRes, err := usersApi.UsersSelectedUserGet(c.AuthContext, selectedUser)
	if err != nil {
		return fmt.Errorf("error reading User (%s): %w", selectedUser, err)
	}

	if userRes.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if userRes.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching user")
	}

	log.Printf("[DEBUG] User: %#v", user)

	d.SetId(user.Uuid)
	d.Set("uuid", user.Uuid)
	d.Set("username", user.Username)
	d.Set("display_name", user.DisplayName)

	return nil
}

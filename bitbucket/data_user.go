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
			"nickname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

func dataReadUser(d *schema.ResourceData, m interface{}) error {
	c := m.(Clients).genClient
	usersApi := c.ApiClient.UsersApi
	var selectedUser string

	if v, ok := d.GetOk("account_id"); ok && v.(string) != "" {
		selectedUser = v.(string)
	}

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
	d.Set("nickname", user.Nickname)
	d.Set("username", user.Username)
	d.Set("display_name", user.DisplayName)
	d.Set("account_id", user.AccountId)
	d.Set("account_status", user.AccountStatus)
	d.Set("is_staff", user.IsStaff)

	return nil
}

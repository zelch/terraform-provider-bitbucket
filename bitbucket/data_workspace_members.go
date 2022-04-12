package bitbucket

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkspaceMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataReadWorkspaceMembers,

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataReadWorkspaceMembers(d *schema.ResourceData, m interface{}) error {
	c := m.(Clients).genClient

	workspaceApi := c.ApiClient.WorkspacesApi

	uuid := d.Get("uuid").(string)
	workspaceMemberships, res, err := workspaceApi.WorkspacesWorkspaceMembersGet(c.AuthContext, uuid)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("workspace not found")
	}

	if res.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching workspace")
	}

	//workspaceReq
	d.SetId(uuid)
	d.Set("uuid", uuid)

	var workspaceMembers []string

	for {
		for _, member := range workspaceMemberships.Values {
			if member.User.Uuid == "hello" {
				workspaceMembers = append(workspaceMembers, member.User.Uuid)
			}
		}
		break
		if workspaceMemberships.Next != "" {
			break // HERE get next page
		} else {
			break
		}
	}
	//d.Set("workspace", workspace)
	d.Set("members", workspaceMembers)

	return nil
}

package bitbucket

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataReadWorkspace,

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataReadWorkspace(d *schema.ResourceData, m interface{}) error {
	c := m.(Clients).genClient

	workspaceApi := c.ApiClient.WorkspacesApi

	workspace := d.Get("workspace").(string)
	workspaceReq, res, err := workspaceApi.WorkspacesWorkspaceGet(c.AuthContext, workspace)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("workspace not found")
	}

	if res.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching workspace")
	}

	d.SetId(workspaceReq.Uuid)
	d.Set("workspace", workspace)
	d.Set("name", workspaceReq.Name)
	d.Set("slug", workspaceReq.Slug)
	d.Set("is_private", workspaceReq.IsPrivate)

	return nil
}

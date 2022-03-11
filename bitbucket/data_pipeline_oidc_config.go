package bitbucket

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPipelineOidcConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataReadPipelineOidcConfig,

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataReadPipelineOidcConfig(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	workspace := d.Get("workspace").(string)
	req, err := c.Get(fmt.Sprintf("2.0/workspaces/%s/pipelines-config/identity/oidc/.well-known/openid-configuration", workspace))
	if err != nil {
		return err
	}

	if req.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if req.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching user")
	}

	body, readerr := ioutil.ReadAll(req.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Pipeline Oidc Config Response JSON: %v", string(body))

	d.SetId(workspace)
	d.Set("workspace", workspace)
	d.Set("oidc_config", string(body))

	return nil
}

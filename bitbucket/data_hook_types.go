package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PaginatedHookTypes struct {
	Values []HookType `json:"values,omitempty"`
	Page   int        `json:"page,omitempty"`
	Size   int        `json:"size,omitempty"`
	Next   string     `json:"next,omitempty"`
}

type HookType struct {
	Event       string `json:"event"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

func dataHookTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataReadHookTypes,

		Schema: map[string]*schema.Schema{
			"subject_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hook_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataReadHookTypes(d *schema.ResourceData, m interface{}) error {
	c := m.(Clients).httpClient

	subjectType := d.Get("subject_type").(string)
	hookTypes, err := c.Get(fmt.Sprintf("2.0/hook_events/%s", subjectType))
	if err != nil {
		return err
	}

	if hookTypes.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if hookTypes.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("internal server error fetching hook types")
	}

	body, readerr := ioutil.ReadAll(hookTypes.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Hook Types Response JSON: %v", string(body))

	var hookTypePages PaginatedHookTypes

	decodeerr := json.Unmarshal(body, &hookTypePages)
	if decodeerr != nil {
		return decodeerr
	}

	log.Printf("[DEBUG] Hook Type Pages Response Decoded: %#v", hookTypePages)

	d.SetId(subjectType)
	d.Set("hook_types", flattenHookTypes(hookTypePages.Values))

	return nil
}

func flattenHookTypes(HookTypes []HookType) []interface{} {
	if len(HookTypes) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, btRaw := range HookTypes {
		log.Printf("[DEBUG] HookType Response Decoded: %#v", btRaw)

		hookType := map[string]interface{}{
			"event":       btRaw.Event,
			"category":    btRaw.Category,
			"label":       btRaw.Label,
			"description": btRaw.Description,
		}

		tfList = append(tfList, hookType)
	}

	return tfList
}

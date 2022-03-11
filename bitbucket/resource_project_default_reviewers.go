package bitbucket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectDefaultReviewers() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectDefaultReviewersCreate,
		Read:   resourceProjectDefaultReviewersRead,
		Update: resourceProjectDefaultReviewersUpdate,
		Delete: resourceProjectDefaultReviewersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reviewers": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceProjectDefaultReviewersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	for _, user := range d.Get("reviewers").(*schema.Set).List() {
		reviewerResp, err := client.Post(fmt.Sprintf("internal/workspaces/%s/projects/%s/default-reviewers/",
			d.Get("workspace").(string),
			d.Get("project").(string),
		), nil)

		if err != nil {
			return err
		}

		if reviewerResp.StatusCode != 200 {
			return fmt.Errorf("failed to project create reviewer %s got code %d", user.(string), reviewerResp.StatusCode)
		}

		defer reviewerResp.Body.Close()
	}

	d.SetId(fmt.Sprintf("%s/%s/reviewers", d.Get("workspace").(string), d.Get("project").(string)))
	return resourceProjectDefaultReviewersRead(d, m)
}

func resourceProjectDefaultReviewersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	workspace, repo, err := defaultReviewersId(d.Id())
	if err != nil {
		return err
	}
	resourceURL := fmt.Sprintf("2.0/workspaces/%s/projects/%s/default-reviewers", workspace, repo)

	res, err := client.Get(resourceURL)
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		log.Printf("[WARN] Default Reviewers (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var reviewers PaginatedReviewers
	var terraformReviewers []string

	for {
		reviewersResponse, err := client.Get(resourceURL)
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(reviewersResponse.Body)
		err = decoder.Decode(&reviewers)
		if err != nil {
			return err
		}

		for _, reviewer := range reviewers.Values {
			terraformReviewers = append(terraformReviewers, reviewer.UUID)
		}

		if reviewers.Next != "" {
			nextPage := reviewers.Page + 1
			resourceURL = fmt.Sprintf("2.0/workspaces/%s/projects/%s/default-reviewers?page=%d", workspace, repo, nextPage)
			reviewers = PaginatedReviewers{}
		} else {
			break
		}
	}

	d.Set("workspace", workspace)
	d.Set("project", repo)
	d.Set("reviewers", terraformReviewers)

	return nil
}

func resourceProjectDefaultReviewersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	oraw, nraw := d.GetChange("reviewers")
	o := oraw.(*schema.Set)
	n := nraw.(*schema.Set)

	add := n.Difference(o)
	remove := o.Difference(n)

	for _, user := range add.List() {
		reviewerResp, err := client.PutOnly(fmt.Sprintf("2.0/workspaces/%s/projects/%s/default-reviewers/%s",
			d.Get("workspace").(string),
			d.Get("project").(string),
			user,
		))

		if err != nil {
			return err
		}

		if reviewerResp.StatusCode != 200 {
			return fmt.Errorf("failed to create reviewer %s got code %d", user.(string), reviewerResp.StatusCode)
		}

		defer reviewerResp.Body.Close()
	}

	for _, user := range remove.List() {
		resp, err := client.Delete(fmt.Sprintf("2.0/workspaces/%s/projects/%s/default-reviewers/%s",
			d.Get("workspace").(string),
			d.Get("project").(string),
			user.(string),
		))

		if err != nil {
			return err
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("[%d] Could not delete %s from default reviewers",
				resp.StatusCode,
				user.(string),
			)
		}
		defer resp.Body.Close()
	}

	return resourceProjectDefaultReviewersRead(d, m)
}

func resourceProjectDefaultReviewersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	for _, user := range d.Get("reviewers").(*schema.Set).List() {
		resp, err := client.Delete(fmt.Sprintf("2.0/workspaces/%s/projects/%s/default-reviewers/%s",
			d.Get("workspace").(string),
			d.Get("project").(string),
			user.(string),
		))

		if err != nil {
			return err
		}

		if resp.StatusCode != 204 {
			return fmt.Errorf("[%d] Could not delete %s from default reviewer",
				resp.StatusCode,
				user.(string),
			)
		}
		defer resp.Body.Close()
	}
	return nil
}

// func defaultReviewersId(id string) (string, string, error) {
// 	parts := strings.Split(id, "/")

// 	if len(parts) != 3 {
// 		return "", "", fmt.Errorf("unexpected format of ID (%q), expected OWNER/REPOSITORY/reviewers", id)
// 	}

// 	return parts[0], parts[1], nil
// }

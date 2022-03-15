package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PiplineSshKey struct {
	PrivateKey string `json:"private_key,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`
}

func resourcePipelineSshKey() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineSshKeysPut,
		Read:   resourcePipelineSshKeysRead,
		Update: resourcePipelineSshKeysPut,
		Delete: resourcePipelineSshKeysDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePipelineSshKeysPut(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	pipeSshKey := expandPipelineSshKey(d)
	log.Printf("[DEBUG] Pipeline Ssh Key Request: %#v", pipeSshKey)
	bytedata, err := json.Marshal(pipeSshKey)

	if err != nil {
		return err
	}

	repo := d.Get("repository").(string)
	workspace := d.Get("workspace").(string)
	_, err = client.Put(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/key_pair",
		workspace, repo),
		bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	d.SetId(string(fmt.Sprintf("%s/%s", workspace, repo)))

	return resourcePipelineSshKeysRead(d, m)
}

func resourcePipelineSshKeysRead(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	workspace, repo, err := pipeSshKeyId(d.Id())
	if err != nil {
		return err
	}
	pipeSshKeysReq, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/key_pair", workspace, repo))

	if pipeSshKeysReq.StatusCode == 404 {
		log.Printf("[WARN] Pipeline Ssh Key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if pipeSshKeysReq.Body == nil {
		return fmt.Errorf("error getting Pipeline Ssh Key (%s): empty response", d.Id())
	}

	var pipeSshKey *PiplineSshKey
	body, readerr := ioutil.ReadAll(pipeSshKeysReq.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Pipeline Ssh Key Response JSON: %v", string(body))

	decodeerr := json.Unmarshal(body, &pipeSshKey)
	if decodeerr != nil {
		return decodeerr
	}

	log.Printf("[DEBUG] Pipeline Ssh Key Response Decoded: %#v", pipeSshKey)

	d.Set("repository", repo)
	d.Set("workspace", workspace)
	d.Set("public_key", pipeSshKey.PublicKey)
	d.Set("private_key", d.Get("private_key").(string))

	return nil
}

func resourcePipelineSshKeysDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	workspace, repo, err := pipeSshKeyId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/key_pair", workspace, repo))

	if err != nil {
		return fmt.Errorf("error deleting Pipeline Ssh Key (%s): %w", d.Id(), err)
	}

	return err
}

func expandPipelineSshKey(d *schema.ResourceData) *PiplineSshKey {
	key := &PiplineSshKey{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
	}

	return key
}

func pipeSshKeyId(id string) (string, string, error) {
	parts := strings.Split(id, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected format of ID (%q), expected WORKSPACE-ID/REPO-ID", id)
	}

	return parts[0], parts[1], nil
}

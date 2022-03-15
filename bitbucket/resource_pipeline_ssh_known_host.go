package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type PiplineSshKnownHost struct {
	UUID      string                        `json:"uuid,omitempty"`
	Hostname  string                        `json:"hostname,omitempty"`
	PublicKey *PiplineSshKnownHostPublicKey `json:"public_key,omitempty"`
}

type PiplineSshKnownHostPublicKey struct {
	KeyType           string `json:"key_type,omitempty"`
	Key               string `json:"key,omitempty"`
	MD5Fingerprint    string `json:"md5_fingerprint,omitempty"`
	SHA256Fingerprint string `json:"sha256_fingerprint,omitempty"`
}

func resourcePipelineSshKnownHost() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineSshKnownHostsCreate,
		Read:   resourcePipelineSshKnownHostsRead,
		Update: resourcePipelineSshKnownHostsUpdate,
		Delete: resourcePipelineSshKnownHostsDelete,
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
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_key": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Ed25519", "ECDSA", "RSA", "DSA"}, false),
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"md5_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha256_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePipelineSshKnownHostsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	pipeSshKnownHost := expandPipelineSshKnownHost(d)
	log.Printf("[DEBUG] Pipeline Ssh Key Request: %#v", pipeSshKnownHost)
	bytedata, err := json.Marshal(pipeSshKnownHost)

	if err != nil {
		return err
	}

	repo := d.Get("repository").(string)
	workspace := d.Get("workspace").(string)
	resp, err := client.Post(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/known_hosts/",
		workspace, repo),
		bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Pipeline Ssh Known Hosts Response JSON: %v", string(body))

	var host PiplineSshKnownHost

	decodeerr := json.Unmarshal(body, &host)
	if decodeerr != nil {
		return decodeerr
	}

	log.Printf("[DEBUG] Pipeline Ssh Known Host Pages Response Decoded: %#v", host)

	d.SetId(string(fmt.Sprintf("%s/%s/%s", workspace, repo, host.UUID)))

	return resourcePipelineSshKnownHostsRead(d, m)
}

func resourcePipelineSshKnownHostsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	workspace, repo, uuid, err := pipeSshKnownHostId(d.Id())
	if err != nil {
		return err
	}

	pipeSshKnownHost := expandPipelineSshKnownHost(d)
	log.Printf("[DEBUG] Pipeline Ssh Key Request: %#v", pipeSshKnownHost)
	bytedata, err := json.Marshal(pipeSshKnownHost)

	if err != nil {
		return err
	}

	_, err = client.Put(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/known_hosts/%s",
		workspace, repo, uuid),
		bytes.NewBuffer(bytedata))

	if err != nil {
		return err
	}

	return resourcePipelineSshKnownHostsRead(d, m)
}

func resourcePipelineSshKnownHostsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	workspace, repo, uuid, err := pipeSshKnownHostId(d.Id())
	if err != nil {
		return err
	}
	pipeSshKnownHostsReq, _ := client.Get(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/known_hosts/%s",
		workspace, repo, uuid))

	if pipeSshKnownHostsReq.StatusCode == 404 {
		log.Printf("[WARN] Pipeline Ssh Key (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if pipeSshKnownHostsReq.Body == nil {
		return fmt.Errorf("error getting Pipeline Ssh Key (%s): empty response", d.Id())
	}

	var pipeSshKnownHost *PiplineSshKnownHost
	body, readerr := ioutil.ReadAll(pipeSshKnownHostsReq.Body)
	if readerr != nil {
		return readerr
	}

	log.Printf("[DEBUG] Pipeline Ssh Key Response JSON: %v", string(body))

	decodeerr := json.Unmarshal(body, &pipeSshKnownHost)
	if decodeerr != nil {
		return decodeerr
	}

	log.Printf("[DEBUG] Pipeline Ssh Key Response Decoded: %#v", pipeSshKnownHost)

	d.Set("repository", repo)
	d.Set("workspace", workspace)
	d.Set("hostname", pipeSshKnownHost.Hostname)
	d.Set("uuid", pipeSshKnownHost.UUID)
	d.Set("public_key", flattenPipelineSshKnownHost(pipeSshKnownHost.PublicKey))

	return nil
}

func resourcePipelineSshKnownHostsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Clients).httpClient

	workspace, repo, uuid, err := pipeSshKnownHostId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(fmt.Sprintf("2.0/repositories/%s/%s/pipelines_config/ssh/known_hosts/%s",
		workspace, repo, uuid))

	if err != nil {
		return fmt.Errorf("error deleting Pipeline Ssh Key (%s): %w", d.Id(), err)
	}

	return err
}

func expandPipelineSshKnownHost(d *schema.ResourceData) *PiplineSshKnownHost {
	key := &PiplineSshKnownHost{
		Hostname:  d.Get("hostname").(string),
		PublicKey: expandPipelineSshKnownHostKey(d.Get("public_key").([]interface{})),
	}

	return key
}

func expandPipelineSshKnownHostKey(pubKey []interface{}) *PiplineSshKnownHostPublicKey {
	tfMap, _ := pubKey[0].(map[string]interface{})

	key := &PiplineSshKnownHostPublicKey{
		KeyType: tfMap["key_type"].(string),
		Key:     tfMap["key"].(string),
	}

	return key
}

func flattenPipelineSshKnownHost(rp *PiplineSshKnownHostPublicKey) []interface{} {
	if rp == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"key_type":           rp.KeyType,
		"key":                rp.Key,
		"md5_fingerprint":    rp.MD5Fingerprint,
		"sha256_fingerprint": rp.SHA256Fingerprint,
	}

	return []interface{}{m}
}

func pipeSshKnownHostId(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("unexpected format of ID (%q), expected WORKSPACE-ID/REPO-ID/UUID", id)
	}

	return parts[0], parts[1], parts[2], nil
}

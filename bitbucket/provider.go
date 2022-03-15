package bitbucket

import (
	"context"
	"net/http"

	"github.com/DrFaust92/bitbucket-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderConfig struct {
	ApiClient   *bitbucket.APIClient
	AuthContext context.Context
}

type Clients struct {
	genClient  ProviderConfig
	httpClient Client
}

// Provider will create the necessary terraform provider to talk to the Bitbucket APIs you should
// specify a USERNAME and PASSWORD or a OAUTH Token
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
				ConflictsWith: []string{"oauth_token"},
				RequiredWith:  []string{"password"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
				ConflictsWith: []string{"oauth_token"},
				RequiredWith:  []string{"username"},
			},
			"oauth_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("BITBUCKET_OAUTH_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_hook":                    resourceHook(),
			"bitbucket_group":                   resourceGroup(),
			"bitbucket_group_membership":        resourceGroupMembership(),
			"bitbucket_default_reviewers":       resourceDefaultReviewers(),
			"bitbucket_repository":              resourceRepository(),
			"bitbucket_repository_variable":     resourceRepositoryVariable(),
			"bitbucket_project":                 resourceProject(),
			"bitbucket_deploy_key":              resourceDeployKey(),
			"bitbucket_pipeline_ssh_key":        resourcePipelineSshKey(),
			"bitbucket_pipeline_ssh_known_host": resourcePipelineSshKnownHost(),
			"bitbucket_ssh_key":                 resourceSshKey(),
			"bitbucket_branch_restriction":      resourceBranchRestriction(),
			"bitbucket_branching_model":         resourceBranchingModel(),
			"bitbucket_deployment":              resourceDeployment(),
			"bitbucket_deployment_variable":     resourceDeploymentVariable(),
			"bitbucket_workspace_hook":          resourceWorkspaceHook(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_ip_ranges":                 dataIPRanges(),
			"bitbucket_pipeline_oidc_config":      dataPipelineOidcConfig(),
			"bitbucket_pipeline_oidc_config_keys": dataPipelineOidcConfigKeys(),
			"bitbucket_hook_types":                dataHookTypes(),
			"bitbucket_user":                      dataUser(),
			"bitbucket_current_user":              dataCurrentUser(),
			"bitbucket_workspace":                 dataWorkspace(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	authCtx := context.Background()

	client := &Client{
		HTTPClient: &http.Client{},
	}

	if v, ok := d.GetOk("username"); ok && v.(string) != "" {
		user := v.(string)
		client.Username = &user
	}

	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		pass := v.(string)
		client.Password = &pass
	}

	if v, ok := d.GetOk("oauth_token"); ok && v.(string) != "" {
		token := v.(string)
		client.OAuthToken = &token
	}

	conf := bitbucket.NewConfiguration()
	apiClient := ProviderConfig{
		ApiClient:   bitbucket.NewAPIClient(conf),
		AuthContext: authCtx,
	}

	clients := Clients{
		genClient:  apiClient,
		httpClient: *client,
	}

	return clients, nil
}

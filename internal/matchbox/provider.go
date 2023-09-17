package matchbox

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a Provider for Matchbox.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ca": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"matchbox_profile": resourceProfile(),
			"matchbox_group":   resourceGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	ca := d.Get("ca").(string)
	clientCert := d.Get("client_cert").(string)
	clientKey := d.Get("client_key").(string)
	endpoint := d.Get("endpoint").(string)

	config := &Config{
		Endpoint:   endpoint,
		ClientCert: []byte(clientCert),
		ClientKey:  []byte(clientKey),
		CA:         []byte(ca),
	}

	client, err := NewMatchboxClient(config)
	if err != nil {
		return client, fmt.Errorf("failed to create Matchbox client or connect to %s: %v", endpoint, err)
	}
	return client, err
}

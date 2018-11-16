package matchbox

import (
	"errors"
	"fmt"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns the provider schema to Terraform.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// Provider configuration
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"endpoints"},
			},
			"endpoints": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"endpoint"},
			},
			"client_cert": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ca": &schema.Schema{
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

	endpoints := []string{}
	if values, hasMultipleEndpoints := d.GetOk("endpoints"); hasMultipleEndpoints {
		for _, value := range values.([]interface{}) {
			endpoints = append(endpoints, value.(string))
		}
	} else if value, hasEndpoint := d.GetOk("endpoint"); hasEndpoint {
		endpoints = append(endpoints, value.(string))
	} else {
		return nil, errors.New("Either endpoints or endpoint has to be set")
	}

	clients := []*matchbox.Client{}
	for _, endpoint := range endpoints {
		config := &Config{
			Endpoint:   endpoint,
			ClientCert: []byte(clientCert),
			ClientKey:  []byte(clientKey),
			CA:         []byte(ca),
		}

		client, err := NewMatchboxClient(config)
		if err != nil {
			return client, fmt.Errorf("failed to create Matchbox client or connect to %s: %v", endpoints, err)
		}
		clients = append(clients, client)
	}

	return clients, nil
}

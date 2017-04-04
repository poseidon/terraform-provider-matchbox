package matchbox

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Delete: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

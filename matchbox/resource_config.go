package matchbox

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceConfigCreate,
		Read:   resourceConfigRead,
		Delete: resourceConfigDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceConfigRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

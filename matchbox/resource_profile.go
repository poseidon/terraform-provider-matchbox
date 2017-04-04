package matchbox

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileCreate,
		Read:   resourceProfileRead,
		Delete: resourceProfileDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceProfileRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

package matchbox

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileCreate,
		Read:   resourceProfileRead,
		Update: resourceProfileUpdate,
		Delete: resourceProfileDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"kernel": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"initrd": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"args": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
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

func resourceProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

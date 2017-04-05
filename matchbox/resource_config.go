package matchbox

import (
	"context"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceConfigCreate,
		Read:   resourceConfigRead,
		Update: resourceConfigUpdate,
		Delete: resourceConfigDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)

	_, err := client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
		Name:   name,
		Config: []byte(d.Get("contents").(string)),
	})

	d.SetId(name)
	return err
}

func resourceConfigRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

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
		Delete: resourceConfigDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contents": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	if err != nil {
		return err
	}

	d.SetId(name)
	return err
}

func resourceConfigRead(d *schema.ResourceData, meta interface{}) error {
	// TODO: Read API is not yet implemented. Must delete and re-create each time.
	d.SetId("")
	return nil
}

func resourceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO: Delete API is not yet implemented
	return nil
}

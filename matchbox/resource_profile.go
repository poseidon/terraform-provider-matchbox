package matchbox

import (
	"context"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
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
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)

	var initrds []string
	for _, initrd := range d.Get("initrd").([]interface{}) {
		initrds = append(initrds, initrd.(string))
	}

	var args []string
	for _, arg := range d.Get("args").([]interface{}) {
		args = append(args, arg.(string))
	}

	profile := &storagepb.Profile{
		Id:         name,
		IgnitionId: d.Get("config").(string),
		Boot: &storagepb.NetBoot{
			Kernel: d.Get("kernel").(string),
			Initrd: initrds,
			Args:   args,
		},
	}

	_, err := client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
		Profile: profile,
	})

	d.SetId(name)
	return err
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

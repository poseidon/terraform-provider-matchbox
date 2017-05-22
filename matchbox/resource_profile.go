package matchbox

import (
	"context"
	"errors"
	"fmt"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
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
			"raw_ignition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"container_linux_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kernel": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"initrd": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"args": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// resourceProfileCreate creates a Profile and its associated configs. Partial
// creates do not modify state and can be retried safely.
func resourceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	if err := validateResourceProfile(d); err != nil {
		return err
	}

	// Profile
	name := d.Get("name").(string)
	clcName, clc := containerLinuxConfig(d)

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
		IgnitionId: clcName,
		Boot: &storagepb.NetBoot{
			Kernel: d.Get("kernel").(string),
			Initrd: initrds,
			Args:   args,
		},
	}
	_, err := client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
		Profile: profile,
	})
	if err != nil {
		return err
	}

	// Container Linux Config
	_, err = client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
		Name:   clcName,
		Config: []byte(clc),
	})
	if err != nil {
		return err
	}

	d.SetId(profile.GetId())
	return err
}

func validateResourceProfile(d *schema.ResourceData) error {
	_, hasRAW := d.GetOk("raw_ignition")
	_, hasCLC := d.GetOk("container_linux_config")
	if hasCLC && hasRAW {
		return errors.New("container_linux_config and raw_ignition are mutually exclusive")
	}

	if !hasCLC && !hasRAW {
		return errors.New("container_linux_config or raw_ignition are required")
	}
	return nil
}

func resourceProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)
	_, err := client.Profiles.ProfileGet(ctx, &serverpb.ProfileGetRequest{
		Id: name,
	})
	if err != nil {
		// resource doesn't exist or is corrupted
		d.SetId("")
		return nil
	}

	_, err = client.Ignition.IgnitionGet(ctx, &serverpb.IgnitionGetRequest{
		Name: containerLinuxConfigName(d),
	})
	if err != nil {
		// resource doesn't exist or is corrupted
		d.SetId("")
		return nil
	}

	return nil
}

// resourceProfileDelete deletes a Profile and its associated configs. Partial
// deletes leave state unchanged and can be retried (deleting resources which
// no longer exist is a no-op).
func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	// Profile
	name := d.Get("name").(string)
	_, err := client.Profiles.ProfileDelete(ctx, &serverpb.ProfileDeleteRequest{
		Id: name,
	})
	if err != nil {
		return err
	}

	// Container Linux Config
	_, err = client.Ignition.IgnitionDelete(ctx, &serverpb.IgnitionDeleteRequest{
		Name: containerLinuxConfigName(d),
	})
	if err != nil {
		return err
	}

	// resource can be destroyed in state
	d.SetId("")
	return nil
}

func containerLinuxConfigName(d *schema.ResourceData) string {
	filename, _ := containerLinuxConfig(d)
	return filename
}

func containerLinuxConfig(d *schema.ResourceData) (filename, config string) {
	name := d.Get("name").(string)

	if content, ok := d.GetOk("container_linux_config"); ok {
		return fmt.Sprintf("%s.yaml.tmpl", name), content.(string)
	}

	if content, ok := d.GetOk("raw_ignition"); ok {
		return fmt.Sprintf("%s.ign", name), content.(string)
	}

	return
}

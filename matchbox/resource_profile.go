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
			// recommended
			"container_linux_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"raw_ignition": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"generic_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// resourceProfileCreate creates a Profile and its associated configs. Partial
// creates do not modify state and can be retried safely.
func resourceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	clients := meta.([]*matchbox.Client)
	ctx := context.TODO()

	if err := validateResourceProfile(d); err != nil {
		return err
	}

	// Profile
	name := d.Get("name").(string)
	// NetBoot
	var initrds []string
	for _, initrd := range d.Get("initrd").([]interface{}) {
		initrds = append(initrds, initrd.(string))
	}
	var args []string
	for _, arg := range d.Get("args").([]interface{}) {
		args = append(args, arg.(string))
	}
	// Container Linux config / Ignition config
	clcName, _ := containerLinuxConfig(d)
	// Generic (experimental) config
	genericName, _ := genericConfig(d)

	profile := &storagepb.Profile{
		Id: name,
		Boot: &storagepb.NetBoot{
			Kernel: d.Get("kernel").(string),
			Initrd: initrds,
			Args:   args,
		},
		IgnitionId: clcName,
		GenericId:  genericName,
	}

	for _, client := range clients {
		// Profile
		_, err := client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
			Profile: profile,
		})
		if err != nil {
			return err
		}

		// Container Linux Config
		if name, content := containerLinuxConfig(d); content != "" {
			_, err = client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
				Name:   name,
				Config: []byte(content),
			})
			if err != nil {
				return err
			}
		}

		// Generic Config
		if name, content := genericConfig(d); content != "" {
			_, err = client.Generic.GenericPut(ctx, &serverpb.GenericPutRequest{
				Name:   name,
				Config: []byte(content),
			})
			if err != nil {
				return err
			}
		}
	}

	d.SetId(profile.GetId())
	return nil
}

func validateResourceProfile(d *schema.ResourceData) error {
	_, hasRAW := d.GetOk("raw_ignition")
	_, hasCLC := d.GetOk("container_linux_config")
	if hasCLC && hasRAW {
		return errors.New("container_linux_config and raw_ignition are mutually exclusive")
	}
	return nil
}

func resourceProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.([]*matchbox.Client)[0]
	ctx := context.TODO()

	// Profile
	name := d.Get("name").(string)
	_, err := client.Profiles.ProfileGet(ctx, &serverpb.ProfileGetRequest{
		Id: name,
	})
	if err != nil {
		// resource doesn't exist or is corrupted and needs creating
		d.SetId("")
		return nil
	}

	// Container Linux Config
	if name, content := containerLinuxConfig(d); content != "" {
		_, err = client.Ignition.IgnitionGet(ctx, &serverpb.IgnitionGetRequest{
			Name: name,
		})
		if err != nil {
			// resource doesn't exist or is corrupted and needs creating
			d.SetId("")
			return nil
		}
	}

	// Generic Config
	if name, content := genericConfig(d); content != "" {
		_, err = client.Generic.GenericGet(ctx, &serverpb.GenericGetRequest{
			Name: name,
		})
		if err != nil {
			// resource doesn't exist or is corrupted and needs creating
			d.SetId("")
			return nil
		}
	}

	return nil
}

// resourceProfileDelete deletes a Profile and its associated configs. Partial
// deletes leave state unchanged and can be retried (deleting resources which
// no longer exist is a no-op).
func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	clients := meta.([]*matchbox.Client)
	ctx := context.TODO()

	for _, client := range clients {
		// Profile
		name := d.Get("name").(string)
		_, err := client.Profiles.ProfileDelete(ctx, &serverpb.ProfileDeleteRequest{
			Id: name,
		})
		if err != nil {
			return err
		}

		// Container Linux Config
		if name, content := containerLinuxConfig(d); content != "" {
			_, err = client.Ignition.IgnitionDelete(ctx, &serverpb.IgnitionDeleteRequest{
				Name: name,
			})
			if err != nil {
				return err
			}
		}

		// Generic Config
		if name, content := genericConfig(d); content != "" {
			_, err = client.Generic.GenericDelete(ctx, &serverpb.GenericDeleteRequest{
				Name: name,
			})
			if err != nil {
				return err
			}
		}
	}

	// resource can be destroyed in state
	d.SetId("")
	return nil
}

func containerLinuxConfig(d *schema.ResourceData) (filename, config string) {
	// use profile name to generate Container Linux and Ignition filenames
	name := d.Get("name").(string)

	if content, ok := d.GetOk("container_linux_config"); ok {
		return fmt.Sprintf("%s.yaml.tmpl", name), content.(string)
	}

	if content, ok := d.GetOk("raw_ignition"); ok {
		return fmt.Sprintf("%s.ign", name), content.(string)
	}

	return
}

func genericConfig(d *schema.ResourceData) (filename, config string) {
	// use profile name to generate generic config filename
	name := d.Get("name").(string)

	if content, ok := d.GetOk("generic_config"); ok {
		return name, content.(string)
	}

	return
}

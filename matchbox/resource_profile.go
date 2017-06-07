package matchbox

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
			"generic": &schema.Schema{
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
	tmplName, tmpl := template(d)

	var initrds []string
	for _, initrd := range d.Get("initrd").([]interface{}) {
		initrds = append(initrds, initrd.(string))
	}
	var args []string
	for _, arg := range d.Get("args").([]interface{}) {
		args = append(args, arg.(string))
	}

	profile := &storagepb.Profile{}
	if strings.HasSuffix(tmplName, ".generic") {
		profile = &storagepb.Profile{
			Id:         name,
			GenericId:  tmplName,
			Boot: &storagepb.NetBoot{
				Kernel: d.Get("kernel").(string),
				Initrd: initrds,
				Args:   args,
			},
		}
	} else {
		profile = &storagepb.Profile{
			Id:         name,
			IgnitionId: tmplName,
			Boot: &storagepb.NetBoot{
				Kernel: d.Get("kernel").(string),
				Initrd: initrds,
				Args:   args,
			},
		}
	}
	_, err := client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
		Profile: profile,
	})
	if err != nil {
		return err
	}

	// Template Generic or Container Linux Config
	if strings.HasSuffix(tmplName, ".generic") {
		_, err = client.Generic.GenericPut(ctx, &serverpb.GenericPutRequest{
			Name:   tmplName,
			Config: []byte(tmpl),
		})
	} else {
		_, err = client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
			Name:   tmplName,
			Config: []byte(tmpl),
		})
	}
	if err != nil {
		return err
	}

	d.SetId(profile.GetId())
	return err
}

func validateResourceProfile(d *schema.ResourceData) error {
	_, hasRAW := d.GetOk("raw_ignition")
	_, hasCLC := d.GetOk("container_linux_config")
	_, hasGEN := d.GetOk("generic")
	if (hasCLC && hasRAW) || (hasCLC && hasGEN) || (hasRAW && hasGEN) {
		return errors.New("container_linux_config, raw_ignition and generic are mutually exclusive")
	}

	if !hasCLC && !hasRAW && !hasGEN {
		return errors.New("container_linux_config or raw_ignition or generic are required")
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

	tmplName := templateName(d)
	if strings.HasSuffix(tmplName, ".generic") {
		_, err = client.Generic.GenericGet(ctx, &serverpb.GenericGetRequest{
			Name: tmplName,
		})
	} else {
		_, err = client.Ignition.IgnitionGet(ctx, &serverpb.IgnitionGetRequest{
			Name: tmplName,
		})
	}
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

	tmplName := templateName(d)
	if strings.HasSuffix(tmplName, ".generic") {
		// Generic Template
		_, err = client.Generic.GenericDelete(ctx, &serverpb.GenericDeleteRequest{
			Name: tmplName,
		})
	} else {
		// Container Linux Config
		_, err = client.Ignition.IgnitionDelete(ctx, &serverpb.IgnitionDeleteRequest{
			Name: tmplName,
		})
	}
	if err != nil {
		return err
	}

	// resource can be destroyed in state
	d.SetId("")
	return nil
}

func templateName(d *schema.ResourceData) string {
	filename, _ := template(d)
	return filename
}

func template(d *schema.ResourceData) (filename, config string) {
	name := d.Get("name").(string)

	if content, ok := d.GetOk("container_linux_config"); ok {
		return fmt.Sprintf("%s.yaml.tmpl", name), content.(string)
	}

	if content, ok := d.GetOk("raw_ignition"); ok {
		return fmt.Sprintf("%s.ign", name), content.(string)
	}

	if content, ok := d.GetOk("generic"); ok {
		return fmt.Sprintf("%s.generic", name), content.(string)
	}

	return
}

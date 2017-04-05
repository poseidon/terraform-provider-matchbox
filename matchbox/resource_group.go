package matchbox

import (
	"context"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"selector": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)

	selectors := map[string]string{}
	for k, v := range d.Get("selector").(map[string]interface{}) {
		selectors[k] = v.(string)
	}

	richGroup := &storagepb.RichGroup{
		Id:       name,
		Profile:  d.Get("profile").(string),
		Selector: selectors,
		Metadata: d.Get("metadata").(map[string]interface{}),
	}
	group, err := richGroup.ToGroup()
	if err != nil {
		return err
	}

	_, err = client.Groups.GroupPut(ctx, &serverpb.GroupPutRequest{
		Group: group,
	})
	d.SetId(name)
	return err
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

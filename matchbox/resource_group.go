package matchbox

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matchbox "github.com/poseidon/matchbox/matchbox/client"
	"github.com/poseidon/matchbox/matchbox/server/serverpb"
	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"selector": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				ForceNew: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*matchbox.Client)
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
		return diag.FromErr(err)
	}

	_, err = client.Groups.GroupPut(ctx, &serverpb.GroupPutRequest{
		Group: group,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.GetId())
	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*matchbox.Client)

	name := d.Get("name").(string)
	groupGetResponse, err := client.Groups.GroupGet(ctx, &serverpb.GroupGetRequest{
		Id: name,
	})

	if err != nil {
		// resource doesn't exist anymore
		d.SetId("")
		return nil
	}

	group := groupGetResponse.Group

	var metadata map[string]string
	if err := json.Unmarshal(group.Metadata, &metadata); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("selector", group.Selector); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("profile", group.Profile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", metadata); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*matchbox.Client)

	name := d.Get("name").(string)
	_, err := client.Groups.GroupDelete(ctx, &serverpb.GroupDeleteRequest{
		Id: name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

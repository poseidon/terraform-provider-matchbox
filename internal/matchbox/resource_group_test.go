package matchbox

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
	"github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

const groupWithAllFields = `
	resource "matchbox_group" "default" {
		name    = "default"
		profile = "worker"
		selector = {
			os = "installed"
		}
		metadata = {
			user = "core"
		}
	}
`

const groupMinimal = `
	resource "matchbox_group" "default" {
		name    = "minimal"
		profile = "worker"
	}
`

func TestResourceGroup(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go func() {
		err := srv.Start()
		if err != nil {
			t.Errorf("fixture server start: %v", err)
		}
	}()
	defer srv.Stop()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: srv.AddProviderConfig(groupWithAllFields),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matchbox_group.default", "id", "default"),
					resource.TestCheckResourceAttr("matchbox_group.default", "profile", "worker"),
					resource.TestCheckResourceAttr("matchbox_group.default", "selector.%", "1"),
					resource.TestCheckResourceAttr("matchbox_group.default", "selector.os", "installed"),
					resource.TestCheckResourceAttr("matchbox_group.default", "metadata.%", "1"),
					resource.TestCheckResourceAttr("matchbox_group.default", "metadata.user", "core"),
					checkMatchboxGroup(srv, &storagepb.Group{
						Id:       "default",
						Profile:  "worker",
						Selector: map[string]string{"os": "installed"},
						Metadata: []byte(`{"user":"core"}`),
					}),
				),
			},
			{
				Config: srv.AddProviderConfig(groupMinimal),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matchbox_group.default", "id", "minimal"),
					resource.TestCheckResourceAttr("matchbox_group.default", "profile", "worker"),
					resource.TestCheckResourceAttr("matchbox_group.default", "selector.%", "0"),
					resource.TestCheckResourceAttr("matchbox_group.default", "metadata.%", "0"),
					checkMatchboxGroup(srv, &storagepb.Group{
						Id:       "minimal",
						Profile:  "worker",
						Metadata: []byte(`{}`),
					}),
				),
			},
		},
	})
}

// TestResourceGroup_Read checks the provider compares the desired state with
// the actual matchbox state
func TestResourceGroup_Read(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go func() {
		err := srv.Start()
		if err != nil {
			t.Errorf("fixture server start: %v", err)
		}
	}()
	defer srv.Stop()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: srv.AddProviderConfig(groupWithAllFields),
			},
			{
				PreConfig: func() {
					// mutate resource on matchbox server
					group, _ := srv.Store.GroupGet("default")
					group.Profile = "altered"
				},
				Config:             srv.AddProviderConfig(groupWithAllFields),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// leave selector and metadata empty
			{
				Config: srv.AddProviderConfig(groupMinimal),
			},
			{
				Config:             srv.AddProviderConfig(groupWithAllFields),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// real matchbox empty metadata is an empty []byte
			{
				PreConfig: func() {
					// mutate resource on matchbox server
					group, _ := srv.Store.GroupGet("minimal")
					group.Metadata = []byte("")
				},
				Config:   srv.AddProviderConfig(groupMinimal),
				PlanOnly: true,
			},
		},
	})
}

func checkMatchboxGroup(srv *FixtureServer, expected *storagepb.Group) resource.TestCheckFunc {
	fn := func(s *terraform.State) error {
		grp, err := srv.Store.GroupGet(expected.Id)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(grp, expected) {
			return fmt.Errorf("expected %+v, got %+v", expected, grp)
		}
		return nil
	}
	return fn
}

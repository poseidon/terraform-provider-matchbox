package matchbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

func TestResourceGroup(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go srv.Start()
	defer srv.Stop()

	hcl := `
		resource "matchbox_group" "default" {
			name    = "default"
			profile = "foo"
			selector = {
				  qux = "baz"
			}

			metadata = {
				foo = "bar"
			}
		}
	`

	check := func(s *terraform.State) error {
		grp, err := srv.Store.GroupGet("default")
		if err != nil {
			return err
		}

		if grp.GetId() != "default" {
			return fmt.Errorf("id, found %q", grp.GetId())
		}

		if grp.GetProfile() != "foo" {
			return fmt.Errorf("profile, found %q", grp.GetProfile())
		}

		selector := grp.GetSelector()
		if len(selector) != 1 || selector["qux"] != "baz" {
			return fmt.Errorf("selector.qux, found %q", selector["qux"])
		}

		if string(grp.GetMetadata()) != "{\"foo\":\"bar\"}" {
			return fmt.Errorf("metadata, found %q", grp.GetProfile())
		}

		return nil
	}

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{{
			Config: srv.AddProviderConfig(hcl),
			Check:  check,
		}},
	})
}

// TestResourceGroup_Read checks the provider compares the desired state with the actual matchbox state and not only
// the Terraform state.
func TestResourceGroup_Read(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go srv.Start()
	defer srv.Stop()

	hcl := `
		resource "matchbox_group" "default" {
			name    = "default"
			profile = "foo"
			selector = {
				  qux = "baz"
			}

			metadata = {
				foo = "bar"
			}
		}
	`

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: srv.AddProviderConfig(hcl),
			},
			{
				PreConfig: func() {
					group, _ := srv.Store.GroupGet("default")
					group.Selector["bux"] = "qux"
				},
				Config:             srv.AddProviderConfig(hcl),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

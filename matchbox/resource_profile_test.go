package matchbox

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/poseidon/matchbox/matchbox/storage/testfakes"
)

func TestResourceProfile(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go func() {
		err := srv.Start()
		if err != nil {
			t.Errorf("fixture server start: %v", err)
		}
	}()
	defer srv.Stop()

	hcl := `
		resource "matchbox_profile" "default" {
			name   = "default"
			kernel = "foo"

			initrd = [
				"bar",
			]

			args = [
				"qux",
			]

			container_linux_config = "baz"
			generic_config = "experimental"
		}
	`

	check := func(s *terraform.State) error {
		profile, err := srv.Store.ProfileGet("default")
		if err != nil {
			return err
		}

		if profile.GetId() != "default" {
			return fmt.Errorf("id, found %q", profile.GetId())
		}

		if profile.GetIgnitionId() != "default.yaml.tmpl" {
			return fmt.Errorf("profile, found %q", profile.GetIgnitionId())
		}

		boot := profile.GetBoot()
		if boot.GetKernel() != "foo" {
			return fmt.Errorf("kernel, found %s", boot.GetKernel())
		}

		initrd := boot.GetInitrd()
		if len(initrd) != 1 || initrd[0] != "bar" {
			return fmt.Errorf("kernel, found %v", initrd)
		}

		args := boot.GetArgs()
		if len(args) != 1 || args[0] != "qux" {
			return fmt.Errorf("args, found %v", args)
		}

		clc, err := srv.Store.IgnitionGet("default.yaml.tmpl")
		if err != nil {
			return fmt.Errorf("failed to get Container Linux config: %v", err)
		}
		if clc != "baz" {
			return fmt.Errorf("want Container Linux config 'baz', got %q", clc)
		}

		genericConfig, err := srv.Store.GenericGet("default")
		if err != nil {
			return fmt.Errorf("failed to get generic config: %v", err)
		}
		if genericConfig != "experimental" {
			return fmt.Errorf("want generic config 'experimental', got %s", genericConfig)
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

func TestResourceProfile_withIgnition(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go func() {
		err := srv.Start()
		if err != nil {
			t.Errorf("fixture server start: %v", err)
		}
	}()
	defer srv.Stop()

	hcl := `
		resource "matchbox_profile" "default" {
			name   = "default"
			kernel = "foo"

			initrd = [
				"bar",
			]

			args = [
				"qux",
			]

			raw_ignition = "baz"
		}
	`

	check := func(s *terraform.State) error {
		profile, err := srv.Store.ProfileGet("default")
		if err != nil {
			return err
		}

		if profile.GetIgnitionId() != "default.ign" {
			return fmt.Errorf("profile, found %q", profile.GetIgnitionId())
		}

		ignition, err := srv.Store.IgnitionGet("default.ign")
		if err != nil {
			return fmt.Errorf("failed to get raw Ignition config: %v", err)
		}
		if ignition != "baz" {
			return fmt.Errorf("want raw Ignition 'baz', got %q", ignition)
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

func TestResourceProfile_withIgnitionAndContainerLinuxConfig(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go func() {
		err := srv.Start()
		if err != nil {
			t.Errorf("fixture server start: %v", err)
		}
	}()
	defer srv.Stop()

	hcl := `
		resource "matchbox_profile" "default" {
			name   = "default"
			container_linux_config = "baz"
			raw_ignition = "baz"
		}
	`

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{{
			Config:      srv.AddProviderConfig(hcl),
			ExpectError: regexp.MustCompile("are mutually exclusive"),
		}},
	})
}

// TestResourceProfile_Read checks the provider compares the desired state with the actual matchbox state and not only
// the Terraform state.
func TestResourceProfile_Read(t *testing.T) {
	srv := NewFixtureServer(clientTLSInfo, serverTLSInfo, testfakes.NewFixedStore())
	go func() {
		err := srv.Start()
		if err != nil {
			t.Errorf("fixture server start: %v", err)
		}
	}()
	defer srv.Stop()

	hcl := `
		resource "matchbox_profile" "default" {
			name   = "default"
			kernel = "foo"

			initrd = [
				"bar",
			]

			args = [
				"qux",
			]

			raw_ignition = "baz"
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
					profile, _ := srv.Store.ProfileGet("default")
					profile.Boot.Args = append(profile.Boot.Args, "bux")

				},
				Config:             srv.AddProviderConfig(hcl),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

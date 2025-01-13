//go:build acceptance.test

package redirect_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	client := openWrtServer.LuCIRPCClient(
		ctx,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()
	options := lucirpc.Options{
		"name":      lucirpc.String("testing"),
		"target":    lucirpc.String("DNAT"),
		"src":       lucirpc.String("wan"),
		"dest":      lucirpc.String("lan"),
		"src_dport": lucirpc.Integer(22),
		"dest_port": lucirpc.Integer(22),
	}
	ok, err := client.CreateSection(ctx, "firwall", "redirect", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_redirect" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "name", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "src", "wan"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "target", "ACCEPT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_redirect.testing", "dest_port", "22"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_redirect" "testing" {
	id        = "testing"
	name      = "example-redirect"
	src       = "wan"
	src_dport = 8080
	dest      = "lan"
	dest_port = 8080
	target    = "DNAT"
	family    = "ipv4"
	proto = [
		"udp",
		"tcp",
	]
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "name", "example-redirect"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src", "wan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src_dport", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest_port", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "target", "DNAT"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "family", "ipv4"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "proto[0]", "udp"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "proto[1]", "tcp"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_firewall_redirect.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_redirect" "testing" {
	id        = "testing"
	name      = "example-redirect"
	src       = "wan"
	src_dport = 8080
	dest      = "lan"
	dest_port = 8080
	target    = "DNAT"
	family    = "any"
	proto = [
		"udp",
		"tcp",
	]
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "name", "example-redirect"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src", "wan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "src_dport", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "dest_port", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "target", "DNAT"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "family", "any"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "proto[0]", "udp"),
			resource.TestCheckResourceAttr("openwrt_firewall_redirect.testing", "proto[1]", "tcp"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}

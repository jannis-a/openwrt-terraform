//go:build acceptance.test

package rule_test

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
		"target":    lucirpc.String("ACCEPT"),
		"src":       lucirpc.String("wan"),
		"dest_port": lucirpc.Integer(22),
	}
	ok, err := client.CreateSection(ctx, "firwall", "rule", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_rule" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "name", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "src", "wan"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "target", "ACCEPT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_rule.testing", "dest_port", "22"),
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

resource "openwrt_firewall_rule" "testing" {
	id        = "testing"
	name      = "example-rule"
	src       = "wan"
	dest_port = 8080
	target    = "DROP"
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
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "name", "example-rule"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "src", "wan"),
			resource.TestCheckNoResourceAttr("openwrt_firewall_rule.testing", "dest"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "dest_port", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "target", "DROP"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "family", "ipv4"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "proto[0]", "udp"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "proto[1]", "tcp"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_firewall_rule.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_rule" "testing" {
	id        = "testing"
	name      = "example-rule"
	src       = "wan"
	dest      = "lan"
	dest_port = 8080
	target    = "ACCEPT"
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
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "name", "example-rule"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "src", "wan"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "dest", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "dest_port", "8080"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "target", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "family", "ipv4"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "proto[0]", "udp"),
			resource.TestCheckResourceAttr("openwrt_firewall_rule.testing", "proto[1]", "tcp"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}

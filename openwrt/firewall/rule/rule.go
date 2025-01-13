package rule

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/firewall/zone"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	nameAttribute            = "name"
	nameAttributeDescription = "Human readable rule name."
	nameUCIOption            = "name"

	targetAttribute            = "target"
	targetAttributeDescription = "Action to take on rule match, e.g. ACCEPT, REJECT, DROP..."
	targetUCIOption            = "target"

	destAttribute            = "dest"
	destAttributeDescription = "Rule applies to traffic entering this zone"
	destUCIOption            = "dest"

	destPortAttribute            = "dest_port"
	destPortAttributeDescription = "Rule applies to traffic targetting this port"
	destPortUCIOption            = "dest_port"

	destIpAttribute            = "dest_ip"
	destIpAttributeDescription = "Rule applies to traffic targetting these IP addresses."
	destIpUCIOption            = "dest_ip"

	srcAttribute            = "src"
	srcAttributeDescription = "Rule applies to traffic from this zone"
	srcUCIOption            = "src"

	srcPortAttribute            = "src_port"
	srcPortAttributeDescription = "Rule applies to traffic originating from this port"
	srcPortUCIOption            = "src_port"

	srcIpAttribute            = "src_ip"
	srcIpAttributeDescription = "Rule applies to traffic originating from any of these IP addresses."
	srcIpUCIOption            = "src_ip"

	familyAttribute            = "family"
	familyAttributeDescription = `Restrict the rule to a single protocol family, must be one of "ipv4" or "ipv6". Applies to both if unset.`
	familyUCIOption            = "family"
	familyIpv4                 = "ipv4"
	familyIpv6                 = "ipv6"

	protocolAttribute            = "proto"
	protocolAttributeDescription = `List of protocols this rule applies to, currently only supports "tcp" and "udp"`
	protocolUCIOption            = "proto"

	schemaDescription = "Firewall traffic rules allowing ports to pass between zones."

	uciConfig = "firewall"
	uciType   = "rule"
)

var (
	nameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       nameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetName, nameAttribute, nameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetName, nameAttribute, nameUCIOption),
	}

	targetSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       targetAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTarget, targetAttribute, targetUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTarget, targetAttribute, targetUCIOption),
		Validators:        zone.TypeValidators,
	}

	destSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDest, destAttribute, destUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDest, destAttribute, destUCIOption),
	}

	destPortSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetDestPort, destPortAttribute, destPortUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetDestPort, destPortAttribute, destPortUCIOption),
		Validators: []validator.Int64{
			int64validator.Any(),
		},
	}

	destIpSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       destIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetDestIp, destIpAttribute, destIpUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetDestIp, destIpAttribute, destIpUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	}

	srcSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSrc, srcAttribute, srcUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSrc, srcAttribute, srcUCIOption),
	}

	srcPortSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetSrcPort, srcPortAttribute, srcPortUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetSrcPort, srcPortAttribute, srcPortUCIOption),
		Validators: []validator.Int64{
			int64validator.Any(),
		},
	}

	srcIpSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       srcIpAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetSrcIp, srcIpAttribute, srcIpUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetSrcIp, srcIpAttribute, srcIpUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	}

	familySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       familyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetFamily, familyAttribute, familyUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetFamily, familyAttribute, familyUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				familyIpv6,
				familyIpv4,
			),
		},
	}

	protocolSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protocolAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetProtocol, protocolAttribute, protocolUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetProtocol, protocolAttribute, protocolUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		srcAttribute:            srcSchemaAttribute,
		srcPortAttribute:        srcPortSchemaAttribute,
		srcIpAttribute:          srcIpSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		destAttribute:           destSchemaAttribute,
		destPortAttribute:       destPortSchemaAttribute,
		destIpAttribute:         destIpSchemaAttribute,
		targetAttribute:         targetSchemaAttribute,
		nameAttribute:           nameSchemaAttribute,
		familyAttribute:         familySchemaAttribute,
		protocolAttribute:       protocolSchemaAttribute,
	}
)

func NewDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

func NewResource() resource.Resource {
	return lucirpcglue.NewResource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	Id       types.String `tfsdk:"id"`
	Src      types.String `tfsdk:"src"`
	SrcPort  types.Int64  `tfsdk:"src_port"`
	SrcIp    types.List   `tfsdk:"src_ip"`
	Dest     types.String `tfsdk:"dest"`
	DestPort types.Int64  `tfsdk:"dest_port"`
	DestIp   types.List   `tfsdk:"dest_ip"`
	Target   types.String `tfsdk:"target"`
	Name     types.String `tfsdk:"name"`
	Family   types.String `tfsdk:"family"`
	Protocol types.List   `tfsdk:"proto"`
}

func modelGetTarget(m model) types.String  { return m.Target }
func modelGetName(m model) types.String    { return m.Name }
func modelGetSrc(m model) types.String     { return m.Src }
func modelGetSrcPort(m model) types.Int64  { return m.SrcPort }
func modelGetSrcIp(m model) types.List     { return m.SrcIp }
func modelGetId(m model) types.String      { return m.Id }
func modelGetDest(m model) types.String    { return m.Dest }
func modelGetFamily(m model) types.String  { return m.Family }
func modelGetDestPort(m model) types.Int64 { return m.DestPort }
func modelGetDestIp(m model) types.List    { return m.DestIp }
func modelGetProtocol(m model) types.List  { return m.Protocol }

func modelSetSrc(m *model, value types.String)     { m.Src = value }
func modelSetSrcPort(m *model, value types.Int64)  { m.SrcPort = value }
func modelSetSrcIp(m *model, value types.List)     { m.SrcIp = value }
func modelSetDest(m *model, value types.String)    { m.Dest = value }
func modelSetId(m *model, value types.String)      { m.Id = value }
func modelSetTarget(m *model, value types.String)  { m.Target = value }
func modelSetName(m *model, value types.String)    { m.Name = value }
func modelSetFamily(m *model, value types.String)  { m.Family = value }
func modelSetDestPort(m *model, value types.Int64) { m.DestPort = value }
func modelSetDestIp(m *model, value types.List)    { m.DestIp = value }
func modelSetProtocol(m *model, value types.List)  { m.Protocol = value }

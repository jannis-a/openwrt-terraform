package networkinterface

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	bringUpOnBootAttribute            = "auto"
	bringUpOnBootAttributeDescription = "Specifies whether to bring up this interface on boot."
	bringUpOnBootUCIOption            = "auto"

	deviceAttribute            = "device"
	deviceAttributeDescription = "Name of the (physical or virtual) device. This name is what the device is known as in LuCI or the `name` field in Terraform. This is not the UCI config name."
	deviceUCIOption            = "device"

	disabledAttribute            = "disabled"
	disabledAttributeDescription = "Disables this interface."
	disabledUCIOption            = "disabled"

	dnsAttribute            = "dns"
	dnsAttributeDescription = "DNS servers"
	dnsUCIOption            = "dns"

	gatewayAttribute            = "gateway"
	gatewayAttributeDescription = "Gateway of the interface"
	gatewayUCIOption            = "gateway"

	ip6AssignAttribute            = "ip6assign"
	ip6AssignAttributeDescription = "Delegate a prefix of given length to this interface"
	ip6AssignUCIOption            = "ip6assign"

	ipAddressAttribute            = "ipaddr"
	ipAddressAttributeDescription = "IP address of the interface"
	ipAddressUCIOption            = "ipaddr"

	macAddressAttribute            = "macaddr"
	macAddressAttributeDescription = "Override the MAC Address of this interface."
	macAddressUCIOption            = "macaddr"

	metricAttribute            = "metric"
	metricAttributeDescription = "Set the metric value for this interface for routing priority."
	metricUCIOption            = "metric"

	mtuAttribute            = "mtu"
	mtuAttributeDescription = "Override the default MTU on this interface."
	mtuUCIOption            = "mtu"

	netmaskAttribute            = "netmask"
	netmaskAttributeDescription = "Netmask of the interface"
	netmaskUCIOption            = "netmask"

	peerDNSAttribute            = "peerdns"
	peerDNSAttributeDescription = "Use DHCP-provided DNS servers."
	peerDNSUCIOption            = "peerdns"

	protocolAttribute            = "proto"
	protocolAttributeDescription = `The protocol type of the interface. Currently, only "dhcp, and "static" are supported.`
	protocolDHCP                 = "dhcp"
	protocolDHCPV6               = "dhcpv6"
	protocolStatic               = "static"
	protocolUCIOption            = "proto"

	requestingAddressAttribute            = "reqaddress"
	requestingAddressAttributeDescription = `Behavior for requesting address. Can only be one of "force", "try", or "none".`
	requestingAddressForce                = "force"
	requestingAddressNone                 = "none"
	requestingAddressTry                  = "try"
	requestingAddressUCIOption            = "reqaddress"

	// The fact we can only support `"auto"` is because we haven't figured out how to represent unions.
	// Once we do,
	// we can support `"auto"`, `no`, or 0-64.
	requestingPrefixAttribute            = "reqprefix"
	requestingPrefixAttributeDescription = `Behavior for requesting prefixes. Currently, only "auto" is supported.`
	requestingPrefixAuto                 = "auto"
	requestingPrefixUCIOption            = "reqprefix"

	ipv4AddressesAttribute            = "ipv4_addresses"
	ipv4AddressesAttributeDescription = "IPv4 addresses assigned to the interface"
	ipv4AddressesUCIOption            = "ipv4-address"

	upAttribute            = "up"
	upAttributeDescription = "Whether the interface is up"
	upUCIOption            = "up"

	pendingAttribute            = "pending"
	pendingAttributeDescription = "Whether the interface is pending"
	pendingUCIOption            = "pending"

	availableAttribute            = "available"
	availableAttributeDescription = "Whether the interface is available"
	availableUCIOption            = "available"

	autostartAttribute            = "autostart"
	autostartAttributeDescription = "Whether the interface starts automatically"
	autostartUCIOption            = "autostart"

	dynamicAttribute            = "dynamic"
	dynamicAttributeDescription = "Whether the interface is dynamically created"
	dynamicUCIOption            = "dynamic"

	uptimeAttribute            = "uptime"
	uptimeAttributeDescription = "Time since the interface was brought up"
	uptimeUCIOption            = "uptime"

	l3DeviceAttribute            = "l3_device"
	l3DeviceAttributeDescription = "Name of the layer 3 device"
	l3DeviceUCIOption            = "l3_device"

	ifnameAttribute            = "ifname"
	ifnameAttributeDescription = "Name of the interface"
	ifnameUCIOption            = "ifname"

	updatedAttribute            = "updated"
	updatedAttributeDescription = "Last update time"
	updatedUCIOption            = "updated"

	dnsMetricAttribute            = "dns_metric"
	dnsMetricAttributeDescription = "DNS metric"
	dnsMetricUCIOption            = "dns_metric"

	dnsServerAttribute            = "dns_server"
	dnsServerAttributeDescription = "DNS servers provided by the interface"
	dnsServerUCIOption            = "dns-server"

	dnsSearchAttribute            = "dns_search"
	dnsSearchAttributeDescription = "DNS search domains"
	dnsSearchUCIOption            = "dns-search"

	ipv6PrefixAttribute            = "ipv6_prefix"
	ipv6PrefixAttributeDescription = "IPv6 prefixes assigned to the interface"
	ipv6PrefixUCIOption            = "ipv6-prefix"

	ipv6PrefixAssignmentAttribute            = "ipv6_prefix_assignment"
	ipv6PrefixAssignmentAttributeDescription = "IPv6 prefix assignments"
	ipv6PrefixAssignmentUCIOption            = "ipv6-prefix-assignment"

	routeAttribute            = "route"
	routeAttributeDescription = "Routes associated with the interface"
	routeUCIOption            = "route"

	errorsAttribute            = "errors"
	errorsAttributeDescription = "Errors reported by the interface"
	errorsUCIOption            = "errors"

	rxBytesAttribute            = "rx_bytes"
	rxBytesAttributeDescription = "Number of received bytes"
	rxBytesUCIOption            = "rx_bytes"

	txBytesAttribute            = "tx_bytes"
	txBytesAttributeDescription = "Number of transmitted bytes"
	txBytesUCIOption            = "tx_bytes"

	rxPacketsAttribute            = "rx_packets"
	rxPacketsAttributeDescription = "Number of received packets"
	rxPacketsUCIOption            = "rx_packets"

	txPacketsAttribute            = "tx_packets"
	txPacketsAttributeDescription = "Number of transmitted packets"
	txPacketsUCIOption            = "tx_packets"

	interfaceAttribute            = "interface"
	interfaceAttributeDescription = "Interface name"
	interfaceUCIOption            = "interface"

	schemaDescription = "A logic network."

	uciConfig = "network"
	uciType   = "interface"
)

var (
	bringUpOnBootSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bringUpOnBootAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetBringUpOnBoot, bringUpOnBootAttribute, bringUpOnBootUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetBringUpOnBoot, bringUpOnBootAttribute, bringUpOnBootUCIOption),
	}

	deviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       deviceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDevice, deviceAttribute, deviceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDevice, deviceAttribute, deviceUCIOption),
	}

	disabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       disabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDisabled, disabledAttribute, disabledUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDisabled, disabledAttribute, disabledUCIOption),
	}

	dnsSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dnsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetDNS, dnsAttribute, dnsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetDNS, dnsAttribute, dnsUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.AnyWithAllWarnings(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCPV6,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolStatic,
				),
			),
		},
	}

	gatewaySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       gatewayAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetGateway, gatewayAttribute, gatewayUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetGateway, gatewayAttribute, gatewayUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid gateway (e.g. "192.168.1.1")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	ip6AssignSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6AssignAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetIP6Assign, ip6AssignAttribute, ip6AssignUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetIP6Assign, ip6AssignAttribute, ip6AssignUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 64),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	ipAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid IP address (e.g. "192.168.3.1")`,
			),
		},
	}

	macAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       macAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMacAddress, macAddressAttribute, macAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMacAddress, macAddressAttribute, macAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$"),
				`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
			),
		},
	}

	metricSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       metricAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMetric, metricAttribute, metricUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMetric, metricAttribute, metricUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 4294967295),
		},
	}

	mtuSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtuAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMTU, mtuAttribute, mtuUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMTU, mtuAttribute, mtuUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
		},
	}

	netmaskSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       netmaskAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNetmask, netmaskAttribute, netmaskUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNetmask, netmaskAttribute, netmaskUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid netmask (e.g. "255.255.255.0")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	peerDNSSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       peerDNSAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetPeerDNS, peerDNSAttribute, peerDNSUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetPeerDNS, peerDNSAttribute, peerDNSUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.AnyBool(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCPV6,
				),
			),
		},
	}

	protocolSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protocolAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetProtocol, protocolAttribute, protocolUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetProtocol, protocolAttribute, protocolUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				protocolDHCP,
				protocolDHCPV6,
				protocolStatic,
			),
		},
	}

	requestingAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       requestingAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRequestingAddress, requestingAddressAttribute, requestingAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRequestingAddress, requestingAddressAttribute, requestingAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				requestingAddressForce,
				requestingAddressNone,
				requestingAddressTry,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolDHCPV6,
			),
		},
	}

	requestingPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       requestingPrefixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRequestingPrefix, requestingPrefixAttribute, requestingPrefixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRequestingPrefix, requestingPrefixAttribute, requestingPrefixUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				requestingPrefixAuto,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolDHCPV6,
			),
		},
	}

	ipv4AddressesSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         ipv4AddressesAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetIPv4Addresses, ipv4AddressesAttribute, ipv4AddressesUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	upSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         upAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionBool(modelSetUp, upAttribute, upUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	pendingSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         pendingAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionBool(modelSetPending, pendingAttribute, pendingUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	availableSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         availableAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionBool(modelSetAvailable, availableAttribute, availableUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	autostartSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         autostartAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionBool(modelSetAutostart, autostartAttribute, autostartUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	dynamicSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         dynamicAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionBool(modelSetDynamic, dynamicAttribute, dynamicUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	uptimeSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         uptimeAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetUptime, uptimeAttribute, uptimeUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	l3DeviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         l3DeviceAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionString(modelSetL3Device, l3DeviceAttribute, l3DeviceUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	ifnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         ifnameAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionString(modelSetIfname, ifnameAttribute, ifnameUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	updatedSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         updatedAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetUpdated, updatedAttribute, updatedUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	dnsMetricSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         dnsMetricAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetDNSMetric, dnsMetricAttribute, dnsMetricUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	dnsServerSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         dnsServerAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetDNSServer, dnsServerAttribute, dnsServerUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	dnsSearchSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         dnsSearchAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetDNSSearch, dnsSearchAttribute, dnsSearchUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	ipv6PrefixSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         ipv6PrefixAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetIPv6Prefix, ipv6PrefixAttribute, ipv6PrefixUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	ipv6PrefixAssignmentSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         ipv6PrefixAssignmentAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetIPv6PrefixAssignment, ipv6PrefixAssignmentAttribute, ipv6PrefixAssignmentUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	routeSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         routeAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetRoute, routeAttribute, routeUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	errorsSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         errorsAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionListString(modelSetErrors, errorsAttribute, errorsUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	rxBytesSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         rxBytesAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetRxBytes, rxBytesAttribute, rxBytesUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	txBytesSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         txBytesAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetTxBytes, txBytesAttribute, txBytesUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	rxPacketsSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         rxPacketsAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetRxPackets, rxPacketsAttribute, rxPacketsUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	txPacketsSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         txPacketsAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionInt64(modelSetTxPackets, txPacketsAttribute, txPacketsUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	interfaceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:         interfaceAttributeDescription,
		ReadResponse:        lucirpcglue.ReadResponseOptionString(modelSetInterface, interfaceAttribute, interfaceUCIOption),
		ResourceExistence:   lucirpcglue.ReadOnly,
		DataSourceExistence: lucirpcglue.ReadOnly,
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		bringUpOnBootAttribute:        bringUpOnBootSchemaAttribute,
		deviceAttribute:               deviceSchemaAttribute,
		disabledAttribute:             disabledSchemaAttribute,
		dnsAttribute:                  dnsSchemaAttribute,
		gatewayAttribute:              gatewaySchemaAttribute,
		ip6AssignAttribute:            ip6AssignSchemaAttribute,
		ipAddressAttribute:            ipAddressSchemaAttribute,
		macAddressAttribute:           macAddressSchemaAttribute,
		mtuAttribute:                  mtuSchemaAttribute,
		metricAttribute:               metricSchemaAttribute,
		netmaskAttribute:              netmaskSchemaAttribute,
		peerDNSAttribute:              peerDNSSchemaAttribute,
		protocolAttribute:             protocolSchemaAttribute,
		requestingAddressAttribute:    requestingAddressSchemaAttribute,
		requestingPrefixAttribute:     requestingPrefixSchemaAttribute,
		ipv4AddressesAttribute:        ipv4AddressesSchemaAttribute,
		upAttribute:                   upSchemaAttribute,
		pendingAttribute:              pendingSchemaAttribute,
		availableAttribute:            availableSchemaAttribute,
		autostartAttribute:            autostartSchemaAttribute,
		dynamicAttribute:              dynamicSchemaAttribute,
		uptimeAttribute:               uptimeSchemaAttribute,
		l3DeviceAttribute:             l3DeviceSchemaAttribute,
		ifnameAttribute:               ifnameSchemaAttribute,
		updatedAttribute:              updatedSchemaAttribute,
		dnsMetricAttribute:            dnsMetricSchemaAttribute,
		dnsServerAttribute:            dnsServerSchemaAttribute,
		dnsSearchAttribute:            dnsSearchSchemaAttribute,
		ipv6PrefixAttribute:           ipv6PrefixSchemaAttribute,
		ipv6PrefixAssignmentAttribute: ipv6PrefixAssignmentSchemaAttribute,
		routeAttribute:                routeSchemaAttribute,
		errorsAttribute:               errorsSchemaAttribute,
		rxBytesAttribute:              rxBytesSchemaAttribute,
		txBytesAttribute:              txBytesSchemaAttribute,
		rxPacketsAttribute:            rxPacketsSchemaAttribute,
		txPacketsAttribute:            txPacketsSchemaAttribute,
		interfaceAttribute:            interfaceSchemaAttribute,
		lucirpcglue.IdAttribute:       lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
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
	BringUpOnBoot        types.Bool   `tfsdk:"auto"`
	Device               types.String `tfsdk:"device"`
	Disabled             types.Bool   `tfsdk:"disabled"`
	DNS                  types.List   `tfsdk:"dns"`
	Gateway              types.String `tfsdk:"gateway"`
	Id                   types.String `tfsdk:"id"`
	IP6Assign            types.Int64  `tfsdk:"ip6assign"`
	IPAddress            types.String `tfsdk:"ipaddr"`
	MacAddress           types.String `tfsdk:"macaddr"`
	MTU                  types.Int64  `tfsdk:"mtu"`
	Netmask              types.String `tfsdk:"netmask"`
	PeerDNS              types.Bool   `tfsdk:"peerdns"`
	Protocol             types.String `tfsdk:"proto"`
	RequestingAddress    types.String `tfsdk:"reqaddress"`
	RequestingPrefix     types.String `tfsdk:"reqprefix"`
	Metric               types.Int64  `tfsdk:"metric"`
	IPv4Addresses        types.List   `tfsdk:"ipv4_addresses"`
	Up                   types.Bool   `tfsdk:"up"`
	Pending              types.Bool   `tfsdk:"pending"`
	Available            types.Bool   `tfsdk:"available"`
	Autostart            types.Bool   `tfsdk:"autostart"`
	Dynamic              types.Bool   `tfsdk:"dynamic"`
	Uptime               types.Int64  `tfsdk:"uptime"`
	L3Device             types.String `tfsdk:"l3_device"`
	Ifname               types.String `tfsdk:"ifname"`
	Updated              types.Int64  `tfsdk:"updated"`
	DNSMetric            types.Int64  `tfsdk:"dns_metric"`
	DNSServer            types.List   `tfsdk:"dns_server"`
	DNSSearch            types.List   `tfsdk:"dns_search"`
	IPv6Prefix           types.List   `tfsdk:"ipv6_prefix"`
	IPv6PrefixAssignment types.List   `tfsdk:"ipv6_prefix_assignment"`
	Route                types.List   `tfsdk:"route"`
	Data                 types.Map    `tfsdk:"data"`
	Errors               types.List   `tfsdk:"errors"`
	RxBytes              types.Int64  `tfsdk:"rx_bytes"`
	TxBytes              types.Int64  `tfsdk:"tx_bytes"`
	RxPackets            types.Int64  `tfsdk:"rx_packets"`
	TxPackets            types.Int64  `tfsdk:"tx_packets"`
	Interface            types.String `tfsdk:"interface"`
}

func modelGetMetric(m model) types.Int64              { return m.Metric }
func modelGetBringUpOnBoot(m model) types.Bool        { return m.BringUpOnBoot }
func modelGetDevice(m model) types.String             { return m.Device }
func modelGetDisabled(m model) types.Bool             { return m.Disabled }
func modelGetDNS(m model) types.List                  { return m.DNS }
func modelGetGateway(m model) types.String            { return m.Gateway }
func modelGetId(m model) types.String                 { return m.Id }
func modelGetIP6Assign(m model) types.Int64           { return m.IP6Assign }
func modelGetIPAddress(m model) types.String          { return m.IPAddress }
func modelGetMacAddress(m model) types.String         { return m.MacAddress }
func modelGetMTU(m model) types.Int64                 { return m.MTU }
func modelGetNetmask(m model) types.String            { return m.Netmask }
func modelGetPeerDNS(m model) types.Bool              { return m.PeerDNS }
func modelGetProtocol(m model) types.String           { return m.Protocol }
func modelGetRequestingAddress(m model) types.String  { return m.RequestingAddress }
func modelGetRequestingPrefix(m model) types.String   { return m.RequestingPrefix }
func modelGetIPv4Addresses(m model) types.List        { return m.IPv4Addresses }
func modelGetUp(m model) types.Bool                   { return m.Up }
func modelGetPending(m model) types.Bool              { return m.Pending }
func modelGetAvailable(m model) types.Bool            { return m.Available }
func modelGetAutostart(m model) types.Bool            { return m.Autostart }
func modelGetDynamic(m model) types.Bool              { return m.Dynamic }
func modelGetUptime(m model) types.Int64              { return m.Uptime }
func modelGetL3Device(m model) types.String           { return m.L3Device }
func modelGetIfname(m model) types.String             { return m.Ifname }
func modelGetUpdated(m model) types.Int64             { return m.Updated }
func modelGetDNSMetric(m model) types.Int64           { return m.DNSMetric }
func modelGetDNSServer(m model) types.List            { return m.DNSServer }
func modelGetDNSSearch(m model) types.List            { return m.DNSSearch }
func modelGetIPv6Prefix(m model) types.List           { return m.IPv6Prefix }
func modelGetIPv6PrefixAssignment(m model) types.List { return m.IPv6PrefixAssignment }
func modelGetRoute(m model) types.List                { return m.Route }
func modelGetData(m model) types.Map                  { return m.Data }
func modelGetErrors(m model) types.List               { return m.Errors }
func modelGetRxBytes(m model) types.Int64             { return m.RxBytes }
func modelGetTxBytes(m model) types.Int64             { return m.TxBytes }
func modelGetRxPackets(m model) types.Int64           { return m.RxPackets }
func modelGetTxPackets(m model) types.Int64           { return m.TxPackets }
func modelGetInterface(m model) types.String          { return m.Interface }

func modelSetMetric(m *model, value types.Int64)              { m.Metric = value }
func modelSetBringUpOnBoot(m *model, value types.Bool)        { m.BringUpOnBoot = value }
func modelSetDevice(m *model, value types.String)             { m.Device = value }
func modelSetDisabled(m *model, value types.Bool)             { m.Disabled = value }
func modelSetDNS(m *model, value types.List)                  { m.DNS = value }
func modelSetGateway(m *model, value types.String)            { m.Gateway = value }
func modelSetId(m *model, value types.String)                 { m.Id = value }
func modelSetIP6Assign(m *model, value types.Int64)           { m.IP6Assign = value }
func modelSetIPAddress(m *model, value types.String)          { m.IPAddress = value }
func modelSetMacAddress(m *model, value types.String)         { m.MacAddress = value }
func modelSetMTU(m *model, value types.Int64)                 { m.MTU = value }
func modelSetNetmask(m *model, value types.String)            { m.Netmask = value }
func modelSetPeerDNS(m *model, value types.Bool)              { m.PeerDNS = value }
func modelSetProtocol(m *model, value types.String)           { m.Protocol = value }
func modelSetRequestingAddress(m *model, value types.String)  { m.RequestingAddress = value }
func modelSetRequestingPrefix(m *model, value types.String)   { m.RequestingPrefix = value }
func modelSetIPv4Addresses(m *model, value types.List)        { m.IPv4Addresses = value }
func modelSetUp(m *model, value types.Bool)                   { m.Up = value }
func modelSetPending(m *model, value types.Bool)              { m.Pending = value }
func modelSetAvailable(m *model, value types.Bool)            { m.Available = value }
func modelSetAutostart(m *model, value types.Bool)            { m.Autostart = value }
func modelSetDynamic(m *model, value types.Bool)              { m.Dynamic = value }
func modelSetUptime(m *model, value types.Int64)              { m.Uptime = value }
func modelSetL3Device(m *model, value types.String)           { m.L3Device = value }
func modelSetIfname(m *model, value types.String)             { m.Ifname = value }
func modelSetUpdated(m *model, value types.Int64)             { m.Updated = value }
func modelSetDNSMetric(m *model, value types.Int64)           { m.DNSMetric = value }
func modelSetDNSServer(m *model, value types.List)            { m.DNSServer = value }
func modelSetDNSSearch(m *model, value types.List)            { m.DNSSearch = value }
func modelSetIPv6Prefix(m *model, value types.List)           { m.IPv6Prefix = value }
func modelSetIPv6PrefixAssignment(m *model, value types.List) { m.IPv6PrefixAssignment = value }
func modelSetRoute(m *model, value types.List)                { m.Route = value }
func modelSetData(m *model, value types.Map)                  { m.Data = value }
func modelSetErrors(m *model, value types.List)               { m.Errors = value }
func modelSetRxBytes(m *model, value types.Int64)             { m.RxBytes = value }
func modelSetTxBytes(m *model, value types.Int64)             { m.TxBytes = value }
func modelSetRxPackets(m *model, value types.Int64)           { m.RxPackets = value }
func modelSetTxPackets(m *model, value types.Int64)           { m.TxPackets = value }
func modelSetInterface(m *model, value types.String)          { m.Interface = value }

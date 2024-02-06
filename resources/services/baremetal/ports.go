package baremetal

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources/internal/utils"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

func Ports() *schema.Table {
	return &schema.Table{
		Name:     "openstack_baremetal_ports",
		Resolver: fetchPort,
		Transform: transformers.TransformWithStruct(
			&Port{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
	}
}

func fetchPort(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	ironic, err := api.GetServiceClient(client.BareMetalV1)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}
	opts := ports.ListOpts{}

	allPages, err := ports.List(ironic, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("opts", format.ToPrettyJSON(opts)).Msg("error listing ports with options")
		return err
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		panic(err)
	}
	for _, port := range allPorts {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		port := port
		api.Logger().Debug().Str("name", port.UUID).Msg("streaming port")
		res <- port
	}
	return nil
}

type Port struct {
	ID             string    `json:"id"`
	UUID           string    `json:"uuid"`
	TenantID       string    `json:"tenant_id"`
	ProjectID      string    `json:"project_id"`
	DeviceID       string    `json:"device_id"`
	NetworkID      string    `json:"network_id"`
	IpAddresses    *[]string `json:"ip_addresses"`
	IpAddress      string    `json:"ip_address"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	AdminStateUp   bool      `json:"admin_state_up"`
	Status         string    `json:"status"`
	Address        string    `json:"mac_address"`
	FixedIPs       *[]string `json:"fixed_ips"`
	DeviceOwner    string    `json:"device_owner"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	AllowedAddressPairs   *[]string   `json:"allowed_address_pairs"`
	Tags                  *[]string   `json:"tags"`
	PropagateUplinkStatus bool        `json:"propagate_uplink_status"`
	ValueSpecs            string      `json:"value_specs"`
	RevisionNumber        int         `json:"revision_number"`
	CreatedAt             *utils.Time `json:"created" cq-name:"created_at" cq-type:"timestamp"`
	UpdatedAt             *utils.Time `json:"updated" cq-name:"updated_at" cq-type:"timestamp"`
}

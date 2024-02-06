package networking

import (
	"context"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

func Networks() *schema.Table {
	return &schema.Table{
		Name:     "openstack_networking_networks",
		Resolver: fetchNetworks,
		Transform: transformers.TransformWithStruct(
			&Network{},
			transformers.WithPrimaryKeys("ID"),
			// transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			// transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			NetworkSubnets(),
			NetworkTags(),
		},
	}
}

func fetchNetworks(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	neutron, err := api.GetServiceClient(client.NetworkV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := networks.ListOpts{}

	allPages, err := networks.List(neutron, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing networks with options")
		return err
	}
	allNetworks := []*Network{}
	if err = networks.ExtractNetworksInto(allPages, &allNetworks); err != nil {
		api.Logger().Error().Err(err).Msg("error extracting networks")
		return err
	}
	api.Logger().Debug().Int("count", len(allNetworks)).Msg("networks retrieved")

	for _, network := range allNetworks {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		network := network
		//api.Logger().Debug().Str("data", format.ToPrettyJSON(network)).Msg("streaming network")
		api.Logger().Debug().Str("id", network.ID).Msg("streaming network")
		res <- network
	}
	return nil
}

type Network struct {
	// UUID for the network
	ID string `json:"id"`

	// Human-readable name for the network. Might not be unique.
	Name string `json:"name"`

	// Description for the network
	Description string `json:"description"`

	// The administrative state of network. If false (down), the network does not
	// forward packets.
	AdminStateUp bool `json:"admin_state_up"`

	// Indicates whether network is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional
	// values.
	Status string `json:"status"`

	// Subnets associated with this network.
	Subnets []string `json:"subnets"`

	// TenantID is the project owner of the network.
	TenantID string `json:"tenant_id"`

	// UpdatedAt and CreatedAt contain ISO-8601 timestamps of when the state of the
	// network last changed, and when it was created.
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`

	// ProjectID is the project owner of the network.
	ProjectID string `json:"project_id"`

	// Specifies whether the network resource can be accessed by any tenant.
	Shared bool `json:"shared"`

	// Availability zone hints groups network nodes that run services like DHCP, L3, FW, and others.
	// Used to make network resources highly available.
	AvailabilityZoneHints []string `json:"availability_zone_hints"`

	// Tags optionally set via extensions/attributestags
	Tags []string `json:"tags"`

	// RevisionNumber optionally set via extensions/standard-attr-revisions
	RevisionNumber int `json:"revision_number"`
}

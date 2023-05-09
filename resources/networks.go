package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

func Networks() *schema.Table {
	return &schema.Table{
		Name:     "openstack_networks",
		Resolver: fetchNetworks,
		Transform: transformers.TransformWithStruct(
			&networks.Network{},
			transformers.WithPrimaryKeys("ID"),
			// transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			// transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
	}
}

func fetchNetworks(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	neutron, err := api.GetServiceClient(client.NetworkV2)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := networks.ListOpts{}

	allPages, err := networks.List(neutron, opts).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing networks with options")
		return err
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error extracting networks")
		return err
	}
	api.Logger.Debug().Int("count", len(allNetworks)).Msg("networks retrieved")

	for _, network := range allNetworks {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		network := network
		//api.Logger.Debug().Str("data", format.ToPrettyJSON(network)).Msg("streaming network")
		api.Logger.Debug().Str("id", network.ID).Msg("streaming network")
		res <- network
	}
	return nil
}

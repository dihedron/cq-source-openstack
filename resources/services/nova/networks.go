package nova

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/networks"
)

func Networks() *schema.Table {
	return &schema.Table{
		Name:     "openstack_nova_networks",
		Resolver: fetchNetworks,
		Transform: transformers.TransformWithStruct(
			&networks.Network{},
		),
	}
}

func fetchNetworks(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.NovaV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allPages, err := networks.List(nova).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error listing networks")
		return err
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting networks")
		return err
	}
	for _, network := range allNetworks {
		res <- network
	}

	return nil
}

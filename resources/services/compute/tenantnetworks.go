package compute

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tenantnetworks"
)

func TenantNetworks() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_tenantnetworks",
		Resolver: fetchTenantNetworks,
		Transform: transformers.TransformWithStruct(
			&tenantnetworks.Network{},
		),
	}
}

func fetchTenantNetworks(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	compute, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allPages, err := tenantnetworks.List(compute).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error listing tenant networks")
		return err
	}
	allTenantNetworks, err := tenantnetworks.ExtractNetworks(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting tenant networks")
		return err
	}
	for _, network := range allTenantNetworks {
		res <- network
	}

	return nil
}

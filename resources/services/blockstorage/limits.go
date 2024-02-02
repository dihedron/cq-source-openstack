package blockstorage

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/limits"
)

func Limits() *schema.Table {
	return &schema.Table{
		Name:     "openstack_limits",
		Resolver: fetchLimits,
		Transform: transformers.TransformWithStruct(
			&limits.Limits{},
		),
	}
}

func fetchLimits(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	cinder, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allLimits, err := limits.Get(cinder).Extract()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error getting limits")
		panic(err)
	}

	res <- allLimits
	return nil
}

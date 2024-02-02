package blockstorage

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/availabilityzones"
)

func AvailabilityZones() *schema.Table {
	return &schema.Table{
		Name:     "openstack_availabilityzones",
		Resolver: fetchAvailabilityZones,
		Transform: transformers.TransformWithStruct(
			&availabilityzones.AvailabilityZone{},
		),
	}
}

func fetchAvailabilityZones(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	cinder, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allPages, err := availabilityzones.List(cinder).AllPages()
	if err != nil {
		panic(err)
	}

	allAvailabilityZones, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error getting availabilityzones")
		panic(err)
	}

	for _, zoneInfo := range allAvailabilityZones {
		res <- zoneInfo
	}

	return nil
}

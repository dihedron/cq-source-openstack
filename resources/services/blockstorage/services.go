package blockstorage

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/services"
)

func Services() *schema.Table {
	return &schema.Table{
		Name:     "openstack_blockstorage_services",
		Resolver: fetchServices,
		Transform: transformers.TransformWithStruct(
			&services.Service{},
		),
	}
}

func fetchServices(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	blockstorage, err := api.GetServiceClient(client.BlockStorageV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving blockstorage client")
		return err
	}

	opts := services.ListOpts{}

	allPages, err := services.List(blockstorage, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error listing services")
		return err
	}
	allServices, err := services.ExtractServices(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting services")
		return err
	}
	api.Logger().Debug().Int("count", len(allServices)).Msg("services retrieved")

	for _, service := range allServices {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		service := service
		res <- service
	}

	return nil
}

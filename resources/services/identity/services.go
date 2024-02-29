package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
)

func Services() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_services",
		Resolver: fetchServices,
		Transform: transformers.TransformWithStruct(
			&services.Service{},
			transformers.WithSkipFields("Links", "Extra"),
		),
	}
}

func fetchServices(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	service_types := [5]string{"blockstorage", "compute", "identity", "imageservice", "networking"}
	result := []services.Service{}
	for _, service_type := range service_types {
		service_type := service_type
		opts := services.ListOpts{
			ServiceType: service_type,
		}

		allPages, err := services.List(identity, opts).AllPages()
		if err != nil {
			api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing services with options")
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
			api.Logger().Debug().Str("service id", service.ID).Str("data", format.ToJSON(service)).Msg("streaming service")
			result = append(result, service)
		}
	}
	res <- result
	return nil
}

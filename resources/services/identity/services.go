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

	opts := services.ListOpts{
		// ServiceType: catalog.Type,
	}

	allPages, err := services.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing services with options")
		return err
	}

	allService, err := services.ExtractServices(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting services")
		return err
	}
	api.Logger().Debug().Int("count", len(allService)).Msg("services retrieved")

	for _, service := range allService {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("service id", service.ID).Str("data", format.ToJSON(service)).Msg("streaming service")
		res <- service
	}
	return nil
}

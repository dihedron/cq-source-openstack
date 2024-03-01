package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/registeredlimits"
)

func RegisteredLimits() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_registeredlimits",
		Resolver: fetchRegisteredLimits,
		Transform: transformers.TransformWithStruct(
			&registeredlimits.RegisteredLimit{},
			transformers.WithSkipFields("Links"),
		),
	}
}

func fetchRegisteredLimits(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger().Debug().Msg("getting list of registered_limits...")

	opts := registeredlimits.ListOpts{}

	allPages, err := registeredlimits.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing registered_limits with options")
		return err
	}

	allLimits, err := registeredlimits.ExtractRegisteredLimits(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting registeredlimits")
		return err
	}
	api.Logger().Debug().Int("count", len(allLimits)).Msg("registeredlimits retrieved")

	for _, registered_limit := range allLimits {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("registered_limit id", registered_limit.ID).Str("data", format.ToJSON(registered_limit)).Msg("streaming registered_limit")
		res <- registered_limit
	}
	return nil
}

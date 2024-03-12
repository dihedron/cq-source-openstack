package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
)

func Roles() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_roles",
		Resolver: fetchRoles,
		Transform: transformers.TransformWithStruct(
			&roles.Role{},
			transformers.WithSkipFields("Links", "Extra"),
		),
	}
}

func fetchRoles(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger().Debug().Msg("getting list of roles...")

	opts := roles.ListOpts{}

	allPages, err := roles.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing roles with options")
		return err
	}

	allRegions, err := roles.ExtractRoles(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting roles")
		return err
	}
	api.Logger().Debug().Int("count", len(allRegions)).Msg("roles retrieved")

	for _, role := range allRegions {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("role id", role.ID).Str("data", format.ToJSON(role)).Msg("streaming role")
		res <- role
	}
	return nil
}

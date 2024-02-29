package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v2/tenants"
)

func Tenants() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_tenants",
		Resolver: fetchTenants,
		Transform: transformers.TransformWithStruct(
			&tenants.Tenant{},
			transformers.WithSkipFields("Links", "Options"),
		),
	}
}

func fetchTenants(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger().Debug().Msg("getting list of tenants...")

	opts := &tenants.ListOpts{}

	allPages, err := tenants.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing tenants with options")
		return err
	}

	allTenants, err := tenants.ExtractTenants(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting tenants")
		return err
	}
	api.Logger().Debug().Int("count", len(allTenants)).Msg("users retrieved")

	for _, tenant := range allTenants {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("tenant id", tenant.ID).Str("data", format.ToJSON(tenant)).Msg("streaming tenant")
		res <- tenant
	}
	return nil
}

package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
)

func Regions() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_regions",
		Resolver: fetchRegions,
		Transform: transformers.TransformWithStruct(
			&regions.Region{},
			transformers.WithSkipFields("Links", "Extra"),
		),
	}
}

func fetchRegions(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger().Debug().Msg("getting list of regions...")

	opts := regions.ListOpts{}

	allPages, err := regions.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing regions with options")
		return err
	}

	allRegions, err := regions.ExtractRegions(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting regions")
		return err
	}
	api.Logger().Debug().Int("count", len(allRegions)).Msg("regions retrieved")

	for _, region := range allRegions {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("region id", region.ID).Str("data", format.ToJSON(region)).Msg("streaming region")
		res <- region
	}
	return nil
}

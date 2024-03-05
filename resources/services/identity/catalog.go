package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/catalog"
)

func Catalog() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_catalog",
		Resolver: fetchCatalog,
		Transform: transformers.TransformWithStruct(
			&tokens.CatalogEntry{},
		),
	}
}

func fetchCatalog(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger().Debug().Msg("getting list of catalog services...")

	allPages, err := catalog.List(identity).AllPages()
	if err != nil {
		api.Logger().Err(err).Msg("error listing catalog services")
		return err
	}

	allCatalogs, err := catalog.ExtractServiceCatalog(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting catalogs")
		return err
	}
	api.Logger().Debug().Int("count", len(allCatalogs)).Msg("catalogs retrieved")

	for _, catalog := range allCatalogs {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		api.Logger().Debug().Str("catalog id", catalog.ID).Str("data", format.ToJSON(catalog)).Msg("streaming catalog")
		res <- catalog
	}
	return nil
}

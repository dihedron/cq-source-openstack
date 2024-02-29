package identity

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
)

func Domains() *schema.Table {
	return &schema.Table{
		Name:     "openstack_identity_domains",
		Resolver: fetchDomains,
		Transform: transformers.TransformWithStruct(
			&domains.Domain{},
			transformers.WithSkipFields("Links"),
		),
	}
}

func fetchDomains(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	identity, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger().Debug().Msg("getting list of domains...")

	opts := domains.ListOpts{}

	allPages, err := domains.List(identity, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing domains with options")
		return err
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting domains")
		return err
	}
	api.Logger().Debug().Int("count", len(allDomains)).Msg("domains retrieved")

	for _, domain := range allDomains {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		domain := domain
		api.Logger().Debug().Str("domain id", domain.ID).Str("data", format.ToJSON(domain)).Msg("streaming domain")
		res <- domain
	}
	return nil
}

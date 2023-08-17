package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

func Aggregates() *schema.Table {
	return &schema.Table{
		Name:     "openstack_aggregates",
		Resolver: fetchAggregates,
		Transform: transformers.TransformWithStruct(
			&aggregates.Aggregate{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Relations: []*schema.Table{
			AggregateHosts(),
		},
	}
}

func fetchAggregates(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	compute, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allPages, err := aggregates.List(compute).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error listing aggregates")
		return err
	}
	allAggregates, err := aggregates.ExtractAggregates(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting aggregates")
		return err
	}
	api.Logger().Debug().Int("count", len(allAggregates)).Msg("aggregates retrieved")

	for _, aggregate := range allAggregates {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		aggregate := aggregate
		// api.Logger().Debug().Str("data", format.ToPrettyJSON(aggregate)).Msg("streaming aggregate")
		api.Logger().Debug().Int("id", aggregate.ID).Msg("streaming aggregate")
		res <- aggregate
	}
	return nil
}

package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

func AggregateHosts() *schema.Table {
	return &schema.Table{
		Name:     "openstack_aggregate_hosts",
		Resolver: fetchAggregateHosts,
		Transform: transformers.TransformWithStruct(
			&Single[string]{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchAggregateHosts(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)
	aggregate := parent.Item.(aggregates.Aggregate)
	for _, v := range aggregate.Hosts {
		host := &Single[string]{Name: v}
		api.Logger.Debug().Int("aggregate id", aggregate.ID).Str("host", v).Msg("streaming aggregate host")
		res <- host
	}
	return nil
}

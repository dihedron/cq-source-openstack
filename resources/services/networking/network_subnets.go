package networking

import (
	"context"

	"github.com/dihedron/cq-plugin-utils/utils"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func NetworkSubnets() *schema.Table {
	return &schema.Table{
		Name:     "openstack_networking_network_subnets",
		Resolver: fetchNetworkSubnets,
		Transform: transformers.TransformWithStruct(
			&utils.Single[string]{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchNetworkSubnets(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)
	network := parent.Item.(*Network)
	for _, v := range network.Subnets {
		subnet := &utils.Single[string]{Name: v}
		api.Logger().Debug().Str("network id", network.ID).Msg("streaming subnet")
		res <- subnet
	}
	return nil
}

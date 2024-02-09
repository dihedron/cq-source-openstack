package baremetal

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
)

func Nodes() *schema.Table {
	return &schema.Table{
		Name:     "openstack_baremetal_nodes",
		Resolver: fetchNode,
		Transform: transformers.TransformWithStruct(
			&nodes.Node{},
			transformers.WithSkipFields("Properties"),
		),
	}
}

func fetchNode(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	ironic, err := api.GetServiceClient(client.BareMetalV1)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}
	opts := nodes.ListOpts{
		ProvisionState: nodes.Deploying,
		Fields:         []string{"name"},
	}

	allPages, err := nodes.List(ironic, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("opts", format.ToPrettyJSON(opts)).Msg("error listing nodes with options")
		return err
	}

	allNodes, err := nodes.ExtractNodes(allPages)
	if err != nil {
		api.Logger().Err(err).Msg("error extracting nodes")
		return err
	}
	for _, node := range allNodes {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		node := node
		api.Logger().Debug().Str("name", node.Name).Msg("streaming node")
		res <- node
	}
	return nil
}

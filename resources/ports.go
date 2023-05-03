package resources

import (
	"context"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func Ports() *schema.Table {
	return &schema.Table{
		Name:     "openstack_ports",
		Resolver: fetchPorts,
		Transform: transformers.TransformWithStruct(
			&ports.Port{},
			transformers.WithPrimaryKeys("ID"),
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links"),
		),
		Columns: []schema.Column{
			{
				Name:        "ip_addresses",
				Type:        schema.TypeStringArray,
				Description: "The collection of IP addresses associated with the port.",
				Resolver: transform.Apply(
					transform.OnObjectField("FixedIPs.IPAddress"),
					transform.NilIfZero(),
				),
			},
			{
				Name:        "ip_address",
				Type:        schema.TypeString,
				Description: "The first IP address associated with the port.",
				Resolver: transform.Apply(
					transform.OnObjectField("FixedIPs.IPAddress"),
					transform.GetElementAt(0),
					transform.NilIfZero(),
				),
			},
		},
	}
}

func fetchPorts(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	neutron, err := api.GetServiceClient(client.NetworkV2)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := ports.ListOpts{}

	allPages, err := ports.List(neutron, opts).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing ports with options")
		return err
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error extracting ports")
		return err
	}
	api.Logger.Debug().Int("count", len(allPorts)).Msg("ports retrieved")

	for _, port := range allPorts {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		port := port
		api.Logger.Debug().Str("data", format.ToPrettyJSON(port)).Msg("streaming port")
		res <- port
	}
	return nil
}

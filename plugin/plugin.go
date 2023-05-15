package plugin

import (
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources"
)

var (
	Version = "development"
)

func Plugin() *source.Plugin {
	return source.NewPlugin(
		"github.com/dihedron-openstack",
		Version,
		schema.Tables{
			resources.Aggregates(),
			resources.Attachments(),
			resources.Flavors(),
			resources.Hypervisors(),
			resources.Instances(),
			resources.Images(),
			resources.KeyPairs(),
			resources.Networks(),
			resources.Ports(),
			resources.Projects(),
			resources.SecurityGroups(),
			resources.SecurityGroupRules(),
			resources.Users(),
			resources.Volumes(),
		},
		client.New,
	)
}

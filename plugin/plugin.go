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
			resources.Attachments(),
			resources.Instances(),
			resources.Images(),
			resources.Networks(),
			resources.Ports(),
			resources.Projects(),
			resources.Users(),
			resources.Volumes(),
		},
		client.New,
	)
}

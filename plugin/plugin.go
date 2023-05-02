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
		//"github.com/dihedron/openstack",
		"github.com/dihedron-openstack",
		Version,
		schema.Tables{
			resources.Instances(),
			resources.Images(),
		},
		client.New,
	)
}

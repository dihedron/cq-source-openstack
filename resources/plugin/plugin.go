package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	internalPlugin "github.com/dihedron/cq-source-openstack/plugin"
)

func Plugin() *plugin.Plugin {
	return plugin.NewPlugin(
		internalPlugin.Name,
		internalPlugin.Version,
		Configure,
	)
}
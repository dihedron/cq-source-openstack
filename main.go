package main

import (
	"github.com/cloudquery/plugin-sdk/serve"
	"github.com/dihedron/cq-source-openstack/plugin"
)

func main() {
	serve.Source(plugin.Plugin())
}

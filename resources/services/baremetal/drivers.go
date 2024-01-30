package baremetal

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/drivers"
)

func Drivers() *schema.Table {
	return &schema.Table{
		Name:     "openstack_drivers",
		Resolver: fetchDriver,
		Transform: transformers.TransformWithStruct(
			&drivers.Driver{},
			transformers.WithSkipFields("Links"),
			transformers.WithSkipFields("Properties"),
		),
	}
}

func fetchDriver(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	ironic, err := api.GetServiceClient(client.BareMetalV1)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}
	opts := drivers.ListDriversOpts{}

	allPages, err := drivers.ListDrivers(ironic, opts).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Str("opts", format.ToPrettyJSON(opts)).Msg("error listing drivers with options")
		return err
	}

	allDrivers, err := drivers.ExtractDrivers(allPages)
	if err != nil{
		panic(err)
	}
	for _, driver := range allDrivers {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		driver := driver
		api.Logger().Debug().Str("name", driver.Name).Msg("streaming driver")
		res <- driver
	}
	return nil
}

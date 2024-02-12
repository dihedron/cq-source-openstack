package nova

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/pointer"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
)

func InstanceAddresses() *schema.Table {
	return &schema.Table{
		Name:     "openstack_nova_instance_addresses",
		Resolver: fetchInstanceAddresses,
		Transform: transformers.TransformWithStruct(
			&Address{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			//transformers.WithSkipFields("OriginalName", "ExtraSpecs"),
		),
	}
}

func fetchInstanceAddresses(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	instance := parent.Item.(*Instance)

	for network, addresses := range instance.Addresses {
		for _, address := range addresses {
			address.Network = pointer.To(network)
			//api.Logger().Debug().Str("instance id", instance.ID).Str("json", format.ToPrettyJSON(address)).Msg("streaming instance addresses")
			api.Logger().Debug().Str("instance id", instance.ID).Msg("streaming instance addresses")
			res <- address
		}
	}

	return nil
}

type Address struct {
	Network    *string `json:"-" cq-name:"network"`
	MACAddress string  `json:"OS-EXT-IPS-MAC:mac_addr" cq-name:"mac_address"`
	IPType     string  `json:"OS-EXT-IPS:type" cq-name:"type"`
	IPAddress  string  `json:"addr" cq-name:"ip_address"`
	IPVersion  int     `json:"version" cq-name:"ip_version"`
}

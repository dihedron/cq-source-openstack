package compute

import (
	"context"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
)

func SecGroups() *schema.Table {
	return &schema.Table{
		Name:     "openstack_compute_secgroups",
		Resolver: fetchSecgroups,
		Transform: transformers.TransformWithStruct(
			&SecurityGroup{},
			transformers.WithSkipFields("Rule"),
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
		),
		Columns: []schema.Column{
			{
				Name:        "rule_id",
				Type:        arrow.BinaryTypes.String,
				Description: "The unique ID of the rule.",
				Resolver:    schema.PathResolver("Rule.ID"),
			},
			{
				Name:        "from_port",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The start port of the rule.",
				Resolver:    schema.PathResolver("Rule.FromPort"),
			},
			{
				Name:        "to_port",
				Type:        arrow.PrimitiveTypes.Int64,
				Description: "The end port of the rule.",
				Resolver:    schema.PathResolver("Rule.ToPort"),
			},
			{
				Name:        "ip_protocol",
				Type:        arrow.BinaryTypes.String,
				Description: "The IP protocol of the rule.",
				Resolver:    schema.PathResolver("Rule.IPProtocol"),
			},
			{
				Name:        "ip_range",
				Type:        arrow.BinaryTypes.String,
				Description: "The IP range of the rule.",
				Resolver:    schema.PathResolver("Rule.IPRange.CIDR"),
			},
			{
				Name:        "group_id",
				Type:        arrow.BinaryTypes.String,
				Description: "The group ID of the rule.",
				Resolver:    schema.PathResolver("Group.TenantID"),
			},
			{
				Name:        "group_name",
				Type:        arrow.BinaryTypes.String,
				Description: "The group name of the rule.",
				Resolver:    schema.PathResolver("Group.Name"),
			},
		},
	}
}

func fetchSecgroups(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	allPages, err := secgroups.List(nova).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving secgroups")
		return err
	}

	allSecgroups, err := secgroups.ExtractSecurityGroups(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving quotasets")
		return err
	}

	for _, secgroup := range allSecgroups {
		api.Logger().Debug().Str("secgroup id", secgroup.ID).Msg("streaming secgroup")
		res <- secgroup
	}

	return nil
}

type SecurityGroup struct {
	ID          string `json:"id" cq-name:"security_group_id"`
	Name        string `json:"name" cq-name:"security_group_name"`
	Description string `json:"description"`
	Rule        []struct {
		ID         string `json:"id" cq-name:"rule_id"`
		FromPort   int    `json:"from_port"`
		ToPort     int    `json:"to_port"`
		IPProtocol string `json:"ip_protocol"`
		IPRange    struct {
			CIDR string `json:"cidr"`
		} `json:"ip_range"`
		ParentGroupID string `json:"parent_group_id"`
		Group         struct {
			TenantID string `json:"tenant_id" cq-name:"group_id"`
			Name     string `json:"name" cq-name:"group_name"`
		}
	} `json:"rules"`
	TenantID string `json:"tenant_id" cq-name:"project_id" cq-type:"tenant"`
}

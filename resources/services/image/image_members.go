package image

import (
	"context"

	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/dihedron/cq-source-openstack/resources/internal/utils"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/members"
)

func ImageMembers() *schema.Table {
	return &schema.Table{
		Name:     "openstack_image_image_members",
		Resolver: fetchImageMembers,
		Transform: transformers.TransformWithStruct(
			&Member{},
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
		),
	}
}

func fetchImageMembers(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	image, err := api.GetServiceClient(client.ImageV2)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error retrieving client")
		return err
	}

	image := parent.Item.(images.Image)

	imageID := image.ID
	allPages, err := members.List(image, imageID).AllPages()
	if err != nil {
		api.Logger().Error().Err(err).Msg("error listing image members")
		return err
	}

	allMembers, err := members.ExtractMembers(allPages)
	if err != nil {
		api.Logger().Error().Err(err).Msg("error extracting image members")
		return err
	}
	api.Logger().Debug().Int("count", len(allMembers)).Msg("image members retrieved")

	for _, member := range allMembers {
		if ctx.Err() != nil {
			api.Logger().Debug().Msg("context done, exit")
			break
		}
		member := member
		api.Logger().Debug().Str("member id", member.MemberID).Str("data", format.ToJSON(member)).Msg("streaming image member")
		res <- member
	}

	return nil
}

type Member struct {
	CreatedAt *utils.Time `json:"created_at" cq-type:"timestamp"`
	ImageID   string      `json:"image_id"`
	MemberID  string      `json:"member_id"`
	Schema    string      `json:"schema"`
	Status    string      `json:"status"`
	UpdatedAt *utils.Time `json:"updated_at" cq-type:"timestamp"`
}

package resources

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-plugin-utils/transform"
	"github.com/dihedron/cq-source-openstack/client"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

func Flavors() *schema.Table {
	return &schema.Table{
		Name:     "openstack_flavors",
		Resolver: fetchFlavors,
		Transform: transformers.TransformWithStruct(
			&Flavor{},
			transformers.WithNameTransformer(transform.TagNameTransformer), // use cq-name tags to translate name
			transformers.WithTypeTransformer(transform.TagTypeTransformer), // use cq-type tags to translate type
			transformers.WithSkipFields("Links", "ExtraSpecs"),
		),
		Relations: []*schema.Table{
			FlavorAccesses(),
		},
	}
}

func fetchFlavors(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	nova, err := api.GetServiceClient(client.ComputeV2)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving client")
		return err
	}

	opts := flavors.ListOpts{
		AccessType: "None",
	}

	allPages, err := flavors.ListDetail(nova, opts).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing flavors with options")
		return err
	}

	allFlavors := []*Flavor{}
	if err = ExtractFlavorsInto(allPages, &allFlavors); err != nil {
		api.Logger.Error().Err(err).Msg("error extracting flavors")
		return err
	}
	api.Logger.Debug().Int("count", len(allFlavors)).Msg("flavors retrieved")
	for _, flavor := range allFlavors {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		flavor := flavor

		// retrieve the extra specs
		extraSpecs := flavors.ListExtraSpecs(nova, *flavor.ID)
		if err != nil {
			api.Logger.Error().Err(err).Str("flavor id", *flavor.ID).Msg("error getting flavor extra specs")
			return err
		}
		flavor.ExtraSpecs, err = ExtractFlavorsExtraSpec(extraSpecs, nil)
		if err != nil {
			api.Logger.Error().Err(err).Str("flavor id", *flavor.ID).Msg("error parsing flavor extra specs")
			return err
		}
		api.Logger.Debug().Str("id", *flavor.ID).Msg("streaming flavor with extra specs")
		res <- flavor
	}
	return nil
}

type Flavor struct {
	ID          *string           `json:"id,omitempty"`
	Disk        int               `json:"disk"`
	RAM         int               `json:"ram"`
	Name        string            `json:"name"`
	RxTxFactor  *float64          `json:"rxtx_factor"`
	Swap        int               `json:"-"`
	VCPUs       int               `json:"vcpus"`
	IsPublic    *bool             `json:"os-flavor-access:is_public" cq-name:"is_public"`
	Ephemeral   int               `json:"OS-FLV-EXT-DATA:ephemeral" cq-name:"ephemeral"`
	Description *string           `json:"description"` // new in version 2.55
	ExtraSpecs  *FlavorExtraSpecs `json:"extra_specs"`
}

type FlavorExtraSpecs struct {
	CPUCores        string `json:"hw:cpu_cores"`
	CPUSockets      string `json:"hw:cpu_sockets"`
	RNGAllowed      string `json:"hw_rng:allowed"`
	WatchdogAction  string `json:"hw:watchdog_action"`
	VGPUs           string `json:"resources:VGPU"`
	TraitCustomVGPU string `json:"trait:CUSTOM_VGPU"`
}

func (r *Flavor) UnmarshalJSON(b []byte) error {
	type tmp Flavor
	var s struct {
		tmp
		Swap interface{} `json:"swap"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Flavor(s.tmp)

	switch t := s.Swap.(type) {
	case float64:
		r.Swap = int(t)
	case string:
		switch t {
		case "":
			r.Swap = 0
		default:
			swap, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
			r.Swap = int(swap)
		}
	}

	return nil
}

func ExtractFlavorsInto(r pagination.Page, v interface{}) error {
	return r.(flavors.FlavorPage).Result.ExtractIntoSlicePtr(v, "flavors")
}

// Extract interprets any extraSpecsResult as ExtraSpecs, if possible.
func ExtractFlavorsExtraSpec(r flavors.ListExtraSpecsResult, v interface{}) (*FlavorExtraSpecs, error) {
	var s struct {
		ExtraSpecs *FlavorExtraSpecs `json:"extra_specs"`
	}
	err := r.ExtractInto(&s)
	return s.ExtraSpecs, err
}

// func ExtractExtraSpecsInto(r pagination.Page, v interface{}) {
// 	var s struct {
// 		ExtraSpecs map[string]string `json:"extra_specs"`
// 	}
// 	err := r.ExtractInto(&s)
// 	return s.ExtraSpecs, err
// }

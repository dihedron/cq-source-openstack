package resources

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"github.com/dihedron/cq-plugin-utils/format"
	"github.com/dihedron/cq-source-openstack/client"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
)

func Users() *schema.Table {
	return &schema.Table{
		Name:     "openstack_users",
		Resolver: fetchUsers,
		Transform: transformers.TransformWithStruct(
			&User{},
			transformers.WithSkipFields("Links", "Options"),
		),
		Relations: []*schema.Table{
			UserKeyPairs(),
		},
		Columns: []schema.Column{
			{
				Name:        "ignore_change_password_upon_first_use",
				Type:        schema.TypeBool,
				Description: "Whether the password should not changed upon first use.",
				Resolver:    schema.PathResolver("Options.IgnoreChangePasswordUponFirstUse"),
			},
			{
				Name:        "ignore_lockout_failure_attempts",
				Type:        schema.TypeBool,
				Description: "Whether the failure attempts should ignored.",
				Resolver:    schema.PathResolver("Options.IgnoreLockoutFailureAttempts"),
			},
			{
				Name:        "ignore_password_expiry",
				Type:        schema.TypeBool,
				Description: "Whether the password expiry should be ignored.",
				Resolver:    schema.PathResolver("Options.IgnorePasswordExpiry"),
			},
		},
	}
}

func fetchUsers(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {

	api := meta.(*client.Client)

	keystone, err := api.GetServiceClient(client.IdentityV3)
	if err != nil {
		api.Logger.Error().Err(err).Msg("error retrieving client")
		return err
	}

	api.Logger.Debug().Msg("getting list of users...")

	opts := users.ListOpts{}

	allPages, err := users.List(keystone, opts).AllPages()
	if err != nil {
		api.Logger.Error().Err(err).Str("options", format.ToPrettyJSON(opts)).Msg("error listing users with options")
		return err
	}

	api.Logger.Debug().Msg("list of users' pages retrieved...")

	allUsers := []*User{}
	if err = ExtractUsersInto(allPages, &allUsers); err != nil {
		api.Logger.Error().Err(err).Msg("error extracting users")
		return err
	}
	api.Logger.Debug().Int("count", len(allUsers)).Msg("instances retrieved")

	if err != nil {
		api.Logger.Error().Err(err).Msg("error extracting users")
		return err
	}
	api.Logger.Debug().Int("count", len(allUsers)).Msg("users retrieved")

	for _, user := range allUsers {
		if ctx.Err() != nil {
			api.Logger.Debug().Msg("context done, exit")
			break
		}
		user := user
		api.Logger.Debug().Str("user id", user.ID).Str("data", format.ToJSON(user)).Msg("streaming user")
		res <- user
	}
	return nil
}

type User struct {
	DefaultProjectID  string                 `json:"default_project_id"`
	Description       string                 `json:"description"`
	DomainID          string                 `json:"domain_id"`
	Enabled           bool                   `json:"enabled"`
	Extra             map[string]interface{} `json:"-"`
	ID                string                 `json:"id"`
	Links             map[string]interface{} `json:"links"`
	Name              string                 `json:"name"`
	PasswordExpiresAt time.Time              `json:"-"`
	Options           struct {
		IgnoreChangePasswordUponFirstUse bool `json:"ignore_change_password_upon_first_use"`
		IgnoreLockoutFailureAttempts     bool `json:"ignore_lockout_failure_attempts"`
		IgnorePasswordExpiry             bool `json:"ignore_password_expiry"`
	} `json:"options"`
	// Options           map[string]interface{} `json:"options"`
}

func (r *User) UnmarshalJSON(b []byte) error {
	type tmp User
	var s struct {
		tmp
		Extra             map[string]interface{}          `json:"extra"`
		PasswordExpiresAt gophercloud.JSONRFC3339MilliNoZ `json:"password_expires_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = User(s.tmp)

	r.PasswordExpiresAt = time.Time(s.PasswordExpiresAt)

	// Collect other fields and bundle them into Extra
	// but only if a field titled "extra" wasn't sent.
	if s.Extra != nil {
		r.Extra = s.Extra
	} else {
		var result interface{}
		err := json.Unmarshal(b, &result)
		if err != nil {
			return err
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			delete(resultMap, "password_expires_at")
			r.Extra = gophercloud.RemainingKeys(User{}, resultMap)
		}
	}

	return err
}

func ExtractUsersInto(r pagination.Page, v interface{}) error {
	return r.(users.UserPage).Result.ExtractIntoSlicePtr(v, "users")
}

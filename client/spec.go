package client

import (
	"net/http"
	"net/url"

	"github.com/gophercloud/gophercloud"
	"github.com/rs/zerolog/log"
)


type Spec struct {
	EndpointUrl                *string  `json:"endpoint_url,omitempty" yaml:"endpoint_url,omitempty"`
	UserID                     *string  `json:"userid,omitempty" yaml:"userid,omitempty"`
	Username                   *string  `json:"username,omitempty" yaml:"username,omitempty"`
	Password                   *string  `json:"password,omitempty" yaml:"password,omitempty"`
	Region                     *string  `json:"region,omitempty" yaml:"region,omitempty"`
	ProjectID                  *string  `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	ProjectName                *string  `json:"project_name,omitempty" yaml:"project_name,omitempty"`
	DomainID                   *string  `json:"domain_id,omitempty" yaml:"domain_id,omitempty"`
	DomainName                 *string  `json:"domain_name,omitempty" yaml:"domain_name,omitempty"`
	AccessToken                *string  `json:"access_token,omitempty" yaml:"access_token,omitempty"`
	AppCredentialID            *string  `json:"app_credential_id,omitempty" yaml:"app_credential_id,omitempty"`
	AppCredentialSecret        *string  `json:"app_credential_secret,omitempty" yaml:"app_credential_secret,omitempty"`
	AllowReauth                *bool    `json:"allow_reauth,omitempty" yaml:"allow_reauth,omitempty"`
	BareMetalV1Microversion    *string  `json:"baremetal_v1_microversion,omitempty" yaml:"baremetal_v1_microversion,omitempty"`
	IdentityV3Microversion     *string  `json:"keyston_v3_microversion,omitempty" yaml:"keyston_v3_microversion,omitempty"`
	ComputeV2Microversion      *string  `json:"compute_v2_microversion,omitempty" yaml:"compute_v2_microversion,omitempty"`
	NetworkingV2Microversion   *string  `json:"networking_v2_microversion,omitempty" yaml:"networking_v2_microversion,omitempty"`
	BlockStorageV3Microversion *string  `json:"blockstorage_v3_microversion,omitempty" yaml:"blockstorage_v3_microversion,omitempty"`
	ImageV2Microversion        *string  `json:"image_v2_microversion,omitempty" yaml:"image_v2_microversion,omitempty"`
	IncludedTables             []string `json:"included_tables,omitempty" yaml:"included_tables,omitempty"`
	ExcludedTables             []string `json:"excluded_tables,omitempty" yaml:"excluded_tables,omitempty"`
}

func (s *Spec) AssignValues() (gophercloud.AuthOptions, error) {
	auth := gophercloud.AuthOptions{}

	if s.EndpointUrl != nil {
		auth.IdentityEndpoint = *s.EndpointUrl
	}
	if s.UserID != nil {
		auth.UserID = *s.UserID
	}
	if s.Username != nil {
		auth.Username = *s.Username
	}
	if s.Password != nil {
		auth.Password = *s.Password
	}
	if s.ProjectID != nil {
		auth.TenantID = *s.ProjectID
	}
	if s.ProjectName != nil {
		auth.TenantName = *s.ProjectName
	}
	if s.DomainID != nil {
		auth.DomainID = *s.DomainID
	}
	if s.DomainName != nil {
		auth.DomainName = *s.DomainName
	}
	if s.AccessToken != nil {
		auth.TokenID = *s.AccessToken
	}
	if s.AppCredentialID != nil {
		auth.ApplicationCredentialID = *s.AppCredentialID
	}
	if s.AppCredentialSecret != nil {
		auth.ApplicationCredentialSecret = *s.AppCredentialSecret
	}
	if s.AllowReauth != nil {
		auth.AllowReauth = *s.AllowReauth
	} else {
		auth.AllowReauth = true
	}

	_, err := auth.ToTokenV3CreateMap(nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize auth options.")
		return auth, err
	}

	return auth, nil
}

func (s *Spec) Validate() error {
	if s.EndpointUrl == nil {
		log.Error().Msg("missing endpoint URL")
		return nil
	}
	// Check that the endpoint URL is a valid URL
	_, err := url.ParseRequestURI(*s.EndpointUrl)
	if err != nil {
		log.Error().Err(err).Msg("invalid endpoint URL.")
		return err
	}
	// Check that the endpoint URL is reachable
	_, err = http.Get(*s.EndpointUrl)
	if err != nil {
		log.Error().Err(err).Msg("unreachable endpoint URL.")
		return err
	}
	return nil
}

func (s *Spec) SetDefaults() {
}

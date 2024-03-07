package client

import (
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type ErrTooManyFields struct{}

func (e *ErrTooManyFields) Error() string {
	return "too many fields has been set up."
}

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

func (s *Spec) Validate() error {
	endpointErr := validateEndpointURL(s.EndpointUrl)
	if endpointErr != nil {
		log.Error().Err(endpointErr).Msg("validation endpoint URL failed.")
		return endpointErr
	}
	appCredentialsErr := s.ValidateAppCredentials()
	if appCredentialsErr != nil {
		log.Error().Err(appCredentialsErr).Msg("too many fields has been set to use app credentials.")
		return appCredentialsErr
	}
	return nil
}

func (s *Spec) ValidateAppCredentials() error {
	// IMPORTANT NOTE: when using App Credentials, it is necessary
	// that all other fields except the endpoint URL be left blank!
	// Check that no other fields are set
	if s.AppCredentialID != nil && s.AppCredentialSecret != nil {
		if s.UserID != nil || s.Username != nil || s.Password != nil || s.Region != nil || s.ProjectID != nil || s.ProjectName != nil || s.DomainID != nil || s.DomainName != nil || s.AccessToken != nil {
			err := &ErrTooManyFields{}
			log.Error().Err(err).Msg("too many fields has been set to use app credentials.")
			return err
		}
	}
	return nil
}

func (s *Spec) SetDefaults() {
}

func validateEndpointURL(endpoint *string) error {
	// Check that the endpoint URL is a valid URL
	_, err := url.ParseRequestURI(*endpoint)
	if err != nil {
		log.Error().Err(err).Msg("invalid endpoint URL.")
		return err
	}
	// Check that the endpoint URL is reachable
	_, err = http.Get(*endpoint)
	if err != nil {
		log.Error().Err(err).Msg("unreachable endpoint URL.")
		return err
	}
	return nil
}

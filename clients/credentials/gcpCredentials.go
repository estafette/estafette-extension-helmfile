package credentials

// GCPInfraCredentials represents the credentials of type gcp-infra as defined in the server config and passed to this trusted image
type GCPInfraCredentials struct {
	Name                 string                                 `json:"name,omitempty"`
	Type                 string                                 `json:"type,omitempty"`
	AdditionalProperties GCPInfraCredentialAdditionalProperties `json:"additionalProperties,omitempty"`
}

// GCPInfraCredentialAdditionalProperties contains the non standard fields for this type of credentials
type GCPInfraCredentialAdditionalProperties struct {
	ServiceAccountKeyfile string `json:"serviceAccountKeyfile,omitempty"`
}

package options

import "encoding/json"

type PolicyType string

const (
	PolicyTypeAllow PolicyType = "allow"
	PolicyTypeDeny  PolicyType = "deny"
)

type HealthCheckType string

const (
	HealthCheckTypeHTTP  HealthCheckType = "http"
	HealthCheckTypeHTTPS HealthCheckType = "https"
)

type AntibotChallengeType string

const (
	AntibotChallenge      AntibotChallengeType = "ANTIBOT_CHALLENGE"
	AntibotChallengeJS    AntibotChallengeType = "ANTIBOT_CHALLENGE_JS"
	AntibotChallengeJSBox AntibotChallengeType = "ANTIBOT_CHALLENGE_JS_BOX"
)

type LocationOptions struct {
	AllowedHTTPMethods    Field[ListOption]                  `json:"allowed_http_methods,omitzero"`
	ASNACL                Field[ASNACLOption]                `json:"asn_acl,omitzero"`
	RefererACL            Field[AccessPolicyOption]          `json:"referer_acl,omitzero"`
	UserAgentACL          Field[AccessPolicyOption]          `json:"user_agent_acl,omitzero"`
	AntibotFloodChallenge Field[AntibotFloodChallengeOption] `json:"antibot_flood_challenge,omitzero"`
}

type ResourceOptions struct {
	LocationOptions
	UseDefaultLEChain                       Field[BooleanOption]               `json:"use_default_le_chain,omitzero"`
	UseRSALECert                            Field[BooleanOption]               `json:"use_rsa_le_cert,omitzero"`
	Blacklist                               Field[ListOption]                  `json:"blacklist,omitzero"`
	Whitelist                               Field[ListOption]                  `json:"whitelist,omitzero"`
	GeoACL                                  Field[GeoACLOption]                `json:"geo_acl,omitzero"`
	IPHash                                  Field[BooleanOption]               `json:"ip_hash,omitzero"`
	RedirectHTTPToHTTPS                     Field[BooleanOption]               `json:"redirect_http_to_https,omitzero"`
	ResponseHeaders                         Field[ResponseHeadersOption]       `json:"response_headers,omitzero"`
	TLSVersions                             Field[ListOption]                  `json:"tls_versions,omitzero"`
	CountryACL                              Field[AccessPolicyOption]          `json:"country_acl,omitzero"`
	WWWRedirects                            Field[BooleanOption]               `json:"www_redirects,omitzero"`
	SourceProtocolHTTPS                     Field[BooleanOption]               `json:"source_protocol_https,omitzero"`
	UpstreamHealthChecker                   Field[UpstreamHealthCheckerOption] `json:"upstream_healthchecker,omitzero"`
	WAF                                     Field[BooleanOption]               `json:"waf,omitzero"`
	MLAnalysis                              Field[BooleanOption]               `json:"ml_analysis,omitzero"`
	WAFMonitoringMode                       Field[BooleanOption]               `json:"waf_monitoring_mode,omitzero"`
	RedirectHTTPSToHTTP                     Field[BooleanOption]               `json:"redirect_https_to_http,omitzero"`
	AntibotChallengeURLExceptions           Field[ListOption]                  `json:"antibot_challenge_url_exceptions,omitzero"`
	AntibotFloodChallengeUserAgentWhitelist Field[ListOption]                  `json:"antibot_flood_challenge_user_agent_whitelist,omitzero"`
}

type BooleanOption struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type ListOption struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type AccessPolicyOption struct {
	Enabled        bool       `json:"enabled"`
	PolicyType     PolicyType `json:"policy_type"`
	ExceptedValues []string   `json:"excepted_values"`
}

type ASNACLOption struct {
	Enabled        bool       `json:"enabled"`
	PolicyType     PolicyType `json:"policy_type"`
	ExceptedValues []int64    `json:"excepted_values"`
}

type GeoACLOption struct {
	Enabled        bool                `json:"enabled"`
	PolicyType     PolicyType          `json:"policy_type"`
	ExceptedValues map[string][]string `json:"excepted_values"`
}

type ResponseHeadersOption struct {
	Enabled bool             `json:"enabled"`
	Value   []ResponseHeader `json:"value"`
}

type ResponseHeader struct {
	Name   string   `json:"name"`
	Value  []string `json:"value"`
	Always bool     `json:"always"`
}

type UpstreamHealthCheckerOption struct {
	Enabled   bool            `json:"enabled"`
	Host      *string         `json:"host"`
	Interval  int             `json:"interval,omitempty"`
	Rise      int             `json:"rise,omitempty"`
	Fall      int             `json:"fall,omitempty"`
	Timeout   int             `json:"timeout,omitempty"`
	Type      HealthCheckType `json:"type,omitempty"`
	HealthURL string          `json:"health_url"`
}

type AntibotFloodChallengeOption struct {
	Enabled     bool                  `json:"enabled"`
	Value       *AntibotChallengeType `json:"value"`
	TimeLimit   int                   `json:"antibot_challenge_timelimit,omitempty"`
	GetOnlyOn   bool                  `json:"antibot_challenge_get_only_on"`
	UTMReferrer bool                  `json:"antibot_challenge_utm_referrer"`
}

func (o ListOption) MarshalJSON() ([]byte, error) {
	type plain ListOption
	if o.Value == nil {
		o.Value = []string{}
	}

	return json.Marshal(plain(o))
}

func (o AccessPolicyOption) MarshalJSON() ([]byte, error) {
	type plain AccessPolicyOption
	if o.ExceptedValues == nil {
		o.ExceptedValues = []string{}
	}

	return json.Marshal(plain(o))
}

func (o ASNACLOption) MarshalJSON() ([]byte, error) {
	type plain ASNACLOption
	if o.ExceptedValues == nil {
		o.ExceptedValues = []int64{}
	}

	return json.Marshal(plain(o))
}

func (o GeoACLOption) MarshalJSON() ([]byte, error) {
	type plain GeoACLOption
	values := make(map[string][]string, len(o.ExceptedValues))
	for country, regions := range o.ExceptedValues {
		if regions == nil {
			regions = []string{}
		}
		values[country] = regions
	}
	o.ExceptedValues = values

	return json.Marshal(plain(o))
}

func (o ResponseHeadersOption) MarshalJSON() ([]byte, error) {
	type plain ResponseHeadersOption
	if o.Value == nil {
		o.Value = []ResponseHeader{}
	}

	return json.Marshal(plain(o))
}

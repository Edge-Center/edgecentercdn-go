package options

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const resourceOptionsResponseJSON = `{
	"use_default_le_chain": {"enabled": true, "value": true},
	"use_rsa_le_cert": null,
	"allowed_http_methods": {"enabled": true, "value": ["GET", "POST", "HEAD"]},
	"asn_acl": {"enabled": true, "policy_type": "allow", "excepted_values": [333, 4294967295]},
	"blacklist": {"enabled": true, "value": ["150.150.150.150/32", "2001:db8::/32"]},
	"whitelist": null,
	"geo_acl": {"enabled": false, "policy_type": "deny", "excepted_values": {"RU": ["MOW", "SPE"], "KZ": []}},
	"ip_hash": {"enabled": true, "value": false},
	"redirect_http_to_https": {"enabled": true, "value": true},
	"referer_acl": {"enabled": true, "policy_type": "deny", "excepted_values": ["", "*.example.com", "example.*"]},
	"response_headers": {
		"enabled": true,
		"value": [
			{"name": "X-Key", "value": ["878"], "always": false},
			{"name": "Cache-Control", "value": ["no-cache", "no-store"], "always": true}
		]
	},
	"tls_versions": {"enabled": true, "value": ["TLSv1.3", "TLSv1.2"]},
	"user_agent_acl": {"enabled": true, "policy_type": "allow", "excepted_values": ["~Mozilla", "^Chromoi.*"]},
	"country_acl": null,
	"www_redirects": {"enabled": false, "value": true},
	"source_protocol_https": {"enabled": true, "value": true},
	"upstream_healthchecker": {
		"enabled": true,
		"host": null,
		"interval": 15000,
		"rise": 3,
		"fall": 2,
		"timeout": 10000,
		"type": "https",
		"health_url": "/health"
	},
	"waf": {"enabled": true, "value": true},
	"ml_analysis": {"enabled": false, "value": false},
	"waf_monitoring_mode": null,
	"redirect_https_to_http": null,
	"antibot_flood_challenge": {
		"enabled": true,
		"value": "ANTIBOT_CHALLENGE",
		"antibot_challenge_timelimit": 5,
		"antibot_challenge_get_only_on": false,
		"antibot_challenge_utm_referrer": true
	},
	"antibot_challenge_url_exceptions": {"enabled": true, "value": ["/static/", "~^/api/"]},
	"antibot_flood_challenge_user_agent_whitelist": {"enabled": true, "value": ["~Insomni(\\w)?", "Mozilla"]}
}`

func marshalJSON(t *testing.T, v interface{}) string {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

func ptr[T any](v T) *T {
	return &v
}

func TestResourceOptions_RoundTripOfResponseBody(t *testing.T) {
	var opts ResourceOptions

	require.NoError(t, json.Unmarshal([]byte(resourceOptionsResponseJSON), &opts))

	assert.JSONEq(t, resourceOptionsResponseJSON, marshalJSON(t, opts))
}

func TestResourceOptions_DecodesResponseBody(t *testing.T) {
	var opts ResourceOptions

	require.NoError(t, json.Unmarshal([]byte(resourceOptionsResponseJSON), &opts))

	assert.Equal(t, Value(BooleanOption{Enabled: true, Value: true}), opts.WAF)
	assert.Equal(t, Value(ListOption{Enabled: true, Value: []string{"150.150.150.150/32", "2001:db8::/32"}}), opts.Blacklist)
	assert.Equal(t, Value(ASNACLOption{
		Enabled:        true,
		PolicyType:     PolicyTypeAllow,
		ExceptedValues: []int64{333, 4294967295},
	}), opts.ASNACL)
	assert.Equal(t, Value(GeoACLOption{
		Enabled:        false,
		PolicyType:     PolicyTypeDeny,
		ExceptedValues: map[string][]string{"RU": {"MOW", "SPE"}, "KZ": {}},
	}), opts.GeoACL)
	assert.Equal(t, Value(ResponseHeadersOption{
		Enabled: true,
		Value: []ResponseHeader{
			{Name: "X-Key", Value: []string{"878"}},
			{Name: "Cache-Control", Value: []string{"no-cache", "no-store"}, Always: true},
		},
	}), opts.ResponseHeaders)
	assert.Equal(t, Value(UpstreamHealthCheckerOption{
		Enabled:   true,
		Interval:  15000,
		Rise:      3,
		Fall:      2,
		Timeout:   10000,
		Type:      HealthCheckTypeHTTPS,
		HealthURL: "/health",
	}), opts.UpstreamHealthChecker)
	assert.Equal(t, Value(AntibotFloodChallengeOption{
		Enabled:     true,
		Value:       ptr(AntibotChallenge),
		TimeLimit:   5,
		UTMReferrer: true,
	}), opts.AntibotFloodChallenge)
}

func TestResourceOptions_MarshalSendsOnlyTouchedKeys(t *testing.T) {
	tests := []struct {
		name string
		opts ResourceOptions
		want string
	}{
		{
			name: "location level options on a resource",
			opts: ResourceOptions{
				LocationOptions: LocationOptions{
					ASNACL:       Value(ASNACLOption{Enabled: true, PolicyType: PolicyTypeDeny, ExceptedValues: []int64{64512}}),
					UserAgentACL: Null[AccessPolicyOption](),
				},
				SourceProtocolHTTPS: Value(BooleanOption{Enabled: true, Value: true}),
			},
			want: `{
				"asn_acl": {"enabled": true, "policy_type": "deny", "excepted_values": [64512]},
				"user_agent_acl": null,
				"source_protocol_https": {"enabled": true, "value": true}
			}`,
		},
		{
			name: "upstream healthchecker sends every sub-field",
			opts: ResourceOptions{UpstreamHealthChecker: Value(UpstreamHealthCheckerOption{
				Enabled:   true,
				Host:      ptr("health.example.com"),
				Interval:  5000,
				Rise:      2,
				Fall:      3,
				Timeout:   1000,
				Type:      HealthCheckTypeHTTP,
				HealthURL: "/status",
			})},
			want: `{"upstream_healthchecker": {
				"enabled": true,
				"host": "health.example.com",
				"interval": 5000,
				"rise": 2,
				"fall": 3,
				"timeout": 1000,
				"type": "http",
				"health_url": "/status"
			}}`,
		},
		{
			name: "upstream healthchecker leaves zero timings and empty type to server defaults",
			opts: ResourceOptions{UpstreamHealthChecker: Value(UpstreamHealthCheckerOption{
				Enabled:   true,
				HealthURL: "/health",
			})},
			want: `{"upstream_healthchecker": {
				"enabled": true,
				"host": null,
				"health_url": "/health"
			}}`,
		},
		{
			name: "response headers send always explicitly",
			opts: ResourceOptions{ResponseHeaders: Value(ResponseHeadersOption{
				Enabled: true,
				Value:   []ResponseHeader{{Name: "X-Frame-Options", Value: []string{"DENY"}}},
			})},
			want: `{"response_headers": {
				"enabled": true,
				"value": [{"name": "X-Frame-Options", "value": ["DENY"], "always": false}]
			}}`,
		},
		{
			name: "antibot dependent options",
			opts: ResourceOptions{
				LocationOptions: LocationOptions{
					AntibotFloodChallenge: Value(AntibotFloodChallengeOption{
						Enabled:   true,
						Value:     ptr(AntibotChallengeJSBox),
						TimeLimit: 30,
						GetOnlyOn: true,
					}),
				},
				AntibotChallengeURLExceptions:           Value(ListOption{Enabled: true, Value: []string{"/static/"}}),
				AntibotFloodChallengeUserAgentWhitelist: Null[ListOption](),
			},
			want: `{
				"antibot_flood_challenge": {
					"enabled": true,
					"value": "ANTIBOT_CHALLENGE_JS_BOX",
					"antibot_challenge_timelimit": 30,
					"antibot_challenge_get_only_on": true,
					"antibot_challenge_utm_referrer": false
				},
				"antibot_challenge_url_exceptions": {"enabled": true, "value": ["/static/"]},
				"antibot_flood_challenge_user_agent_whitelist": null
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.JSONEq(t, tt.want, marshalJSON(t, tt.opts))
		})
	}
}

func TestLocationOptions_MarshalSendsOnlyTouchedKeys(t *testing.T) {
	tests := []struct {
		name string
		opts LocationOptions
		want string
	}{
		{
			name: "antibot challenge without time limit leaves it to the server default",
			opts: LocationOptions{AntibotFloodChallenge: Value(AntibotFloodChallengeOption{
				Enabled: true,
				Value:   ptr(AntibotChallengeJS),
			})},
			want: `{"antibot_flood_challenge": {
				"enabled": true,
				"value": "ANTIBOT_CHALLENGE_JS",
				"antibot_challenge_get_only_on": false,
				"antibot_challenge_utm_referrer": false
			}}`,
		},
		{
			name: "disabled antibot challenge sends value null",
			opts: LocationOptions{AntibotFloodChallenge: Value(AntibotFloodChallengeOption{TimeLimit: 10})},
			want: `{"antibot_flood_challenge": {
				"enabled": false,
				"value": null,
				"antibot_challenge_timelimit": 10,
				"antibot_challenge_get_only_on": false,
				"antibot_challenge_utm_referrer": false
			}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.JSONEq(t, tt.want, marshalJSON(t, tt.opts))
		})
	}
}

func TestOptions_NilCollectionsAreSentAsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		option interface{}
		want   string
	}{
		{
			name:   "list option",
			option: ListOption{},
			want:   `{"enabled": false, "value": []}`,
		},
		{
			name:   "access policy option",
			option: AccessPolicyOption{PolicyType: PolicyTypeAllow},
			want:   `{"enabled": false, "policy_type": "allow", "excepted_values": []}`,
		},
		{
			name:   "asn acl option",
			option: ASNACLOption{PolicyType: PolicyTypeDeny},
			want:   `{"enabled": false, "policy_type": "deny", "excepted_values": []}`,
		},
		{
			name:   "geo acl option without countries",
			option: GeoACLOption{PolicyType: PolicyTypeDeny},
			want:   `{"enabled": false, "policy_type": "deny", "excepted_values": {}}`,
		},
		{
			name: "geo acl country without regions",
			option: GeoACLOption{
				Enabled:        true,
				PolicyType:     PolicyTypeAllow,
				ExceptedValues: map[string][]string{"RU": nil, "KZ": {"ALA"}},
			},
			want: `{"enabled": true, "policy_type": "allow", "excepted_values": {"RU": [], "KZ": ["ALA"]}}`,
		},
		{
			name:   "response headers option",
			option: ResponseHeadersOption{},
			want:   `{"enabled": false, "value": []}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.JSONEq(t, tt.want, marshalJSON(t, tt.option))
		})
	}
}

func TestGeoACLOption_MarshalDoesNotMutateCaller(t *testing.T) {
	geo := GeoACLOption{ExceptedValues: map[string][]string{"RU": nil}}

	marshalJSON(t, geo)

	assert.Nil(t, geo.ExceptedValues["RU"])
}

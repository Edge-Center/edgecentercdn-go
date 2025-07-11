package edgecenter

type LocationOptions struct {
	AllowedHTTPMethods          *AllowedHTTPMethods          `json:"allowed_http_methods"`
	BrotliCompression           *BrotliCompression           `json:"brotli_compression"`
	BrowserCacheSettings        *BrowserCacheSettings        `json:"browser_cache_settings"`
	Cors                        *Cors                        `json:"cors"`
	DisableProxyForceRanges     *DisableProxyForceRanges     `json:"disable_proxy_force_ranges"`
	EdgeCacheSettings           *EdgeCacheSettings           `json:"edge_cache_settings"`
	FetchCompressed             *FetchCompressed             `json:"fetch_compressed"`
	FollowOriginRedirect        *FollowOriginRedirect        `json:"follow_origin_redirect"`
	ForceReturn                 *ForceReturn                 `json:"force_return"`
	ForwardHostHeader           *ForwardHostHeader           `json:"forward_host_header"`
	GeoAcl                      *GeoAccessPolicy             `json:"geo_acl"`
	GzipCompression             *GzipCompression             `json:"gzip_compression"`
	HostHeader                  *HostHeader                  `json:"hostHeader"`
	IgnoreCookie                *IgnoreCookie                `json:"ignore_cookie"`
	IgnoreQueryString           *IgnoreQueryString           `json:"ignore_query_string"`
	ImageStack                  *ImageStack                  `json:"image_stack"`
	IPAddressACL                *IPAddressACL                `json:"ip_address_acl"`
	LimitBandwidth              *LimitBandwidth              `json:"limit_bandwidth"`
	ProxyCacheMethodsSet        *ProxyCacheMethodsSet        `json:"proxy_cache_methods_set"`
	QueryParamsBlacklist        *QueryParamsBlacklist        `json:"query_params_blacklist"`
	QueryParamsWhitelist        *QueryParamsWhitelist        `json:"query_params_whitelist"`
	RedirectHttpToHttps         *RedirectHttpToHttps         `json:"redirect_http_to_https"`
	RedirectHttpsToHttp         *RedirectHttpsToHttp         `json:"redirect_https_to_http"`
	RefererACL                  *RefererACL                  `json:"referer_acl"`
	ResponseHeadersHidingPolicy *ResponseHeadersHidingPolicy `json:"response_headers_hiding_policy"`
	Rewrite                     *Rewrite                     `json:"rewrite"`
	SecureKey                   *SecureKey                   `json:"secure_key"`
	Slice                       *Slice                       `json:"slice"`
	SNI                         *SNIOption                   `json:"sni"`
	Stale                       *Stale                       `json:"stale"`
	StaticRequestHeaders        *StaticRequestHeaders        `json:"static_request_headers"`
	StaticResponseHeaders       *StaticResponseHeaders       `json:"static_response_headers"`
	UserAgentACL                *UserAgentACL                `json:"user_agent_acl"`
	WebSockets                  *WebSockets                  `json:"websockets"`
}

type ResourceOptions struct {
	LocationOptions
	HTTP3Enabled      *HTTP3Enabled      `json:"http3_enabled"`
	TLSVersions       *TLSVersions       `json:"tls_versions"`
	UseDefaultLEChain *UseDefaultLEChain `json:"use_default_le_chain"`
}

type AllowedHTTPMethods struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type BrotliCompression struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type BrowserCacheSettings struct {
	Enabled bool   `json:"enabled"`
	Value   string `json:"value"`
}

type Cors struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
	Always  bool     `json:"always"`
}

type DisableProxyForceRanges struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type EdgeCacheSettings struct {
	Enabled      bool              `json:"enabled"`
	Value        string            `json:"value"`
	CustomValues map[string]string `json:"custom_values"`
	Default      string            `json:"default"`
}

type FetchCompressed struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type FollowOriginRedirect struct {
	Enabled bool  `json:"enabled"`
	Codes   []int `json:"codes"`
	UseHost bool  `json:"use_host"`
}

type ForceReturn struct {
	Enabled bool   `json:"enabled"`
	Code    int    `json:"code"`
	Body    string `json:"body"`
}

type ForwardHostHeader struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type GeoAccessPolicy struct {
	Enabled  bool                `json:"enabled"`
	Default  string              `json:"policy_type"`
	Excepted map[string][]string `json:"excepted_values"`
}

type GzipCompression struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type HostHeader struct {
	Enabled bool   `json:"enabled"`
	Value   string `json:"value"`
}

type HTTP3Enabled struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type IgnoreCookie struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type IgnoreQueryString struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type ImageStack struct {
	Enabled     bool `json:"enabled"`
	AvifEnabled bool `json:"avif_enabled"`
	WebpEnabled bool `json:"webp_enabled"`
	Quality     int  `json:"quality"`
	PngLossless bool `json:"png_lossless"`
}

type IPAddressACL struct {
	Enabled        bool     `json:"enabled"`
	PolicyType     string   `json:"policy_type"`
	ExceptedValues []string `json:"excepted_values"`
}

type LimitBandwidth struct {
	Enabled   bool   `json:"enabled"`
	LimitType string `json:"limit_type"`
	Speed     int    `json:"speed"`
	Buffer    int    `json:"buffer"`
}

type ProxyCacheMethodsSet struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type QueryParamsBlacklist struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type QueryParamsWhitelist struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type RedirectHttpToHttps struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type RedirectHttpsToHttp struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type RefererACL struct {
	Enabled        bool     `json:"enabled"`
	PolicyType     string   `json:"policy_type"`
	ExceptedValues []string `json:"excepted_values"`
}

type ResponseHeadersHidingPolicy struct {
	Enabled  bool     `json:"enabled"`
	Mode     string   `json:"mode"`
	Excepted []string `json:"excepted"`
}

type Rewrite struct {
	Enabled bool   `json:"enabled"`
	Body    string `json:"body"`
	Flag    string `json:"flag"`
}

type SecureKey struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
	Type    int    `json:"type"`
}

type Slice struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type SNIOption struct {
	Enabled        bool   `json:"enabled"`
	SNIType        string `json:"sni_type"`
	CustomHostname string `json:"custom_hostname"`
}

type Stale struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type StaticRequestHeaders struct {
	Enabled bool              `json:"enabled"`
	Value   map[string]string `json:"value"`
}

type StaticResponseHeadersItem struct {
	Name   string   `json:"name"`
	Value  []string `json:"value"`
	Always bool     `json:"always"`
}

type StaticResponseHeaders struct {
	Enabled bool                        `json:"enabled"`
	Value   []StaticResponseHeadersItem `json:"value"`
}

type TLSVersions struct {
	Enabled bool     `json:"enabled"`
	Value   []string `json:"value"`
}

type UseDefaultLEChain struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

type UserAgentACL struct {
	Enabled        bool     `json:"enabled"`
	PolicyType     string   `json:"policy_type"`
	ExceptedValues []string `json:"excepted_values"`
}

type WebSockets struct {
	Enabled bool `json:"enabled"`
	Value   bool `json:"value"`
}

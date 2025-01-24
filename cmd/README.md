# EdgeCenter CDN API client
   
## Build the binary:
```bash
go build -o edge-cli
```

## Usage

#### Retrieve whoami
```bash
./edge-cli whoami
```

#### Retrieve metrics for sent bytes and 5xx responses
```bash
./edge-cli metrics get --metric=sent_bytes --metric=responses_5xx --from=2025-01-01T00:00:00Z --to=2025-01-24T23:59:00Z --granularity=1d
```

#### Retrieve aggregate metrics for sent bytes and 5xx responses
```bash
./edge-cli metrics get --metric=sent_bytes --metric=responses_5xx --from=2025-01-01T00:00:00Z --to=2025-01-24T23:59:00Z
```

#### Group metrics by resource
```bash
./edge-cli metrics get --groupby=resource --metric=sent_bytes --metric=origin_response_time --metric=responses_5xx --from=2025-01-01T00:00:00Z --to=2025-01-24T23:59:00Z
```

#### Group metrics by country
```bash
./edge-cli metrics get --groupby=country --metric=sent_bytes --metric=origin_response_time --metric=responses_5xx --from=2025-01-01T00:00:00Z --to=2025-01-24T23:59:00Z
```

#### Group metrics by resource output as table
```bash
./edge-cli metrics get --output=table --groupby=resource --groupby=vhost --metric=request_time --metric=origin_response_time --metric=cdn_bytes  --from=2025-01-01T00:00:00Z --to=2025-01-24T23:59:00Z
```

#### Filter by virtual hosts
```bash
./edge-cli metrics get -o=table --groupby=vhost --metric=sent_bytes  --from=2025-01-01T00:00:00Z --to=2025-01-24T23:59:00Z
```

#### Resources
```bash
./edge-cli resource list
```

#### Resources fields
```bash
./edge-cli resource list --field=id --field=cname --status=active --search=preprod
```

#### Resources IDS
```bash
./edge-cli resource list | jq '.[].id'
```

#### Create origin
```bash
./edge-cli origins create --name=example.com --source=$example.com
```

#### Create origin with auth
```bash
./edge-cli origins create \
    --name="example2.com" \
    --source="example2.com" \
    --auth-type="aws_signature_v2" \
    --bucket-name="test" \
    --access-key-id="123123234" \
    --secret-key="asdnaskjdnajdasdjansndjasndjknsdksdna"
```


module github.com/brienze1/crypto-robot-operation-hub

go 1.19

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go-v2 v1.16.12
	github.com/aws/aws-sdk-go-v2/config v1.17.3
	github.com/aws/aws-sdk-go-v2/credentials v1.12.16
	github.com/aws/aws-sdk-go-v2/service/sns v1.17.15
	github.com/cucumber/godog v0.12.5
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/lib/pq v1.10.6
	github.com/stretchr/testify v1.7.2
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.15 // indirect
	github.com/aws/smithy-go v1.13.0 // indirect
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-memdb v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gogo/protobuf v1.2.1 => github.com/gogo/protobuf v1.3.2
	github.com/miekg/dns v1.0.14 => github.com/miekg/dns v1.1.50
	github.com/prometheus/client_golang v0.9.3 => github.com/prometheus/client_golang v1.13.0
	github.com/prometheus/client_golang v1.4.0 => github.com/prometheus/client_golang v1.13.0
	golang.org/x/text v0.3.3 => golang.org/x/text v0.3.7
)

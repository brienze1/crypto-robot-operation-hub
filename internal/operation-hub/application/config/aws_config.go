package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"sync"
)

var (
	sessionInit  sync.Once
	snsInit      sync.Once
	dynamoDBInit sync.Once
)

var (
	awsConfig      *aws.Config
	snsClient      *sns.Client
	dynamoDbClient *dynamodb.Client
)

func getConfig() *aws.Config {
	if awsConfig == nil {
		sessionInit.Do(
			func() {
				if properties.Properties().Aws.Config.OverrideConfig {
					newAwsConfig, err := config.LoadDefaultConfig(context.TODO(),
						config.WithEndpointResolverWithOptions(NewEndpointResolver()),
						config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
							properties.Properties().Aws.Config.AccessKey,
							properties.Properties().Aws.Config.AccessSecret,
							properties.Properties().Aws.Config.Token)))
					if err != nil {
						panic("configuration error, " + err.Error())
					}
					awsConfig = &newAwsConfig
				} else {
					newAwsConfig, err := config.LoadDefaultConfig(context.TODO())
					if err != nil {
						panic("configuration error, " + err.Error())
					}
					awsConfig = &newAwsConfig
				}
			})
	}

	return awsConfig
}

func NewEndpointResolver() aws.EndpointResolverWithOptionsFunc {
	return func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               properties.Properties().Aws.Config.URL,
			SigningRegion:     properties.Properties().Aws.Config.Region,
			HostnameImmutable: true,
		}, nil
	}
}

func SNSClient() *sns.Client {
	if awsConfig == nil {
		snsInit.Do(func() {
			cfg := getConfig()

			snsClient = sns.NewFromConfig(*cfg)
		})
	}

	return snsClient
}

func DynamoDBClient() *dynamodb.Client {
	if awsConfig == nil {
		dynamoDBInit.Do(func() {
			cfg := getConfig()

			dynamoDbClient = dynamodb.NewFromConfig(*cfg)
		})
	}

	return dynamoDbClient
}

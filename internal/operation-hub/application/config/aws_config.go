package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"sync"
)

var (
	sessionInit sync.Once
	snsInit     sync.Once
)

var (
	awsSession *session.Session
	snsClient  *sns.SNS
)

func newSession() *session.Session {
	if awsSession == nil {
		sessionInit.Do(
			func() {
				if properties.Properties().Aws.OverrideConfig {
					awsSession = session.Must(session.NewSession(&aws.Config{
						Credentials: credentials.NewStaticCredentials(
							properties.Properties().Aws.AccessKey,
							properties.Properties().Aws.AccessSecret,
							properties.Properties().Aws.Token,
						),
						Endpoint: properties.Properties().Aws.URL,
						Region:   properties.Properties().Aws.Region,
					}))
				} else {
					awsSession = session.Must(session.NewSessionWithOptions(session.Options{
						SharedConfigState: session.SharedConfigEnable,
					}))
				}
			})
	}

	return awsSession
}

func SNSClient() *sns.SNS {
	if awsSession == nil {
		snsInit.Do(
			func() {
				newSession := newSession()

				snsClient = sns.New(newSession)
			})
	}

	return snsClient
}

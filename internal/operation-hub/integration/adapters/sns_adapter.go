package adapters

import "github.com/aws/aws-sdk-go/service/sns"

type SNSAdapter interface {
	Publish(input *sns.PublishInput) (*sns.PublishOutput, error)
}

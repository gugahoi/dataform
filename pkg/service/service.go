package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// Service is an abstraction from rds.Service so that consumers dont need the AWS SDK
type Service interface {
	rdsiface.RDSAPI
}

// New returns an aws session with a region when a non nil region is provided, defaults otherwise
func New(region string) Service {
	return rds.New(NewSession(region))
}

// NewSession returns an aws session
func NewSession(region string) *session.Session {
	if region != "" {
		return session.New(&aws.Config{
			Region: aws.String(region),
		})
	}
	return session.New(aws.NewConfig())
}

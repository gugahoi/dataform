package main

import (
	"flag"

	"github.com/MYOB-Technology/dataform/pkg/dfm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	var region string

	region = *flag.String("region", "", "AWS Region")
	flag.Parse()

	svc := rds.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))

	manager := dfm.NewManager(svc)

	manager.Create("some-db")
}

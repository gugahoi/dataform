package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "dfm",
	Short: "Give me ALL your RDS",
}

var (
	awsRegion string = ""
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&awsRegion, "region", "", "", "AWS Region")
}

// getAwsSession check for provided region and return an aws session
func getAwsSession() *session.Session {
	if awsRegion != "" {
		return session.New(&aws.Config{
			Region: aws.String(awsRegion),
		})
	}
	return session.New(aws.NewConfig())
}

// getAwsError return just the error message
func getAwsError(err error) string {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Message()
	}
	return err.Error()
}

// Execute ...
func Execute() {
	RootCmd.Execute()
}

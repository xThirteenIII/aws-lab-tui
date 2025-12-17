package model

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// commands.go collects all asyncronous functions.

// fetchS3FilesCmd returns a S3FilesMsg to be catched by Update() before entering S3 state.
func fetchS3FilesCmd(s string) tea.Cmd {

	return func() tea.Msg {

		// Load aws configurations from .aws files
		// User needs to have .aws/credentials and .aws/config files
		ctx := context.Background()
		awsConf, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			// Send an ErrorMsg if cannot read .aws files
			return ErrorMsg{Err: fmt.Errorf("cannot load aws configs: %v", err)}
		}

		// TODO: Load variables into a map when InitModel
		s3Bucket, _ := MustEnv("S3_TEST_OTA_BUCKET")
		s3Client := s3.NewFromConfig(awsConf)
		listObjectsV2Input := &s3.ListObjectsV2Input{
			Bucket:    aws.String(s3Bucket),
			Prefix:    aws.String(s),
			Delimiter: aws.String("/"),
		}

		listObjectsV2Output, err := s3Client.ListObjectsV2(ctx, listObjectsV2Input)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		files := []list.Item{}

		if listObjectsV2Output == nil {
			return ErrorMsg{Err: fmt.Errorf("no files retrieved")}
		}

		// Add folders
		for _, folder := range listObjectsV2Output.CommonPrefixes {
			newFolder := item{title: path.Base(*folder.Prefix)}
			files = append(files, newFolder)
		}

		// Add files
		for _, file := range listObjectsV2Output.Contents {
			newFile := item{title: path.Base(*file.Key)}
			files = append(files, newFile)
		}

		// Return files but not the first one, which is the current folder.
		return S3FilesMsg{Files: files[1:]}
	}
}

func sendIoTJob(m *Model) tea.Cmd {

	// Valide data. Just checks if anything is an empty string for now.
	// TODO: validate that name is not used it's far too expensive, since it would mean to
	// check every item in the Jobs AWS page.
	if !m.iotTool.IsValid() {
		return func() tea.Msg {
			return ErrorMsg{Err: fmt.Errorf("missing required job data")}
		}
	}

	return func() tea.Msg {

		// Create context
		ctx := context.Background()
		// Load aws configurations from .aws
		awsConf, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			// Send an ErrorMsg if cannot read .aws files
			return ErrorMsg{Err: fmt.Errorf("cannot load aws configs: %v", err)}
		}

		// Load environment variables from .env file
		awsAccountID, _ := MustEnv("AWS_ACCOUNT_ID")
		presignedRole, _ := MustEnv("IOT_JOB_PRESIGNED_ROLE")
		s3Bucket, _ := MustEnv("S3_TEST_OTA_BUCKET")

		// Configurate presignedUrl
		// This is used for more security.
		roleArn := "arn:aws:iam::" + awsAccountID + ":role/" + presignedRole
		presignedUrl := &types.PresignedUrlConfig{
			ExpiresInSec: &m.iotTool.JobExpiresInSeconds,
			RoleArn:      &roleArn,
		}

		// jobInput is needed to send the job
		jobInput := &iot.CreateJobInput{}
		jobInput.TargetSelection = "SNAPSHOT"
		jobInput.PresignedUrlConfig = presignedUrl
		jobInput.Targets = []string{"arn:aws:iot:eu-west-2:" + awsAccountID + ":thing/" + m.iotTool.ThingName}
		jobInput.JobId = &m.iotTool.JobName
		jobInput.DocumentSource = aws.String("s3://" + s3Bucket + "/" + m.iotTool.DocumentKey)

		// Create iot Client
		iotClient := iot.NewFromConfig(awsConf)
		_, err = iotClient.CreateJob(ctx, jobInput)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return IoTJobMsg{
			&iot.DescribeJobExecutionInput{
				JobId:           aws.String(m.iotTool.JobName),
				ThingName:       aws.String(m.iotTool.ThingName),
				ExecutionNumber: aws.Int64(1),
			},
		}
	}
}

func monitorIoTJobStatus() tea.Cmd {

	return nil
}

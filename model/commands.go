package model

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

		return S3FilesMsg{Files: files}
	}
}

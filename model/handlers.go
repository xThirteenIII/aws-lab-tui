package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

// GetS3Files returns a tea.Cmd and a string for the error
func (m *model) loadS3Files() tea.Cmd {

	/*
		// awsConf reads in the .aws user folder
		awsConf, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			m.errorChannel <- err
			return nil
		}

		//s3Client := s3.NewFromConfig(awsConf)
	*/
	return nil
}

func (m *model) waitForS3Files() tea.Cmd {
	return func() tea.Msg {
		select {
		case res := <-m.s3FilesLoadedChannel:
			return res
		case err := <-m.errorChannel:
			m.err = err.Error()
			return err
		}
	}
}

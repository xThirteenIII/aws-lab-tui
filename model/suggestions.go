package model

type suggestions struct {
	jobSuggestions []string
}

func (s *suggestions) addJobSuggestion(sug string) {

	if s.jobSuggestions == nil {
		s.jobSuggestions = make([]string, 0)
	}
	s.jobSuggestions = append(s.jobSuggestions, sug)
}

package model

import (
	"aws-iot-tui/cache"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
)

type suggestions struct {
	jobSuggestions []string
	fileName       string
}

func (s *suggestions) addJobSuggestion(sug string) {
	s.jobSuggestions = nil //reset list
	s.loadFromCache()

	// Prepend
	s.jobSuggestions = append([]string{sug}, s.jobSuggestions...)
	if len(s.jobSuggestions) > 100 {
		s.jobSuggestions = s.jobSuggestions[:100]
	}

	s.saveToCache()
}

func (s *suggestions) loadFromCache() {
	data, err := os.ReadFile(s.fileName)
	if err != nil || len(data) == 0 {
		return // return if file does not exist or it's empty
	}

	// cacheRoot is of type Cache (see ../cache/Cache.go), which is a FlatBuffers object
	cacheRoot := cache.GetRootAsCache(data, 0)
	length := cacheRoot.JobsLength()

	// Does data have some bytes written in it?
	if len(data) > 0 {

		for i := range length {
			// If so, take the i-th string from FlatBuffers vector as []byte.
			jobBytes := cacheRoot.Jobs(i)

			// Append to jobSuggestions if there's something in there
			if jobBytes != nil {
				s.jobSuggestions = append(s.jobSuggestions, string(jobBytes))
			}
		}
	}
}

func (s *suggestions) saveToCache() {

	// Create an initial FlatBuffers builder of 1024 B
	builder := flatbuffers.NewBuilder(1024)

	var offsets []flatbuffers.UOffsetT
	for _, job := range s.jobSuggestions {
		offsets = append(offsets, builder.CreateString(job))
	}
	cache.CacheStartJobsVector(builder, len(offsets))
	for i := len(offsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(offsets[i])
	}

	jobOffset := builder.EndVector(len(offsets))
	cache.CacheStart(builder)
	cache.CacheAddJobs(builder, jobOffset)
	cacheOffset := cache.CacheEnd(builder)
	builder.Finish(cacheOffset)

	// Save on file, if file does not exists, WriteFile creates it.
	err := os.WriteFile(s.fileName, builder.FinishedBytes(), 0644)
	if err != nil {
		panic(err)
	}
}

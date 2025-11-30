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

	// Create a new slice of strings and add the new job suggestion at the head
	s.jobSuggestions = append([]string{sug}, s.jobSuggestions...)

	// Limit to 100 records
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

	// Does data have some bytes written in it?
	if len(data) > 0 {
		// how many jobs are already saved in the cache file?
		length := cacheRoot.JobsLength()

		// Cool, let's read each job
		for i := range length {

			// Read a job as a []byte
			jobBytes := cacheRoot.Jobs(i)

			// Convert []byte to string, then append it to jobSuggestions
			// NB: this DOES NOT add the new job suggestion, we're just creating
			// loading already existing suggestions and appending it to our field
			if jobBytes != nil {
				s.jobSuggestions = append(s.jobSuggestions, string(jobBytes))
			}
		}
	}
}

func (s *suggestions) saveToCache() {

	// Create an initial FlatBuffers builder of 1024 B
	builder := flatbuffers.NewBuilder(1024)

	// offsets are like pointers, they tell us where to search for a specific string
	// in the buffer
	var offsets []flatbuffers.UOffsetT

	// for every suggestion in the slice (i.e. job), append its offset
	for _, job := range s.jobSuggestions {
		offsets = append(offsets, builder.CreateString(job))
	}

	// Start jobs vector creation in the FlatBuffer
	cache.CacheStartJobsVector(builder, len(offsets))

	// FlatBuffers rule: offsets go prepending
	for i := len(offsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(offsets[i])
	}

	// Close the vector and save its offset
	jobOffset := builder.EndVector(len(offsets))

	// Create Cache object that contains the vector
	cache.CacheStart(builder)
	cache.CacheAddJobs(builder, jobOffset)
	cacheOffset := cache.CacheEnd(builder)

	// Complete the buffer
	builder.Finish(cacheOffset)

	// Save on file, if file does not exists, WriteFile creates it.
	err := os.WriteFile(s.fileName, builder.FinishedBytes(), 0644)
	if err != nil {
		panic(err)
	}
}

package refactor

import (
	"aws-iot-tui/cache"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
)

// suggestions holds job and mac suggestions, and the cache file name
type suggestions struct {
	jobSuggestions []string
	macSuggestions []string
	cacheFile      string
}

// addJobSuggestion prepend the new suggestion to the jobSuggestion slice
func (s *suggestions) addJobSuggestion(sug string) {
	s.loadFromCache()

	// Create a new slice of strings and add the new job suggestion at the head
	s.jobSuggestions = append([]string{sug}, s.jobSuggestions...)

	// Limit to 50 records
	if len(s.jobSuggestions) > 50 {
		s.jobSuggestions = s.jobSuggestions[:50]
	}

	s.saveToCache()
}

// addJobSuggestion prepend the new suggestion to the macSuggestions slice
func (s *suggestions) addMacSuggestion(sug string) {
	s.loadFromCache()

	// Create a new slice of strings and add the new mac suggestion at the head
	s.macSuggestions = append([]string{sug}, s.macSuggestions...)

	// Limit to 100 records
	if len(s.macSuggestions) > 50 {
		s.macSuggestions = s.macSuggestions[:50]
	}

	s.saveToCache()
}

// WARNING: for now, we read and re-write the whole cache.bin file.
//
//	We can do it since flatBuffer performance is not affected, only 100 records in our fb.
//	If the file is *K or 10*K +, then we can start reading and writing only necessary parts.
func (s *suggestions) loadFromCache() {

	// These two assignments make sure the cache.bin file does not get read and written everytime over
	// and over.
	s.jobSuggestions = nil
	s.macSuggestions = nil
	data, err := os.ReadFile(s.cacheFile)
	if err != nil || len(data) == 0 {
		return // return if file does not exist or it's empty
	}

	// cacheRoot is of type Cache (see ../cache/Cache.go), which is a FlatBuffers object
	cacheRoot := cache.GetRootAsCache(data, 0)

	// Does data have some bytes written in it?
	if len(data) > 0 {
		// how many jobs are already saved in the cache file?
		jobsLength := cacheRoot.JobsLength()

		// how many macs are already saved in the cache file?
		macsLength := cacheRoot.MacsLength()

		// Cool, let's read each job
		for i := range jobsLength {

			// Read a job as a []byte
			jobBytes := cacheRoot.Jobs(i)

			// Convert []byte to string, then append it to jobSuggestions
			// NB: this DOES NOT add the new job suggestion, we're just creating
			// loading already existing suggestions and appending it to our field
			if jobBytes != nil {
				s.jobSuggestions = append(s.jobSuggestions, string(jobBytes))
			}
		}

		// Now let's read each mac
		for i := range macsLength {

			// Read a job as a []byte
			macBytes := cacheRoot.Macs(i)

			// Convert []byte to string, then append it to jobSuggestions
			// NB: this DOES NOT add the new job suggestion, we're just creating
			// loading already existing suggestions and appending it to our field
			if macBytes != nil {
				s.macSuggestions = append(s.macSuggestions, string(macBytes))
			}
		}
	}
}

// saveToCache creates a flatbuffer vector for suggestions and writes it to the cacheFile, as bytes
func (s *suggestions) saveToCache() {

	// Create an initial FlatBuffers builder of 1024 B
	builder := flatbuffers.NewBuilder(1024)

	// offsets are like pointers, they tell us where to search for a specific string
	// in the buffer
	var jobOffsets []flatbuffers.UOffsetT

	// for every suggestion in the slice (i.e. job), append its offset
	for _, job := range s.jobSuggestions {
		jobOffsets = append(jobOffsets, builder.CreateString(job))
	}

	var macOffsets []flatbuffers.UOffsetT

	// for every suggestion in the slice (i.e. mac), append its offset
	for _, mac := range s.macSuggestions {
		macOffsets = append(macOffsets, builder.CreateString(mac))
	}

	// Start jobs vector creation in the FlatBuffer
	cache.CacheStartJobsVector(builder, len(jobOffsets))
	// FlatBuffers rule: offsets go prepending
	for i := len(jobOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(jobOffsets[i])
	}
	jobOffset := builder.EndVector(len(jobOffsets))

	// Start macs vector creation in the FlatBuffer
	cache.CacheStartMacsVector(builder, len(macOffsets))
	for i := len(macOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(macOffsets[i])
	}
	// Close the vectors and save theirs offsets
	macOffset := builder.EndVector(len(macOffsets))

	// Create Cache object that contains the vectors
	cache.CacheStart(builder)
	cache.CacheAddJobs(builder, jobOffset)
	cache.CacheAddMacs(builder, macOffset)
	cacheOffset := cache.CacheEnd(builder)

	// Complete the buffers
	builder.Finish(cacheOffset)

	// Save on file, if file does not exists, WriteFile creates it.
	err := os.WriteFile(s.cacheFile, builder.FinishedBytes(), 0644)
	if err != nil {
		panic(err)
	}
}

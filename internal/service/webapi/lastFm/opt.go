package lastFm

type opt struct {
	flowSize          int
	maxPerSource      int
	maxPerArtist      int
	maxSimilarArtists int
	maxTopPerArtist   int
}

type Set func(e *opt)

// Parameter to set the max number of similar songs per item source.
// Example: if the items source has two songs, a value of 3 will result in 6 songs in the response.
func MaxPerSource(num int) Set {
	return func(e *opt) { e.maxPerSource = num }
}

// Parameter to set the number of songs from an artist.
// Example: when set to 2, the response will contain 2 songs from the same artist.
func MaxPerArtist(num int) Set {
	return func(e *opt) { e.maxPerArtist = num }
}

// Parameter to set the number of items handled by a single thread that were retrieved by the search.
func FlowSize(num int) Set {
	return func(e *opt) { e.flowSize = num }
}

// Parameter to limit the number when requesting similar artists.
func MaxSimilarArtists(num int) Set {
	return func(e *opt) { e.maxSimilarArtists = num }
}

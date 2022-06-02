package lastFm

type options struct {
	quantityFlow            int
	maxAudioAmountPerSource int
	maxAudioAmountPerArtist int
	maxNumSimiliarArtists   int
	maxNumTopPerArtist      int
}

type ProcessingOptions func(e *options)

// Parameter to set the max number of similar songs per item source.
// Example: if the items source has two songs, a value of 3 will result in 6 songs in the response.
func SetMaxAudioAmountPerSource(num int) ProcessingOptions {
	return func(e *options) {
		e.maxAudioAmountPerSource = num
	}
}

// Parameter to set the number of songs from an artist.
// Example: when set to 2, the response will contain 2 songs from the same artist.
func SetMaxAudioAmountPerArtist(num int) ProcessingOptions {
	return func(e *options) {
		e.maxAudioAmountPerArtist = num
	}
}

// Parameter to set the number of items handled by a single thread that were retrieved by the search.
func SetQuantityFlow(num int) ProcessingOptions {
	return func(e *options) {
		e.quantityFlow = num
	}
}

// Parameter to limit the number when requesting similar artists.
func SetMaxNumSimiliarArtists(num int) ProcessingOptions {
	return func(e *options) {
		e.maxNumSimiliarArtists = num
	}
}

// Parameter to limit the number of top songs from an artist.
func SetMaxNumTopPerArtist(num int) ProcessingOptions {
	return func(e *options) {
		e.maxNumTopPerArtist = num
	}
}

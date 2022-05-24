package lastFm

type options struct {
	quantityFlow            int
	maxNumberOfSong         int
	maxAudioAmountPerSource int
	maxAudioAmountPerArtist int
	maxNumSimiliarArtists   int
	maxNumTopPerArtist      int
}

type processingOptions func(e *options)

func setMaxAudioAmountPerSource(num int) processingOptions {
	return func(e *options) {
		e.maxAudioAmountPerSource = num
	}
}

func setMaxAudioAmountPerArtist(num int) processingOptions {
	return func(e *options) {
		e.maxAudioAmountPerArtist = num
	}
}

func setQuantityFlow(num int) processingOptions {
	return func(e *options) {
		e.quantityFlow = num
	}
}

func setNumberOfSong(num int) processingOptions {
	return func(e *options) {
		e.maxNumberOfSong = num
	}
}

func maxNumSimiliarArtists(num int) processingOptions {
	return func(e *options) {
		e.maxNumSimiliarArtists = num
	}
}
func maxNumTopPerArtist(num int) processingOptions {
	return func(e *options) {
		e.maxNumTopPerArtist = num
	}
}

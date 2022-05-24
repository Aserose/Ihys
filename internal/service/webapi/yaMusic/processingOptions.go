package yaMusic

type options struct {
	quantityFlow         int
	numberOfSong         int
	audioAmountPerSource int
	audioAmountPerArtist int
}

type processingOptions func(e *options)

func setAudioAmountPerSource(num int) processingOptions {
	return func(e *options) {
		e.audioAmountPerSource = num
	}
}

func setAudioAmountPerArtist(num int) processingOptions {
	return func(e *options) {
		e.audioAmountPerArtist = num
	}
}

func setQuantityFlow(num int) processingOptions {
	return func(e *options) {
		e.quantityFlow = num
	}
}

func setNumberOfSong(num int) processingOptions {
	return func(e *options) {
		e.numberOfSong = num
	}
}

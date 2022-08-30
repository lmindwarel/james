package spotify

import (
	"sync"
	"unsafe"

	"github.com/xlab/portaudio-go/portaudio"
)

// PortAudio helpers
func PAError(err portaudio.Error) bool {
	return portaudio.ErrorCode(err) != portaudio.PaNoError

}

func PAErrorText(err portaudio.Error) string {
	return "PortAudio error: " + portaudio.GetErrorText(err)
}

func paCallback(wg *sync.WaitGroup, channels int, samples <-chan [][]float32) portaudio.StreamCallback {
	wg.Add(1)
	return func(_ unsafe.Pointer, output unsafe.Pointer, sampleCount uint,
		_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

		log.Debugf("paCallback")

		const (
			statusContinue = int32(portaudio.PaContinue)
			statusComplete = int32(portaudio.PaComplete)
		)

		frame, ok := <-samples
		if !ok {
			wg.Done()
			return statusComplete
		}
		if len(frame) > int(sampleCount) {
			frame = frame[:sampleCount]
		}

		var idx int
		out := (*(*[1 << 32]float32)(unsafe.Pointer(output)))[:int(sampleCount)*channels]
		for _, sample := range frame {
			if len(sample) > channels {
				sample = sample[:channels]
			}
			for i := range sample {
				out[idx] = sample[i]
				idx++
			}
		}

		return statusContinue
	}
}

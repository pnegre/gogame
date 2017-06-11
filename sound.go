package gogame

/*
#cgo pkg-config: sdl2
#cgo LDFLAGS: -lm
#include "SDL.h"

extern void setAmplitude(int id, int amp);
extern void setFreq(int id, int freq);
extern void audioCallback(void* userdata, Uint8* stream, int len);
extern SDL_AudioDeviceID newAudioDevice();

*/
import "C"
import "errors"
import "unsafe"
import "math"

var toneCache = make(map[int]*ToneGenerator)

//export soundGoCallback
func soundGoCallback(id int, ptr unsafe.Pointer, len int) {
	tg := toneCache[id]
	len /= 4
	slice := (*[1 << 30]float32)(unsafe.Pointer(ptr))[:len:len]
	tg.feedSamples(slice)
}


type ToneGenerator struct {
	dev C.SDL_AudioDeviceID
	amp float32
	freq float32
}

func NewToneGenerator() (*ToneGenerator, error) {
	dev := C.newAudioDevice()
	if dev == 0 {
		return nil, errors.New("Can't open tone generator")
	}
	sd := new(ToneGenerator)
	sd.dev = dev
	toneCache[int(sd.dev)] = sd
	return sd, nil
}

func (self *ToneGenerator) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *ToneGenerator) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *ToneGenerator) SetFreq(freq int) {
	self.freq = float32(freq)
	// C.setFreq(C.int(self.dev), C.int(freq))
}

func (self *ToneGenerator) SetAmplitude(amp int) {
	self.amp = float32(amp)
	// C.setAmplitude(C.int(self.dev), C.int(amp))
}

func (self *ToneGenerator) Close() {
	C.SDL_CloseAudioDevice(self.dev)
}

func (self *ToneGenerator) feedSamples(data []float32) {
	var v float32
	for i:=0; i<len(data); i++ {
		data[i] = self.amp * float32(math.Sin(float64(v * 2 * math.Pi / 44100)))
		v += self.freq
	}
	//
	// for(int i=0, j=0; i<len / 4; i++) {
	// 	data[i] = amp * sin(v * 2 * M_PI / FREQUENCY);
	// 	if (j > N) {
	// 		v -= N*freq;
	// 		j = 0;
	// 	} else {
	// 		v += freq;
	// 		j++;
	// 	}
	// }
}

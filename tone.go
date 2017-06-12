package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

extern SDL_AudioDeviceID newAudioDevice(int frequency);
extern void closeAudioDevice(SDL_AudioDeviceID id);

*/
import "C"
import "errors"
import "unsafe"
import "math"
import "math/rand"

const FREQUENCY = 44100

var toneCache = make(map[int]*ToneGenerator)

const (
	GENERATOR_TYPE_TONE = iota
	GENERATOR_TYPE_NOISE
)

//export soundGoCallback
func soundGoCallback(id int, ptr unsafe.Pointer, len int) {
	len /= 4
	// To create a Go slice backed by a C array
	// (without copying the original data), one needs to acquire this length
	// at runtime and use a type conversion to a pointer to a very big array and then
	// slice it to the length that you want (also remember to set
	// the cap if you're using Go 1.2 or later)
	slice := (*[1 << 30]float32)(ptr)[:len:len]
	toneCache[id].feedSamples(slice)
}

type ToneGenerator struct {
	genType int
	dev     C.SDL_AudioDeviceID
	amp     float32
	freq    float32
	v       float32
	period  int
	j       int
}

func NewToneGenerator(genType int) (*ToneGenerator, error) {
	dev := C.newAudioDevice(FREQUENCY)
	if dev == 0 {
		return nil, errors.New("Can't open tone generator")
	}
	sd := new(ToneGenerator)
	sd.dev = dev
	sd.genType = genType
	toneCache[int(sd.dev)] = sd
	return sd, nil
}

func (self *ToneGenerator) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *ToneGenerator) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *ToneGenerator) SetFreq(freq float32) {
	self.freq = freq
	self.period = int(2 * FREQUENCY / self.freq)
}

func (self *ToneGenerator) SetAmplitude(amp float32) {
	self.amp = amp
}

func (self *ToneGenerator) Close() {
	C.closeAudioDevice(self.dev)
	delete(toneCache, int(self.dev))
}

func (self *ToneGenerator) feedSamples(data []float32) {
	if self.genType == GENERATOR_TYPE_TONE {
		for i := 0; i < len(data); i++ {
			data[i] = self.amp * float32(math.Sin(float64(self.v*2*math.Pi/FREQUENCY)))
			if self.j > self.period {
				self.v -= float32(self.period) * self.freq
				self.j = 0
				} else {
					self.v += self.freq
					self.j++
				}
			}
	}  else if self.genType == GENERATOR_TYPE_NOISE {
		for i := 0; i < len(data); i++ {
			data[i] = self.amp * rand.Float32();
		}
	}
}

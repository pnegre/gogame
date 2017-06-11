package gogame

/*
#cgo pkg-config: sdl2
#cgo LDFLAGS: -lm
#include "SDL.h"

extern void setAmplitude(int id, int amp);
extern void setFreq(int id, int freq);
extern void audioCallback(void* userdata, Uint8* stream, int len);
extern SDL_AudioDeviceID newAudioDevice(int frequency);
extern void closeAudioDevice(SDL_AudioDeviceID id);

*/
import "C"
import "errors"
import "unsafe"
import "math"

const FREQUENCY = 44100

var toneCache = make(map[int]*ToneGenerator)

//export soundGoCallback
func soundGoCallback(id int, ptr unsafe.Pointer, len int) {
	tg := toneCache[id]
	len /= 4
	slice := (*[1 << 30]float32)(unsafe.Pointer(ptr))[:len:len]
	tg.feedSamples(slice)
}

type ToneGenerator struct {
	dev    C.SDL_AudioDeviceID
	amp    float32
	freq   float32
	v      float32
	period int
	j int
}

func NewToneGenerator() (*ToneGenerator, error) {
	dev := C.newAudioDevice(FREQUENCY)
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
	self.v = 0
	self.j = 0
	self.period = int(2 * FREQUENCY / self.freq)
}

func (self *ToneGenerator) SetAmplitude(amp int) {
	self.amp = float32(amp)
}

func (self *ToneGenerator) Close() {
	C.closeAudioDevice(self.dev)
}

func (self *ToneGenerator) feedSamples(data []float32) {
	for i:= 0; i < len(data); i++ {
		data[i] = self.amp * float32(math.Sin(float64(self.v*2*math.Pi/FREQUENCY)))
		if self.j > self.period {
			self.v -= float32(self.period) * self.freq
			self.j = 0
		} else {
			self.v += self.freq
			self.j++
		}
	}
}

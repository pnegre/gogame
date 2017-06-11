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

type ToneGenerator struct {
	dev C.SDL_AudioDeviceID
}

func NewToneGenerator() (*ToneGenerator, error) {
	dev := C.newAudioDevice()
	if dev == 0 {
		return nil, errors.New("Can't open tone generator")
	}
	sd := new(ToneGenerator)
	sd.dev = dev
	return sd, nil
}

func (self *ToneGenerator) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *ToneGenerator) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *ToneGenerator) SetFreq(freq int) {
	C.setFreq(C.int(self.dev), C.int(freq))
}

func (self *ToneGenerator) SetAmplitude(amp int) {
	C.setAmplitude(C.int(self.dev), C.int(amp))
}

func (self *ToneGenerator) Close() {
	C.SDL_CloseAudioDevice(self.dev)
}

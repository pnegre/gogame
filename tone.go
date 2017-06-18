package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

extern SDL_AudioDeviceID newAudioDevice(int frequency);
extern void closeAudioDevice(SDL_AudioDeviceID id);
extern void mixAudio(void *src, void* dst, int len);

*/
import "C"
import "errors"
import "unsafe"


const (
	GENERATOR_TYPE_TONE = iota
	GENERATOR_TYPE_NOISE
)

var theCallback func([]int16)

//export soundGoCallback
func soundGoCallback(id int, ptr unsafe.Pointer, len int) {
	len /= 2
	// To create a Go slice backed by a C array
	// (without copying the original data), one needs to acquire this length
	// at runtime and use a type conversion to a pointer to a very big array and then
	// slice it to the length that you want (also remember to set
	// the cap if you're using Go 1.2 or later)
	slice := (*[1 << 30]int16)(ptr)[:len:len]
	if theCallback != nil {
		theCallback(slice)
	}
}

func RegisterSoundCallback(fnc func([]int16)) {
	theCallback = fnc
}

// func MixAudio(src unsafe.Pointer, dst unsafe.Pointer, lg int) {
// 	C.mixAudio(src, dst, C.int(lg))
// }

type AudioDevice struct {
	dev       C.SDL_AudioDeviceID
	frequency int
}

func NewAudioDevice(frequency int) (*AudioDevice, error) {
	dev := C.newAudioDevice(C.int(frequency))
	if dev == 0 {
		return nil, errors.New("Can't open audio device")
	}
	sd := new(AudioDevice)
	sd.dev = dev
	sd.frequency = frequency
	return sd, nil
}

func (self *AudioDevice) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *AudioDevice) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *AudioDevice) Close() {
	C.closeAudioDevice(self.dev)
}

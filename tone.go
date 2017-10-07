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

type funcCallbackType func([]int16)

var mapa map[int]*AudioDevice

func init() {
	mapa = make(map[int]*AudioDevice)
}

//export soundGoCallback
func soundGoCallback(id int, ptr unsafe.Pointer, len int) {
	len /= 2
	// To create a Go slice backed by a C array
	// (without copying the original data), one needs to acquire this length
	// at runtime and use a type conversion to a pointer to a very big array and then
	// slice it to the length that you want (also remember to set
	// the cap if you're using Go 1.2 or later)
	slice := (*[1 << 30]int16)(ptr)[:len:len]
	for k, v := range mapa {
		if k == id {
			v.callback(slice)
		}
	}
}

type AudioDevice struct {
	dev      C.SDL_AudioDeviceID
	callback funcCallbackType
}

func NewAudioDevice(frequency int) (*AudioDevice, error) {
	dev := C.newAudioDevice(C.int(frequency))
	if dev == 0 {
		return nil, errors.New("Can't open audio device")
	}
	sd := new(AudioDevice)
	sd.dev = dev
	return sd, nil
}

func (self *AudioDevice) SetCallback(fnc funcCallbackType) {
	self.callback = fnc
	mapa[int(self.dev)] = self
}

func (self *AudioDevice) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *AudioDevice) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *AudioDevice) Close() {
	C.SDL_PauseAudioDevice(self.dev, 1)
	C.closeAudioDevice(self.dev)
}

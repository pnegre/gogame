package gogame

/*
#cgo pkg-config: sdl2
#cgo LDFLAGS: -lm
#include "SDL.h"
#include <stdlib.h>
#include <stdio.h>
#include <math.h>

int FREQUENCY = 44100;

float vvs[100];
int amplitudesCallback[100];
int freqsCallback[100];

void audioCallback(void* userdata, Uint8* stream, int len) {
	int *id = (int*) userdata;
	float v = freqsCallback[*id];
	float *data = (float*) stream;
	for(int i=0; i<len / 4; i++) {
		data[i] = amplitudesCallback[*id] * sin(v * 2 * M_PI / FREQUENCY);
		v += freqsCallback[*id];
	}
}

SDL_AudioDeviceID newAudioDevice() {
	int *id = (int*) malloc(sizeof(int));
    SDL_AudioSpec want, have;
    SDL_AudioDeviceID dev;
    SDL_memset(&want, 0, sizeof(want));
    want.freq = FREQUENCY;
    want.format = AUDIO_F32;
    want.channels = 1;
    want.samples = 4096;
    want.callback = audioCallback;
	want.userdata = id;

    dev = SDL_OpenAudioDevice(NULL, 0, &want, &have, SDL_AUDIO_ALLOW_FORMAT_CHANGE);
    if (dev == 0) {
        SDL_Log("Failed to open audio: %s", SDL_GetError());
    }
    printf("freq: %d\n", have.freq);

	*id = dev;
    return dev;
}

*/
import "C"
import "errors"

type SoundDevice struct {
	dev C.SDL_AudioDeviceID
}

func NewSoundDevice() (*SoundDevice, error) {
	dev := C.newAudioDevice()
	if dev == 0 {
		return nil, errors.New("Can't open audio device")
	}
	sd := new(SoundDevice)
	sd.dev = dev
	return sd, nil
}

func (self *SoundDevice) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *SoundDevice) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *SoundDevice) SetFreq(freq int) {
	C.freqsCallback[self.dev] = C.int(freq)
}

func (self *SoundDevice) SetAmplitude(amp int) {
	C.amplitudesCallback[self.dev] = C.int(amp)
}

func (self *SoundDevice) Close() {
	C.SDL_CloseAudioDevice(self.dev)
}

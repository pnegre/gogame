package gogame

/*
#cgo pkg-config: sdl2
#cgo LDFLAGS: -lm
#include "SDL.h"
#include <stdlib.h>
#include <stdio.h>
#include <math.h>

#define MAXTONEGENERATORS 100

int FREQUENCY = 44100;

float vvs[MAXTONEGENERATORS];
int amplitudesCallback[MAXTONEGENERATORS];
int freqsCallback[MAXTONEGENERATORS];

void audioCallback(void* userdata, Uint8* stream, int len) {
	int *id = (int*) userdata;
	float v = vvs[*id];
	float freq = freqsCallback[*id];
	float amp = amplitudesCallback[*id];
	float *data = (float*) stream;
	int N = (int) (2*FREQUENCY/freq);
	for(int i=0, j=0; i<len / 4; i++) {
		data[i] = amp * sin(v * 2 * M_PI / FREQUENCY);
		if (j > N) {
			v -= N*freq;
			j = 0;
		} else {
			v += freq;
			j++;
		}
	}
	vvs[*id] = v;
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
	if (dev >= MAXTONEGENERATORS) {
		SDL_Log("Error: MAXTONEGENERATORS reached!!");
		return 0;
	}

	*id = dev;
    return dev;
}

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
	C.freqsCallback[self.dev] = C.int(freq)
}

func (self *ToneGenerator) SetAmplitude(amp int) {
	C.amplitudesCallback[self.dev] = C.int(amp)
}

func (self *ToneGenerator) Close() {
	C.SDL_CloseAudioDevice(self.dev)
}

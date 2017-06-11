#include "SDL.h"
#include <stdlib.h>
#include <stdio.h>
#include <math.h>

#include "_cgo_export.h"

#define MAXTONEGENERATORS 100

int FREQUENCY = 44100;

float vvs[MAXTONEGENERATORS];
int amplitudesCallback[MAXTONEGENERATORS];
int freqsCallback[MAXTONEGENERATORS];

void setAmplitude(int id, int amp) {
	amplitudesCallback[id] = amp;
}

void setFreq(int id, int freq) {
	freqsCallback[id] = freq;
}

void audioCallback(void* userdata, Uint8* stream, int len) {
	int *id = (int*) userdata;
	// float v = vvs[*id];
	// float freq = freqsCallback[*id];
	// float amp = amplitudesCallback[*id];
	// float *data = (float*) stream;
	// int N = (int) (2*FREQUENCY/freq);
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
	// vvs[*id] = v;
	soundGoCallback(*id, stream, len);
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

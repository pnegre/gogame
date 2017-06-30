#include "SDL.h"
#include <stdlib.h>
#include <stdio.h>
#include <math.h>

#include "_cgo_export.h"

void audioCallback(void* userdata, Uint8* stream, int len) {
	int *id = (int*) userdata;
	soundGoCallback(*id, stream, len);
}

SDL_AudioDeviceID newAudioDevice(int frequency) {
	int *id = (int*) malloc(sizeof(int));
    SDL_AudioSpec want, have;
    SDL_AudioDeviceID dev;
    SDL_memset(&want, 0, sizeof(want));
    want.freq = frequency;
    want.format = AUDIO_S16; // Signed 16 bits
    want.channels = 1;
    want.samples = 256;
    want.callback = audioCallback;
	want.userdata = id;

    dev = SDL_OpenAudioDevice(NULL, 0, &want, &have, SDL_AUDIO_ALLOW_FORMAT_CHANGE);
	printf("Freq: %d\n", have.freq);
    if (dev == 0) {
        SDL_Log("Failed to open audio: %s", SDL_GetError());
    }

	*id = dev;
    return dev;
}

void closeAudioDevice(SDL_AudioDeviceID dev) {
	SDL_CloseAudioDevice(dev);
	// TODO: free de l'integer id
}

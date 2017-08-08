#include "SDL.h"

int initSDL() {
	if (SDL_Init(SDL_INIT_AUDIO | SDL_INIT_VIDEO) != 0) {
		SDL_Log("Unable to initialize SDL: %s", SDL_GetError());
		return 1;
	}
	return 0;
}

SDL_Window * newScreen(char *title, int h, int v) {
    return SDL_CreateWindow(title, SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, h, v, SDL_WINDOW_RESIZABLE);
}

SDL_Renderer * newRenderer( SDL_Window * screen ) {
    SDL_Renderer * r = SDL_CreateRenderer(screen, -1, SDL_RENDERER_PRESENTVSYNC | SDL_RENDERER_ACCELERATED); // SDL_RENDERER_SOFTWARE ); // SDL_RENDERER_ACCELERATED  );
	return r;
}

void setScaleQuality(int n) {
	switch(n) {
	case 1:
		SDL_SetHint(SDL_HINT_RENDER_SCALE_QUALITY, "1");
		break;
	case 2:
		SDL_SetHint(SDL_HINT_RENDER_SCALE_QUALITY, "2");
	}
}

int isNull(void *pointer) {
	if (pointer == NULL) {
		return 1;
	} else {
		return 0;
	}
}

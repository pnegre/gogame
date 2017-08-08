#include "SDL.h"
#include "SDL_image.h"
#include <stdlib.h>
#include <stdio.h>

SDL_Texture * makeTexture( char *f, SDL_Renderer *ren ) {
    SDL_Texture *tex = IMG_LoadTexture(ren, f);
    return tex;
}

void renderGOOD( SDL_Renderer *ren, SDL_Texture *tex, int ox, int oy, int x, int y, int w, int h, int dw, int dh) {
    static SDL_Rect org;
    static SDL_Rect dst;
    org.x = ox;
    org.y = oy;
    org.w = w;
    org.h = h;
    dst.x = x;
    dst.y = y;
    dst.w = dw;
    dst.h = dh;
    SDL_RenderCopy(ren, tex, &org, &dst);
}

void queryTexture(SDL_Texture *t, int *h, int *v) {
    SDL_QueryTexture(t, NULL, NULL, h, v);
}

int intersects(int x1, int y1, int w1, int h1, int x2, int y2, int w2, int h2) {
    static SDL_Rect a;
    static SDL_Rect b;
    a.x = x1; a.y = y1; a.w = w1; a.h = h1;
    b.x = x2; b.y = y2; b.w = w2; b.h = h2;
    return SDL_HasIntersection(&a, &b);
}

SDL_Texture *makeEmptyTexture(SDL_Renderer *ren, int w, int h) {
	SDL_Texture *t = SDL_CreateTexture(ren, SDL_PIXELFORMAT_RGB24, SDL_TEXTUREACCESS_STREAMING, w, h);
	if (t == NULL) {
		printf("Error creating empty texture: %s\n", SDL_GetError());
	}
	return t;
}

unsigned char *lockTexture(SDL_Texture *t) {
	void *texture_data;
	int texture_pitch;
	if (SDL_LockTexture(t, 0, &texture_data, &texture_pitch) == -1) {
		printf("Error: %s\n", SDL_GetError());
	}
	return (unsigned char*) texture_data;
}

void unlockTexture(SDL_Texture *t) {
	SDL_UnlockTexture(t);
}

void pixel(unsigned char *data, int w, int h, int x, int y, int r, int g, int b) {
	data += (x+y*w)*3;
	*data++ = (unsigned char) r;
	*data++ = (unsigned char) g;
	*data++ = (unsigned char) b;
}

void clear(unsigned char *data, int w, int h) {
	for(int i=0; i<h; i++)
		for(int j=0; j<w; j++) {
			*data++ = 0;
			*data++ = 0;
			*data++ = 0;
		}
}

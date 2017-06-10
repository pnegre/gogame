/*

Package gogame is a set of functions and modules designed for writing games. Gogame uses SDL2 internally.

This software is free. It's released under the LGPL license. You can create open source and commercial games with it. See the license for full details.

OPENGL is required. Developer libraries of SDL2, SDL2-image and SDL2-TTF are required also.

*/
package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

int initSDL() {
	if (SDL_Init(SDL_INIT_AUDIO | SDL_INIT_VIDEO) != 0) {
		SDL_Log("Unable to initialize SDL: %s", SDL_GetError());
		return 1;
	}
	return 0;
}

SDL_Window * newScreen(char *title, int h, int v) {
    return SDL_CreateWindow(title, SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, h, v, SDL_WINDOW_OPENGL | SDL_WINDOW_RESIZABLE);
}

SDL_Renderer * newRenderer( SDL_Window * screen ) {
    return SDL_CreateRenderer(screen, -1, SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC );
}


*/
import "C"
import "errors"

var screen *C.SDL_Window
var renderer *C.SDL_Renderer

type Color struct {
	R, G, B, A uint8
}

var COLOR_WHITE = &Color{255, 255, 255, 255}
var COLOR_BLACK = &Color{0, 0, 0, 255}
var COLOR_RED = &Color{255, 0, 0, 255}
var COLOR_BLUE = &Color{0, 0, 255, 255}

// Use this function to create a window and a renderer (not visible to user)
func Init(title string, h, v int) error {
	if i := C.initSDL(); i != 0 {
		return errors.New("Error initializing SDL")
	}
	screen = C.newScreen(C.CString(title), C.int(h), C.int(v))
	renderer = C.newRenderer(screen)
	if screen == nil || renderer == nil {
		return errors.New("Error on initializing SDL2")
	}
	return nil
}

// Full Screen mode
func SetFullScreen(fs bool) {
	if fs {
		C.SDL_SetWindowFullscreen(screen, C.SDL_WINDOW_FULLSCREEN_DESKTOP)
	} else {
		C.SDL_SetWindowFullscreen(screen, 0)
	}
}

// Get window size
func GetWindowSize() (int, int) {
	var w, h C.int
	C.SDL_GetWindowSize(screen, &w, &h)
	return int(w), int(h)
}

// Set window size
func SetWindowSize(h, v int) {
	C.SDL_SetWindowSize(screen, C.int(h), C.int(v))
}

// Set a device independent resolution for rendering
func SetLogicalSize(h, v int) {
	C.SDL_RenderSetLogicalSize(renderer, C.int(h), C.int(v))
}

// Destroys renderer and window
func Quit() {
	C.SDL_DestroyRenderer(renderer)
	C.SDL_DestroyWindow(screen)
	C.SDL_Quit()
}

// Clear the current rendering target with black color
func RenderClear() {
	C.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 0)
	C.SDL_RenderClear(renderer)
}

// Update the screen with rendering performed
func RenderPresent() {
	C.SDL_RenderPresent(renderer)
}

// Wait specified number of milliseconds before returning
func Delay(s int) {
	C.SDL_Delay(C.Uint32(s))
}

// Draw pixel at position x,y
func DrawPixel(x, y int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	C.SDL_RenderDrawPoint(renderer, C.int(x), C.int(y))
}

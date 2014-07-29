package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"


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

func Init(title string, h, v int) error {
	screen = C.newScreen(C.CString(title), C.int(h), C.int(v))
	renderer = C.newRenderer(screen)
	if screen == nil || renderer == nil {
		return errors.New("Error on initializing SDL2")
	}
	return nil
}

func SetLogicalSize(h, v int) {
	C.SDL_RenderSetLogicalSize(renderer, C.int(h), C.int(v))
}

func Quit() {
	C.SDL_DestroyRenderer(renderer)
	C.SDL_DestroyWindow(screen)
	C.SDL_Quit()
}

func RenderClear() {
	C.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 0)
	C.SDL_RenderClear(renderer)
}

func RenderPresent() {
	C.SDL_RenderPresent(renderer)
}

func Delay(s int) {
	C.SDL_Delay(C.Uint32(s))
}

func DrawPixel(x, y int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	C.SDL_RenderDrawPoint(renderer, C.int(x), C.int(y))
}

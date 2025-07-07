/*

Package gogame is a set of functions and modules designed for writing games. Gogame uses SDL2 internally.

This software is free. It's released under the LGPL license. You can create open source and commercial games with it. See the license for full details.

OPENGL is required. Developer libraries of SDL2, SDL2-image and SDL2-TTF are required also.

*/
package gogame

/*
#cgo pkg-config: sdl2
#cgo pkg-config: SDL2_gfx
#include "SDL.h"
#include "SDL2_gfxPrimitives.h"

extern int initSDL();
extern SDL_Window * newScreen(char *title, int h, int v);
extern SDL_Renderer * newRenderer( SDL_Window * screen );
extern void setScaleQuality(int n);
extern int isNull(void *pointer);
extern void getDesktopDisplayResolution(int displayIndex, int *w, int *h);
*/
import "C"
import (
	"errors"
	"unsafe"
)

var screen *C.SDL_Window
var renderer *C.SDL_Renderer

type Color struct {
	R, G, B, A uint8
}

var COLOR_WHITE = &Color{255, 255, 255, 255}
var COLOR_BLACK = &Color{0, 0, 0, 255}
var COLOR_RED = &Color{255, 0, 0, 255}
var COLOR_BLUE = &Color{0, 0, 255, 255}

// InitSDL initializes the SDL library. It must be called before any other function in this package.
// It returns an error if SDL could not be initialized.
func InitSDL() error {
	if i := C.initSDL(); i != 0 {
		return errors.New("Error initializing SDL")
	}
	return nil
}

// Use this function to create a window and a renderer (not visible to user)
func Init(title string, h, v int) error {
	screen = C.newScreen(C.CString(title), C.int(h), C.int(v))
	if C.isNull(unsafe.Pointer(screen)) == 1 {
		return errors.New("Error initalizing SCREEN")
	}
	renderer = C.newRenderer(screen)
	if C.isNull(unsafe.Pointer(renderer)) == 1 {
		return errors.New("Error initalizing RENDERER")
	}
	if screen == nil || renderer == nil {
		return errors.New("Error on initializing SDL2")
	}
	return nil
}

func SetScaleQuality(n int) {
	C.setScaleQuality(C.int(n))
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

func GetDesktopDisplayResolution() (int, int) {
	var w, h C.int
	C.getDesktopDisplayResolution(0, &w, &h)
	return int(w), int(h)
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

// Draw line
func DrawLine(x1, y1, x2, y2 int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	C.SDL_RenderDrawLine(renderer, C.int(x1), C.int(y1), C.int(x2), C.int(y2))
}

// Draw rectangle
func DrawRect(x1, y1, x2, y2 int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	r := &C.SDL_Rect{x: C.int(x1), y: C.int(y1), w: C.int(x2 - x1), h: C.int(y2 - y1)}
	C.SDL_RenderDrawRect(renderer, r)
}

// Draw filled triangle
func DrawFilledTriangle(x1, y1, x2, y2, x3, y3 int, color *Color) {
	C.filledTrigonRGBA(renderer, C.short(x1), C.short(y1),
		C.short(x2), C.short(y2), C.short(x3), C.short(y3),
		C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
}

// Draw filled rectangle
func DrawFilledRectangle(x1, y1, x2, y2 int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	r := &C.SDL_Rect{x: C.int(x1), y: C.int(y1), w: C.int(x2 - x1), h: C.int(y2 - y1)}
	C.SDL_RenderFillRect(renderer, r)
}

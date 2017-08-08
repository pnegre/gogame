package gogame

/*
#cgo pkg-config: sdl2 SDL2_image
#include "SDL.h"
#include "SDL_image.h"
#include <stdlib.h>
#include <stdio.h>

extern SDL_Texture * makeTexture( char *f, SDL_Renderer *ren );
extern void renderGOOD( SDL_Renderer *ren, SDL_Texture *tex, int ox, int oy, int x, int y, int w, int h, int dw, int dh);
extern void queryTexture(SDL_Texture *t, int *h, int *v);
extern int intersects(int x1, int y1, int w1, int h1, int x2, int y2, int w2, int h2);
extern SDL_Texture *makeEmptyTexture(SDL_Renderer *ren, int w, int h);
extern unsigned char *lockTexture(SDL_Texture *t);
extern void unlockTexture(SDL_Texture *t);
extern void pixel(unsigned char *data, int w, int h, int x, int y, int r, int g, int b);
extern void clear(unsigned char *data, int w, int h);
extern int isNull(void *pointer);


*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Drawable interface {
	BlitRect(*Rect)
	GetDimensions() (int, int)
}

type Texture struct {
	tex   *C.SDL_Texture
	realw int
	realh int
	dstw  int
	dsth  int
	data  *C.uchar
}

type Rect struct {
	X, Y, W, H int
}

// Set center of Rect to x,y
func (self *Rect) SetCenter(x, y int) {
	self.X = x - self.W/2
	self.Y = y - self.H/2
}

func (self *Rect) GetCenter() (x, y int) {
	return self.X + self.W/2, self.Y + self.H/2
}

// Determine if two rectangles intersect
func (self *Rect) Intersects(r2 *Rect) bool {
	if 0 == C.intersects(C.int(self.X), C.int(self.Y), C.int(self.W), C.int(self.H),
		C.int(r2.X), C.int(r2.Y), C.int(r2.W), C.int(r2.H)) {
		return false
	}

	return true
}

// Construct a new texture from a image file
func NewTexture(filename string) *Texture {
	tex := C.makeTexture(C.CString(filename), renderer)
	return getNewTexture(tex)
}

func NewEmptyTexture(w, h int) (*Texture, error) {
	tex := C.makeEmptyTexture(renderer, C.int(w), C.int(h))
	if C.isNull(unsafe.Pointer(tex)) == 1 {
		return nil, fmt.Errorf("Error creating texture")
	}
	return getNewTexture(tex), nil
}

func getNewTexture(tex *C.SDL_Texture) *Texture {
	t := new(Texture)
	t.tex = tex
	var w C.int
	var h C.int
	C.queryTexture(t.tex, &w, &h)
	t.realw, t.realh = int(w), int(h)
	t.dstw, t.dsth = t.realw, t.realh
	return t
}

func (self *Texture) Lock() {
	self.data = C.lockTexture(self.tex)
}

func (self *Texture) Unlock() {
	C.unlockTexture(self.tex)
}

func (self *Texture) Clear() {
	C.clear(self.data, C.int(self.realw), C.int(self.realh))
}

func (self *Texture) Pixel(x, y int, color *Color) {
	if x < 0 || x > self.realw || y < 0 || y > self.realh {
		return
	}
	C.pixel(self.data, C.int(self.realw), C.int(self.realh), C.int(x), C.int(y), C.int(color.R), C.int(color.G), C.int(color.B))
}

func (self *Texture) SetDimensions(w, h int) {
	self.dstw = w
	self.dsth = h
}

// Destroy texture. Must be called explicitly, no automatic free for texture data
func (self *Texture) Destroy() {
	C.SDL_DestroyTexture(self.tex)
}

// Get texture dimensions, (horizontal, vertical)
func (self *Texture) GetDimensions() (int, int) {
	return self.dstw, self.dsth
}

func (self *Texture) Blit(x, y int) {
	C.renderGOOD(renderer, self.tex, C.int(0), C.int(0), C.int(x), C.int(y), C.int(self.realw),
		C.int(self.realh), C.int(self.dstw), C.int(self.dsth))
}

// Blit texture to screen, using provided rect
func (self *Texture) BlitRect(r *Rect) {
	C.renderGOOD(renderer, self.tex, C.int(0), C.int(0), C.int(r.X), C.int(r.Y), C.int(self.realw),
		C.int(self.realh), C.int(r.W), C.int(r.H))
}

// Get subtexture
func (self *Texture) SubTex(x, y, w, h int) *SubTexture {
	return &SubTexture{self, &Rect{x, y, w, h}, w, h}
}

type SubTexture struct {
	tex        *Texture
	rect       *Rect
	dstw, dsth int
}

func (self *SubTexture) SetDimensions(w, h int) {
	self.dstw = w
	self.dsth = h
}

func (self *SubTexture) Blit(x, y int) {
	C.renderGOOD(renderer, self.tex.tex, C.int(self.rect.X), C.int(self.rect.Y), C.int(x), C.int(y), C.int(self.rect.W),
		C.int(self.rect.H), C.int(self.dstw), C.int(self.dsth))
}

// Blit subtexture to screen, using provided rect
func (self *SubTexture) BlitRect(r *Rect) {
	C.renderGOOD(renderer, self.tex.tex, C.int(self.rect.X), C.int(self.rect.Y), C.int(r.X), C.int(r.Y), C.int(self.rect.W),
		C.int(self.rect.H), C.int(r.W), C.int(r.H))
}

// Get subtexture dimensions
func (self *SubTexture) GetDimensions() (int, int) {
	return self.rect.W, self.rect.H
}

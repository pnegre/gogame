package gogame

/*
#cgo pkg-config: sdl2 SDL2_image
#include "SDL.h"
#include "SDL_image.h"

SDL_Texture * makeTexture( char *f, SDL_Renderer *ren ) {
    SDL_Texture *tex = IMG_LoadTexture(ren, f);
    return tex;
}

void render( SDL_Renderer *ren, SDL_Texture *tex, int x, int y, int w, int h) {
    static SDL_Rect dst;
    dst.x = x;
    dst.y = y;
    dst.w = w;
    dst.h = h;
    SDL_RenderCopy(ren, tex, NULL, &dst);
}

void render2( SDL_Renderer *ren, SDL_Texture *tex, int x, int y, int w, int h, int ox, int oy) {
    static SDL_Rect org;
    static SDL_Rect dst;
    dst.x = x;
    dst.y = y;
    dst.w = w;
    dst.h = h;
    org.x = ox;
    org.y = oy;
    org.w = w;
    org.h = h;
    SDL_RenderCopy(ren, tex, &org, &dst);


}

void queryTexture(SDL_Texture *t, int *h, int *v) {
    SDL_QueryTexture(t, NULL, NULL, h, v);
}


*/
import "C"

type Drawable interface {
	BlitRect(*Rect)
	GetDimensions() (int, int)
}

type Texture struct {
	tex *C.SDL_Texture
	w   int
	h   int
}

type Rect struct {
	X, Y, W, H int
}

// Set center of Rect to x,y
func (self *Rect) SetCenter(x, y int) {
	self.X = x - self.W/2
	self.Y = y - self.H/2
}

// Construct a new texture from a image file
func NewTexture(filename string) *Texture {
	tex := C.makeTexture(C.CString(filename), renderer)
	return getNewTexture(tex)
}

func getNewTexture(tex *C.SDL_Texture) *Texture {
	t := new(Texture)
	t.tex = tex
	var w C.int
	var h C.int
	C.queryTexture(t.tex, &w, &h)
	t.w, t.h = int(w), int(h)
	// El finalizer no funciona. Suposo, fent proves, que és perquè el Garbage Collector
	// s'executa dins la seva pròpia goroutine, i sembla que no es pot cridar a
	// SDL_DestroyTexture si no és des del thread original
	// 	runtime.SetFinalizer(t, func(t *Texture) {
	// 		log.Printf("Finalizing texture...")
	// 		C.SDL_DestroyTexture(t.tex)
	// 	})
	return t
}

// Destroy texture. Must be called explicitly, no automatic free for texture data
func (self *Texture) Destroy() {
	C.SDL_DestroyTexture(self.tex)
}

// Get texture dimensions, (horizontal, vertical)
func (self *Texture) GetDimensions() (int, int) {
	return self.w, self.h
}

// Blit texture to screen, using provided rect
func (self *Texture) BlitRect(r *Rect) {
	C.render(renderer, self.tex, C.int(r.X), C.int(r.Y),
		C.int(r.W), C.int(r.H))
}

// Blit texture to screen, using x,y coordinates (top-left corner)
func (self *Texture) BlitToScreen(x, y int) {
	C.render(renderer, self.tex, C.int(x), C.int(y),
		C.int(self.w), C.int(self.h))
}

// Get subtexture
func (self *Texture) SubTex(x, y, w, h int) *SubTexture {
	return &SubTexture{self, &Rect{x, y, w, h}}
}

type SubTexture struct {
	tex  *Texture
	rect *Rect
}

// Blit subtexture to screen, using provided rect
func (self *SubTexture) BlitRect(r *Rect) {
	C.render2(renderer, self.tex.tex, C.int(r.X), C.int(r.Y), C.int(r.W),
		C.int(r.H), C.int(self.rect.X), C.int(self.rect.Y))
}

// Get subtexture dimensions
func (self *SubTexture) GetDimensions() (int, int) {
	return self.rect.W, self.rect.H
}

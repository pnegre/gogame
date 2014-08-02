package cache

import "github.com/pnegre/gogame"
import "log"

var texCache = make(map[string]*gogame.Texture)
var fontCache = make(map[namesize]*gogame.Font)

// Gets texture from cache. Creates it if necessary.
func GetTexture(name string) *gogame.Texture {
	if v, ok := texCache[name]; ok {
		return v
	}

	tex := gogame.NewTexture(name)
	texCache[name] = tex
	log.Printf("Loaded new texture: %s", name)
	return tex
}

// Destroy all textures from cache
func DestroyTextures() {
	for _, v := range texCache {
		log.Printf("Destroying texture %v", v)
		v.Destroy()
	}
}

type namesize struct {
	name string
	size int
}

// Gets font from cache. Creates it if necessary.
func GetFont(name string, size int) *gogame.Font {
	ns := namesize{name, size}
	if v, ok := fontCache[ns]; ok {
		return v
	}

	f := gogame.NewFont(name, size)
	fontCache[ns] = f
	log.Printf("Loaded new font %s with size %d", name, size)
	return f
}

// Destroy all fonts from cache
func DestroyFonts() {
	for k, v := range fontCache {
		log.Printf("Destroying font %v", k)
		v.Destroy()
	}
}

// Destroy all textures and fonts from cache
func DestroyAll() {
	DestroyFonts()
	DestroyTextures()
}

# GOGAME

This is a simple 2D game library for Go. It's incomplete and I'm writing it for my personal projects. Feel free to suggest features!!

## Install instructions (linux only)

I use debian. If you're like me, you will need to install the following packages:

    libsdl2-dev
    libsdl2-image-dev
    libsdl2-ttf-dev

Others distros shoud have equivalent libraries avaliable.

And now, the easy part (make sure you have a valid $GOPATH):

    go get github.com/pnegre/gogame

Finally, compile the library and install it:

    go install github.com/pnegre/gogame

## Documentation

Once you have installed the package, just run:

    godoc github.com/pnegre/gogame | less
    godoc github.com/pnegre/gogame/sprite | less
    godoc github.com/pnegre/gogame/cache | less

Or you may prefer to view on the web browser:

    godoc -http=:6060

Point your browser to http://localhost:6060/pkg/github.com/pnegre/gogame

## Usage

Simple example:

    package main

    import (
        "github.com/pnegre/gogame"
        "github.com/pnegre/gogame/cache"
        "github.com/pnegre/gogame/sprite"
        "log"
    )

    const (
        WINTITLE = "test"
        WIN_W    = 800
        WIN_H    = 600
        IMAGE    = "someimage.png"
    )

    type Target struct {
        *sprite.Simple
    }

    func NewTarget() *Target {
        tex := cache.GetTexture(IMAGE)
        s := sprite.NewSimple(tex)
        s.Rect.SetCenter(WIN_W/2, WIN_H/2)
        return &Target{s}
    }

    func (self *Target) Update() {
        if gogame.IsKeyPressed(gogame.K_LEFT) {
            self.Rect.X -= 10
        }
        if gogame.IsKeyPressed(gogame.K_RIGHT) {
            self.Rect.X += 10
        }
    }

    func main() {
        if err := gogame.Init(WINTITLE, WIN_W, WIN_H); err != nil {
            log.Fatal(err)
        }
        defer gogame.Quit()
        defer cache.DestroyAll()

        target := NewTarget()

        for {
            if quit := gogame.SlurpEvents(); quit == true {
                break
            }

            target.Update()

            gogame.RenderClear()
            target.Draw()
            gogame.RenderPresent()

            gogame.Delay(50)
        }

    }



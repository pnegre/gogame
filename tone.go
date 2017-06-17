package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

extern SDL_AudioDeviceID newAudioDevice(int frequency);
extern void closeAudioDevice(SDL_AudioDeviceID id);

*/
import "C"
import "errors"
import "unsafe"
// import "math"
import "math/rand"

const FREQUENCY = 44100
const INVFREQUENCY = float32(1) / float32(FREQUENCY)

var toneCache = make(map[int]*ToneGenerator)

const (
	GENERATOR_TYPE_TONE = iota
	GENERATOR_TYPE_NOISE
)

//export soundGoCallback
func soundGoCallback(id int, ptr unsafe.Pointer, len int) {
	len /= 4
	// To create a Go slice backed by a C array
	// (without copying the original data), one needs to acquire this length
	// at runtime and use a type conversion to a pointer to a very big array and then
	// slice it to the length that you want (also remember to set
	// the cap if you're using Go 1.2 or later)
	slice := (*[1 << 30]float32)(ptr)[:len:len]
	toneCache[id].feedSamples(slice)
}

type ToneGenerator struct {
	genType int
	dev     C.SDL_AudioDeviceID
	amp     float32
	freq    float32
	v       float32
	period  int
	j       int
	count int
}

func NewToneGenerator(genType int) (*ToneGenerator, error) {
	dev := C.newAudioDevice(FREQUENCY)
	if dev == 0 {
		return nil, errors.New("Can't open tone generator")
	}
	sd := new(ToneGenerator)
	sd.dev = dev
	sd.genType = genType
	toneCache[int(sd.dev)] = sd
	return sd, nil
}

func (self *ToneGenerator) Start() {
	C.SDL_PauseAudioDevice(self.dev, 0)
}

func (self *ToneGenerator) Stop() {
	C.SDL_PauseAudioDevice(self.dev, 1)
}

func (self *ToneGenerator) SetFreq(freq float32) {
	self.freq = freq
	self.period = int(0x10000 * self.freq / FREQUENCY)
}

func (self *ToneGenerator) SetAmplitude(amp float32) {
	self.amp = amp
}

func (self *ToneGenerator) Close() {
	C.closeAudioDevice(self.dev)
	delete(toneCache, int(self.dev))
}

func (self *ToneGenerator) feedSamples(data []float32) {
	if self.genType == GENERATOR_TYPE_TONE {

		/*

				Tret de EMULIB


				if(WaveCH[J].Freq>=SndRate/2) break;
		          K=0x10000*WaveCH[J].Freq/SndRate;
		          L1=WaveCH[J].Count;

				  for(I=0;I<Samples;I++,L1+=K)
		          {
		            L2 = L1+K;
		            A1 = L1&0x8000? 127:-128;
		            if((L1^L2)&0x8000)
		              A1=A1*(0x8000-(L1&0x7FFF)-(L2&0x7FFF))/K;
		            Wave[I]+=A1*V;
		          }
		          WaveCH[J].Count=L1&0xFFFF;

		*/

		K := int(0x10000*self.freq/FREQUENCY);
		L1:=self.count;
		var A1 int
		for i := 0; i < len(data); i, L1 = i+1, L1+K {
			L2 := L1+K;
			if L1&0x8000 != 0 {
				A1 = 127
			} else {
				A1 = -128
			}
			if((L1^L2)&0x8000 != 0) {
				A1=A1*(0x8000-(L1&0x7FFF)-(L2&0x7FFF))/K;
			}
			data[i]=float32(A1)*self.amp;

			// data[i] = self.amp * float32(math.Sin(float64(self.v*2*math.Pi*self.freq)))
			// self.v += INVFREQUENCY
			// if self.j > self.period {
			// 	self.v -= float32(self.j) * INVFREQUENCY
			// 	self.j = 0
			// } else {
			// 	self.j++
			// }
		}
		self.count = L1&0xFFFF
	} else if self.genType == GENERATOR_TYPE_NOISE {
		for i := 0; i < len(data); i++ {
			data[i] = self.amp * rand.Float32()
		}
	}
}

package main

// #cgo CFLAGS: -I /usr/include/flite/
// #cgo LDFLAGS: -lflite -lflite_cmu_us_kal
// #include "flite.h"
// cst_voice* register_cmu_us_kal(const char *voxdir);
import "C"
import (
	"flag"
	"fmt"
	"unsafe"
)

var voice *C.cst_voice
var path, speech string

func init() {
	C.flite_init()
	voice = C.register_cmu_us_kal(nil)

	flag.StringVar(&path, "path", "hello.wav", "file path")
	flag.StringVar(&speech, "speech", "hello", "say the speech")
}

func main() {
	flag.Parse()
	if err := textToSpeech(path, speech); err != nil {
		fmt.Println("err:", err)
	}
}

func textToSpeech(path, text string) error {
	if voice == nil {
		return fmt.Errorf("could not find default voice")
	}
	ctext := C.CString(text)
	cout := C.CString(path)
	C.flite_text_to_speech(ctext, voice, cout)
	C.free(unsafe.Pointer(ctext))
	C.free(unsafe.Pointer(cout))
	return nil
}

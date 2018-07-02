package main

// extern void SayHello(_GoString_ s);
import "C"
import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	defer trace.Stop()

	go C.SayHello("Hello, World\n")
	time.Sleep(time.Millisecond)
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}

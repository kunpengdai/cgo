package main

//#include <stdlib.h>
import "C"
import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter generate string num:")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		fmt.Println(text)
		num, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			fmt.Printf("wrong num input,with err:%v", err)
			os.Exit(1)
		}
		if num == 0 {
			os.Exit(0)
		}
		memLeak(int(num))
	}

}

func memLeak(n int) {
	for i := 0; i < n; i++ {
		str := randStringRunes(100)
		// no memleak with free
		// cs := C.CString(str)
		// defer C.free(unsafe.Pointer(cs))

		//mem leak
		_ = C.CString(str)
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("memStatus : sys:%d,heapSys:%d,heapAlloc:%d,heapIdle:%d,heapReleased:%d\n", m.Sys, m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

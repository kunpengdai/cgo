# cgo编程初识

## 1. why cgo

1. GO语言有自己的擅长的领域 [web后端,分布式,区块链] ,但许多传统领域仍是C的主场
1. 通过cgo可以继承C语言近半个世纪的软件积累
1. cgo是GO语言直接和其他语言通讯的桥梁,C=>GO ,GO=>C
1. BUT:可以直接用纯粹的GO语言解决的问题不用cgo,能让其他组提供RPC调用的话,不使用cgo [fix one problem，cause another]

## 2. 快速入门
### 2.1 小例子
``` golang
package main

//#include <stdio.h>
import "C"

func main() {
	C.puts(C.CString("Hello world!\n"))
}
```

import "C" 

> When the Go tool sees that one or more Go files use the special import "C", it will look for other non-Go files in the directory and compile them as part of the Go package.

调用C函数，C.cfunc(C.cparams..)

tips:

* import "C"  与其上的声明语句之间不要有空行，C函数传递的需要是C类型的参数。
* cgo支持C,C++,Fortran语言

### 2.2 又一个小例子
 
 ``` golang
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
	defer C.free(unsafe.Pointer(ctext))
	defer C.free(unsafe.Pointer(cout))
	C.flite_text_to_speech(ctext, voice, cout)
	return nil
}
 ```

### 2.3 类型转换
C语言中可以被golang引用的类型有：
```
C.char, C.schar (signed char), C.uchar (unsigned char), C.short, C.ushort (unsigned short), C.int, C.uint (unsigned int), C.long, C.ulong (unsigned long), C.longlong (long long), C.ulonglong (unsigned long long), C.float, C.double, C.complexfloat (complex float), and C.complexdouble (complex double).
特别的，void* 可以用 Go的unsafe.pointer 来表示，__int128 和 __uint128 可以用[16]byte来表示。
struct,union,enum类型，需要添加struct_,union_,enum_
```
示例：
``` go
package main

/*#include <stdlib.h>
#include <stdio.h>
void printNum(int a){
	printf("input num:%d\n",a);
}

void printStr(char* c){
	printf("input string:%s\n",c);
}

void square(int *s,int n) {
    for(int i=1;i<=n;i++){
		s[i-1] = i*i;
	}
}
*/
import (
	"C"
)
import (
	"fmt"
	"unsafe"
)

func main() {
	//int
	num := 42
	C.printNum(C.int(num))

	//string
	str := "test"
	cstr := C.CString(str) //go string to c pointer char*
	defer C.free(unsafe.Pointer(cstr))
	C.printStr(cstr)

	//array point
	var ints [10]int32
	C.square((*C.int)(unsafe.Pointer(&ints[0])), 10)
	fmt.Println(ints)

}

```

指针：注意指针指向数据的size;golang中int一般64位，直接转成c中的int是不行的。

## 3. C=>GO GO=>C
利用cgo，不止可以在go中调用c代码，也可以导出方法给c来使用;
比如可以导出日志打印函数来供c代码打印日志;

golang中导出日志打印方法供c代码使用

``` go
//export go_log
func go_log(level C.int, msg *C.char) {
	levelInt := int(level)
	goMsg := C.GoString(msg)
	//c日志等级转化为golang日志
	switch levelInt {
	case 4, 5, 6:
		log.WInfo("cl", "record", "clog", goMsg)
	case 3:
		log.WWarn("cl", "record", "clog", goMsg)
	case 2:
		log.WError("cl", "record", "clog", goMsg)
	case 1, 0:
		log.WFatal("cl", "record", "clog", goMsg)
	default:
		log.WInfo("cl", "record", "clog", goMsg)
	}
}
```

在c的头文件中声明golang导出的方法

```c
#ifdef GOLANG
      extern void go_log(int level, char* msg);
#endif
```

* thread

## 4. 内存模型

### 4.1 内存管理
golang的内存管理器无法管理C代码中申请的内存，所以申请的c内存都需要手动进行释放。

内存释放示例：
``` go
str := "c mem"
cstr := C.CString(str) //go string to c pointer char*
defer C.free(unsafe.Pointer(cstr))
```

* 内存泄露示例：
``` golang
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

```
### 4.2 线程 goroutine
在cgo中，启动的c线程不同于goroutine，设置的runtime.GOMAXPROCS(cpuNum)只对goroutine有效，对c的线程没有限制作用，需要额外控制c并发的数量。

## 5. revover
在cgo中，golang引发的panic异常可以通过recover来捕获，但是c代码中引发的异常捕获不到。此问题会造成程序直接崩溃，需要注意下。

``` golang
package main

// int div(int a,int b)
// {return a/b;}
import "C"
import "fmt"

func main() {
	goPanic() //will recover
	cPanic()  //won't recover
}

func goPanic() {
	defer func() {
		recover()
		fmt.Println("go recoverd!")
	}()
	panic("panic in go")
}

func cPanic() {
	defer func() {
		recover()
		fmt.Println("c recoverd!")
	}()
	fmt.Println("div:", C.div(1, 0))
}
```

## 6. 库引用 与 部署

通过使用编译参数可以控制编译行为。

``` golang
// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"
```
CFLAGS, CPPFLAGS, CXXFLAGS, FFLAGS[头文件路径] and LDFLAGS[库文件路径] 分别可以控制 c，c++，fortran

如果cgo需要引入大量的库，可以使用pkg-config方式来管理依赖库

``` go
// #include "librecorder/recorder.h"
// #cgo pkg-config: librecorder
```
设置PKG_CONFIG_PATH变量
export PKG_CONFIG_PATH=your-pkgconfig-path
pkg-config配置文件示例

```
prefix=/usr/local
exec_prefix=${prefix}
libdir=${prefix}/lib
includedir_old=${prefix}/include
includedir_new=${prefix}/include/librecorder

Name: librecorder
Description: FFMPEG pull stream
Version: 0.9.9
Requires: libavformat,libavcodec,libswscale,libavutil
Requires.private: 
Conflicts:
Libs: -L${libdir} -lrecorder -lavformat -lavcodec -lswscale -lswresample -lavutil -lm -lpthread
Libs.private: 
Cflags: -I${includedir_old} -I${includedir_new} -Wall -g -shared
```

## 7. References

[1]  [cgo](https://golang.org/cmd/cgo/)

[2]  [深入cgo编程](https://github.com/chai2010/gopherchina2018-cgo-talk)

[3]  [cgo is not go](https://dave.cheney.net/2016/01/18/cgo-is-not-go)

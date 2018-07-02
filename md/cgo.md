# cgo编程初识

## 1. why cgo

1. GO语言有自己的擅长的领域 [web后端,分布式] ,但许多传统领域仍是C的主场
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

### 2.2 又一个小例子,库调用
 * flite

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
}*/
import (
	"C"
)
import "unsafe"

func main() {
	num := 42
	C.printNum(C.int(num))
	str := "test"
	cstr := C.CString(str) //go string to c pointer char*
	defer C.free(unsafe.Pointer(cstr))
	C.printStr(cstr)
}

```

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

## 4. 内存管理

golang的内存管理器无法管理C代码中申请的内存，所以申请的c内存都需要手动进行释放。

内存释放示例：
``` go
str := "c mem"
cstr := C.CString(str) //go string to c pointer char*
defer C.free(unsafe.Pointer(cstr))
```

内存泄露示例：
``` go

```

## 5. 线程模型

## 6. revover

## 7. 库引用 与 部署

## 8. 项目介绍
### 8.1 连麦视频录制
### 8.2 截图服务
### 8.2.1 pipeline
### 8.2.2 channel监控

## 9. References

[1]  [cgo](https://golang.org/cmd/cgo/)

[2]  [深入cgo编程](https://github.com/chai2010/gopherchina2018-cgo-talk)

[3]  [cgo is not go](https://dave.cheney.net/2016/01/18/cgo-is-not-go)

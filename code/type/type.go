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

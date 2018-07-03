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

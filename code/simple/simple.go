package main

//#include <stdio.h>
import "C"

func main() {
	C.puts(C.CString("Hello world!\n"))
	C.printf(C.CString("%s%d"), C.CString("string:"), 1)
}

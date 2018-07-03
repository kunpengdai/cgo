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

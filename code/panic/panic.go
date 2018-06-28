package main

// int div(int a,int b)
// {return a/b;}
import "C"
import "fmt"

func main() {
	defer func() {
		recover()
		fmt.Println("recoverd!")
	}()
	fmt.Println("div:", C.div(1, 0))
}

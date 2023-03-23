package main

func main() {
	ready := make(chan struct{})
	close(ready)
	<-ready
}

package main

type Publisher interface {
	Publish(msg []byte) error
	Start()
	Shutdown(immediately bool)
	Cleanup()
}

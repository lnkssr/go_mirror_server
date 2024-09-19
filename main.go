package main

import (
	"io"
	"log"
	"net"
)

func mirrorConnection(src net.Conn, target string) {
	defer src.Close()

	dst, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to target: %v", err)
		return
	}
	defer dst.Close()

	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			log.Printf("Error copying from src to dst: %v", err)
		}
	}()
	_, err = io.Copy(src, dst)
	if err != nil {
		log.Printf("Error copying from dst to src: %v", err)
	}
}

func main() {
	adress_back := "localhost:8080"
	adress_front := "0.0.0.0:8080"

	listener, err := net.Listen("tcp", adress_front)
	if err != nil {
		log.Fatalf("Failed to set up listener: %v", err)
	}
	defer listener.Close()
	log.Println("Listening on , mirroring to", adress_back)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go mirrorConnection(conn, adress_back)
	}
}

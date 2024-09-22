package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func mirrorConnection(src net.Conn, target string) {
	defer src.Close()

	dst, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to target %s: %v", target, err)
		return
	}
	defer dst.Close()

	log.Printf("Connected to source %s and target %s", src.RemoteAddr(), dst.RemoteAddr())

	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			log.Printf("Error copying from source %s to target %s: %v", src.RemoteAddr(), dst.RemoteAddr(), err)
		}
	}()
	_, err = io.Copy(src, dst)
	if err != nil {
		log.Printf("Error copying from target %s to source %s: %v", dst.RemoteAddr(), src.RemoteAddr(), err)
	}
}

func main() {
	logFlag := flag.Bool("log", false, "Include/exclude logging to file")
	flag.Parse()

	adress_back := "localhost:8080"
	adress_front := "0.0.0.0:8080"

	if *logFlag {
		f, _ := os.Create("main.log")
		defer f.Close()
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	} else {
		log.SetOutput(os.Stdout)
	}

	listener, err := net.Listen("tcp", adress_front)
	if err != nil {
		log.Fatalf("Failed to set up listener on %s: %v", adress_front, err)
	}
	defer listener.Close()
	log.Println("Listening on", adress_front, "mirroring to", adress_back)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go mirrorConnection(conn, adress_back)
	}
}

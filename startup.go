package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var startup_success bool = false
var startup_port io.ReadWriteCloser

func startup(port io.ReadWriteCloser) {
	startup_port = port
	go readStartup(port)
	for startup_success == false {
		startup_port.Write([]byte("AT\r\n"))
		time.Sleep(time.Second)
	}
	fmt.Println("Modem detected!")
}

func readStartup(reader io.Reader) {
	//buffered_reader := bufio.NewReader(startup_port)
	var buf []byte
	for startup_success == false {
		var read []byte = make([]byte, 1)
		reader.Read(read)
		buf = append(buf, read[0])
		casted := string(buf)
		if strings.HasSuffix(casted, "OK\r\n") {
			//fmt.Println("SUCCESS")
			startup_success = true
		}
	}
}

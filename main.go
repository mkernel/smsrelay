package main

import (
	"fmt"
	"strings"

	"github.com/jacobsa/go-serial/serial"
)

func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/ttyS0",
		BaudRate:        115200,
		MinimumReadSize: 2,
		StopBits:        1,
		DataBits:        8,
	}

	port, err := serial.Open(options)
	if err != nil {
		fmt.Println("Unable to open serial port. Bailing out...")
		fmt.Print(err)
	} else {
		startup(port)
		go reader(port)
		check_return = true
		//as we now have a modem that talks to us, we should get to the point of telling it not to echo.
		port.Write([]byte("ATE0\r\n"))
		//fmt.Println(">ATE0")
		<-proceed
		//and we should go forward and fire up the network
		port.Write([]byte("AT+COPS=1\r\n"))
		//fmt.Println(">AT+COPS=1")
		<-proceed
		port.Write([]byte("AT+CMGF=1\r\n"))
		//fmt.Println(">AT+CMGF=1")
		<-proceed
		check_return = false
		for true {
			idx := <-incoming_sms
			fmt.Println("Incoming sms detected! Index: " + idx)
			collect_reply = true
			port.Write([]byte("AT+CMGR=" + idx + "\r\n"))
			content := <-collect
			header := strings.Split(content[1], ",")
			number := header[1]
			sms_content := content[2 : len(content)-2]
			sms := strings.Join(sms_content, "\n")
			fmt.Println("SMS from " + number + ":" + sms)
			port.Write([]byte("AT+CMGD=" + idx + "\r\n"))
			//TODO: now as we have the message, we have to call out to the world.
		}
	}
}

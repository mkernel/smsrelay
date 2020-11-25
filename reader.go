package main

import (
	"io"
	"strings"
)

var check_return bool = false
var collect_reply bool = false
var incoming_sms chan string = make(chan string, 1)
var proceed chan string = make(chan string, 1)
var collect chan []string = make(chan []string, 1)

func reader(reader io.Reader) {
	//buffered := bufio.NewReader(reader)
	var buf []byte
	var collect_buf []string
	for true {
		var input []byte = make([]byte, 1)
		reader.Read(input)
		buf = append(buf, input[0])
		casted := string(buf)
		if strings.HasSuffix(casted, "\r\n") {
			casted = casted[:len(casted)-2]
			//fmt.Println(casted)
			if collect_reply {
				collect_buf = append(collect_buf, casted)
				if casted == "OK" {
					//this was the last line to collect.
					collect <- collect_buf
					collect_buf = nil
					collect_reply = false
				}
			}
			if check_return && casted == "OK" {
				proceed <- casted
			}
			if check_return && casted == "ERROR" {
				proceed <- casted
			}
			if strings.HasPrefix(casted, "+CMTI: \"SM\",") {
				//let's get the index of the string.
				components := strings.Split(casted, ",")
				idx := components[1]
				incoming_sms <- idx
			}
			//TODO: we have to send out data from here. But we will go with that one by one...
			buf = nil
		}
	}
}

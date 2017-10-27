package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var addr string
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addr = ipnet.IP.String()

			}

		}
	}

	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		host, _ = os.Hostname()
	}

	io.WriteString(w, "ip is "+addr+"\n"+"hostname is "+host+"\n")
}

func main() {
	http.HandleFunc("/", helloHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

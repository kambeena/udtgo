package main

import (
	"fmt"
	"github.com/kambeena/udtgo"
	"encoding/json"
	"os"
	"strconv"
)

func main()  {

	portno := 9000
	network := "ip4"
	isStream := true
	host := "localhost"
	fileName := "myfile"

	if len(os.Args) <= 1 {
		fmt.Println("Please provide filname to be uploaded.usage : filename port server")
		return
	}


	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	}

	if len(os.Args) >= 3 {
		portno, _ = strconv.Atoi(os.Args[2])
	}

	if len(os.Args) >= 3 {
		host = os.Args[3]
	}

	fmt.Printf("staring client to %s %d \n", host, portno)

	s, err := startClient(network, host, portno, isStream)

	if err != nil {
		fmt.Errorf("Unable to start client")
	}
	defer udtgo.Close(s)

	fi, err := os.Lstat(fileName)
	if err != nil {
		fmt.Errorf("Unable read file %s", fileName)
		return
	}

	reqContent := map[string]interface{}{
		"fileName": fileName+"_copy",
		"fileSize" : fi.Size(),
	}

	msg, err := json.Marshal(reqContent)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}

	n, err := udtgo.Send(s, &msg[0], len(msg))

	if err != nil {
		fmt.Errorf("Unable to send request %s %d", err, n)
	}

	fmt.Printf("Request sent %s \n", string(msg))

	//send file
	var offset int64 = 0

	datasent, err := udtgo.Sendfile(s, fileName, &offset, fi.Size())

	if err != nil {
		fmt.Errorf("Unable to send file %s %d", err, datasent)
	}


}

func startClient(network string, host string, portno int, isStream bool) (socket *udtgo.Socket, err error) {

	socket, err = udtgo.CreateSocket(network, isStream)
	if err != nil {
		return nil, fmt.Errorf("Unable to create socket :%s", err)
	}

	n, err := udtgo.Connect(socket, host, portno)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to the socket :%s %d", err, n)
	}

	return
}

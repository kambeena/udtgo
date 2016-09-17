package main

import (
	"github.com/kambeena/udtgo"
	"fmt"
	"strings"
	"encoding/json"
	"os"
	"log"
	"strconv"
	"path"
)



func main() {

	//start UDT system
	udtgo.Startup()

	portno := 9000
	network := "ip4"
	isStream := true
	uploadDir := ""

	if len(os.Args) >= 2 {
		portno, _ = strconv.Atoi(os.Args[1])
	}

	if len(os.Args) >= 3 {
		uploadDir = os.Args[2]
	}

	s, err := startServer(portno, network, isStream)
	if err != nil {
		fmt.Printf("Unable to start server at %d %s", portno)
		log.Fatal(err)
	}

	defer udtgo.Close(s)

	fmt.Printf("Started Server at port %d \n", portno)


	for {
		ns, err := udtgo.Accept(s)
		if err != nil {
			fmt.Printf("Unable to accept request on socket")
			continue
		}

		go handleRequest(ns, uploadDir)
	}

}

func handleRequest(socket *udtgo.Socket, uploadDir string) {

	defer udtgo.Close(socket)
	//receive message
	request, err := receiveRequest(socket)

	if err != nil {
		fmt.Printf("Unable to get request %s", err)
		return
	}

	//ummarshall request
	reqObject := make(map[string]interface{})
	err = json.Unmarshal([]byte(request), &reqObject)

	if err != nil {
		fmt.Printf("Unable to Unmarshal request %s", err)
		return
	}

	fileName := reqObject["fileName"].(string)

	if uploadDir != "" {
		fileName = fmt.Sprintf("%s%s", uploadDir, fileName)
	}

	base := path.Dir(fileName)

	err = os.MkdirAll(base, os.ModePerm)

	if err != nil {
		fmt.Printf("Unable to create dir path %s", err)
		return
	}

	fileSize := int64(reqObject["fileSize"].(float64))
	fmt.Printf("Received request fileName %s fileSize %v \n", fileName, fileSize)

	var offset int64 = 0

	n, err := udtgo.Recvfile(socket, fileName, &offset, fileSize)

	if err != nil {
		fmt.Printf("Unable to receive file %s %d", err, n)
		return
	}

	fi, err := os.Lstat(fileName)
	if err != nil {
		fmt.Printf("Unable read file %s", fileName)
		return
	}

	fmt.Printf("Successfully recived file %s \n", fi.Name())
}

func receiveRequest(socket *udtgo.Socket) (request string, err error){

	data := make([]byte, 10000)

	n, err := udtgo.Recv(socket, &data[0], len(data))

	if err != nil {
		return "", fmt.Errorf("Unable to receive data %s %d", err, n)
	}
	request = strings.TrimSpace(string(data[:n]))

	return request, nil

}


func startServer(portno int, network string, isStream bool) (socket *udtgo.Socket, err error) {

	socket, err = udtgo.CreateSocket(network, isStream)
	if err != nil {
		return nil, fmt.Errorf("Unable to create socket :%s", err)
	}
	n, err := udtgo.Bind(socket, portno)
	if err != nil {
		return nil, fmt.Errorf("Unable to bind socket :%d", n)
	}
	n, err = udtgo.Listen(socket, 4)
	if err != nil {
		return nil, fmt.Errorf("Unable to listen socket :%d", n)
	}

	return
}

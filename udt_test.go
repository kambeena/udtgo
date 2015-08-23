/*****************************************************************************
Copyright (c) 2015, Kamlesh Sharma at kambeena@gmail.com.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above
  copyright notice, this list of conditions and the
  following disclaimer.

* Redistributions in binary form must reproduce the
  above copyright notice, this list of conditions
  and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the creator nor the names of its contributors may be used to
  endorse or promote products derived from this
  software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*****************************************************************************/

/*****************************************************************************
written by
   Kamlesh Sharma, last updated 06/21/2015
*****************************************************************************/

package udtgo

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	PORT9000 = 9000 + iota
	PORT9001
	PORT9002
	PORT9003
	PORT9004
	PORT9005
	PORT9006
	PORT9007
	PORT9008
)

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	shutdown()
	os.Exit(exitCode)
}

func setup() {
	Startup()
}

func shutdown() {
	Cleanup()
}

func TestSendRecvData(t *testing.T) {
	portno := PORT9000
	network := "ip4"
	isStream := true
	s, err := startServer(portno, network, isStream)
	if err != nil {
		t.Errorf("Unable to start server")
	}
	defer Close(s)

	message := "Hello from Kamlesh"
	go sendData(t, network, "localhost", portno, isStream, message)

	ns, err := Accept(s)
	if err != nil {
		t.Errorf("Unable to accept request on socket")
	}
	defer Close(ns)

	data := make([]byte, 100)

	n, err := Recv(ns, &data[0], len(data))

	if err != nil {
		t.Errorf("Unable to receive data %s %d", err, n)
	}
	recvMsg := strings.TrimSpace(string(data[:n]))

	if !(bytes.Equal(([]byte)(message), ([]byte)(recvMsg))) {
		t.Errorf("Unable to verify the message")
	}

	return
}

func TestSendRecvFile(t *testing.T) {
	portno := PORT9001
	network := "ip4"
	isStream := true
	s, err := startServer(portno, network, isStream)
	if err != nil {
		t.Errorf("Unable to start server")
	}
	defer Close(s)

	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Unable to get working director")
	}
	fmt.Println(pwd)

	filepath := strings.TrimSpace(pwd + string(os.PathSeparator) + "test1.jar")
	destFilepath := strings.TrimSpace(pwd + string(os.PathSeparator) + "test1_copy.jar")
	fi, err := os.Lstat(filepath)
	if err != nil {
		t.Errorf("Unable read file %s", filepath)
	}

	size := fi.Size()

	go sendFile(t, network, "localhost", portno, isStream, filepath, size)

	ns, err := Accept(s)
	if err != nil {
		t.Errorf("Unable to accept request on socket")
	}
	defer Close(ns)

	//data := ([]byte)(destFilepath)
	var offset int64 = 0

	n, err := Recvfile(ns, destFilepath, &offset, size)

	if err != nil {
		t.Errorf("Unable to receive file %s %d", err, n)
	}

	fi, err = os.Lstat(destFilepath)
	if err != nil {
		t.Errorf("Unable read file %s", destFilepath)
	}

	return
}

func TestSendRecvMsgServer(t *testing.T) {
	portno := PORT9002
	network := "ip4"
	isStream := false
	s, err := startServer(portno, network, isStream)
	if err != nil {
		t.Errorf("Unable to start server")
	}
	defer Close(s)

	message := "Hello from Kamlesh"

	go sendMsg(t, network, "localhost", portno, isStream, message)

	ns, err := Accept(s)
	if err != nil {
		t.Errorf("Unable to accept request on socket")
	}
	defer Close(ns)

	data := make([]byte, 100)

	n, err := RecvMsg(ns, &data[0], len(data))

	if err != nil {
		t.Errorf("Unable to receive data %s %d", err, n)
	}
	recvMsg := strings.TrimSpace(string(data[:n]))

	if !(bytes.Equal(([]byte)(message), ([]byte)(recvMsg))) {
		t.Errorf("Unable to verify the message")
	}

	return
}

func TestGetsockstate(t *testing.T) {
	portno := PORT9003
	network := "ip4"
	isStream := true

	socket, err := CreateSocket(network, isStream)
	if err != nil {
		t.Errorf("Unable to create socket :%s", err)
	}

	socksstate, err := Getsockstate(socket)
	if socksstate != INIT {
		t.Errorf("Socket status should be %d got :%d", INIT, socksstate)
	}

	n, err := Bind(socket, portno)
	if err != nil {
		t.Errorf("Unable to bind socket :%d", n)
	}

	socksstate, err = Getsockstate(socket)
	if socksstate != OPENED {
		t.Errorf("Socket status should be %d got :%d", OPENED, socksstate)
	}

	n, err = Listen(socket, 4)
	if err != nil {
		t.Errorf("Unable to listen socket :%d", n)
	}

	socksstate, err = Getsockstate(socket)
	if socksstate != LISTENING {
		t.Errorf("Socket status should be %d got :%d", LISTENING, socksstate)
	}

	n, err = Close(socket)
	if err != nil {
		t.Errorf("Unable to close socket :%d", n)
	}

	socksstate, err = Getsockstate(socket)
	if socksstate != BROKEN {
		t.Errorf("Socket status should be %d got :%d", BROKEN, socksstate)
	}

}

func TestGetsockopt(t *testing.T) {
	portno := PORT9004
	network := "ip4"
	isStream := true

	socket, err := CreateSocket(network, isStream)
	if err != nil {
		t.Errorf("Unable to create socket :%s", err)
	}

	n, err := Bind(socket, portno)
	if err != nil {
		t.Errorf("Unable to bind socket :%d", n)
	}

	n, err = Listen(socket, 4)
	if err != nil {
		t.Errorf("Unable to listen socket :%d", n)
	}

	socksopt, err := Getsockopt(socket, UDT_SNDSYN)
	if socksopt.(uint64) != 1 {
		t.Errorf("Return value should be 1 but returned :%d", socksopt)
	}
	socksopt, err = Getsockopt(socket, UDT_MSS)
	if socksopt.(uint16) != 1500 {
		t.Errorf("Return value should be 1500 but returned :%d", socksopt)
	}
	socksopt, err = Getsockopt(socket, UDT_MAXBW)
	if socksopt.(uint64) <= 0 {
		t.Errorf("Return value should be more than 0 but returned :%d", socksopt)
	}
	socksopt, err = Getsockopt(socket, UDT_FC)
	if socksopt.(uint16) != 25600 {
		t.Errorf("Return value should be 25600 but returned :%d", socksopt)
	}
	socksopt, err = Getsockopt(socket, UDT_RENDEZVOUS)
	if socksopt.(uint64) != 0 {
		t.Errorf("Return value should be 1 but returned :%d", socksopt)
	}
	socksopt, err = Getsockopt(socket, UDT_SNDBUF)
	t.Logf("Socket option %d", socksopt)
	socksopt, err = Getsockopt(socket, UDT_RCVBUF)
	t.Logf("Socket option %d", socksopt)
	socksopt, err = Getsockopt(socket, UDT_REUSEADDR)
	if socksopt.(uint64) != 1 {
		t.Errorf("Return value should be 1 but returned :%d", socksopt)
	}
	socksopt, err = Getsockopt(socket, UDT_SNDTIMEO)
	t.Logf("Socket option %d", socksopt)

}

func TestSetsockopt(t *testing.T) {
	network := "ip4"
	isStream := true
	var MSS_VAL uint16 = 500

	socket, err := CreateSocket(network, isStream)
	if err != nil {
		t.Errorf("Unable to create socket :%s", err)
	}

	var optval uint16 = MSS_VAL //don't change this value

	n, err := Setsockopt(socket, UDT_MSS, optval)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err := Getsockopt(socket, UDT_MSS)
	if optvalupdated.(uint16) != optval {
		t.Errorf("Return value should be %d but returned :%d", optval, optvalupdated)
	}

	var optval1 uint64 = 0

	n, err = Setsockopt(socket, UDT_SNDSYN, optval1)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_SNDSYN)
	if optvalupdated.(uint64) != optval1 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optval1 = 0

	n, err = Setsockopt(socket, UDT_RCVSYN, optval1)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_RCVSYN)
	if optvalupdated.(uint64) != optval1 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optval = 1500

	n, err = Setsockopt(socket, UDT_FC, optval)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_FC)
	if optvalupdated.(uint16) != optval {
		t.Errorf("Return value should be %d but returned :%d", optval, optvalupdated)
	}

	optval = 20000
	n, err = Setsockopt(socket, UDT_SNDBUF, optval)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_SNDBUF)
	var sndbufVal uint16 = 20000 / (MSS_VAL - 28)
	if optvalupdated.(uint16) != sndbufVal*(MSS_VAL-28) {
		t.Errorf("Return value should be %d but returned :%d", sndbufVal*(MSS_VAL-28), optvalupdated)
	}

	optval = 20
	n, err = Setsockopt(socket, UDT_RCVBUF, optval)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_RCVBUF)
	//value should be 32*(MSS_VAL-28)
	if optvalupdated.(uint16) != 32*(MSS_VAL-28) {
		t.Errorf("Return value should be %d but returned :%d", optval, optvalupdated)
	}

	optval = 1500

	n, err = Setsockopt(socket, UDP_SNDBUF, optval)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDP_SNDBUF)
	if optvalupdated.(uint16) != optval {
		t.Errorf("Return value should be %d but returned :%d", optval, optvalupdated)
	}

	optval = 1500

	n, err = Setsockopt(socket, UDP_RCVBUF, optval)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDP_RCVBUF)
	if optvalupdated.(uint16) != optval {
		t.Errorf("Return value should be %d but returned :%d", optval, optvalupdated)
	}

	optval1 = 1

	n, err = Setsockopt(socket, UDT_RENDEZVOUS, optval1)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_RENDEZVOUS)
	if optvalupdated.(uint64) != optval1 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optval1 = 0

	n, err = Setsockopt(socket, UDT_REUSEADDR, optval1)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_REUSEADDR)
	if optvalupdated.(uint64) != optval1 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optval1 = 20971520

	n, err = Setsockopt(socket, UDT_MAXBW, optval1)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}
	optvalupdated, err = Getsockopt(socket, UDT_MAXBW)
	if optvalupdated.(uint64) != optval1 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	linger := Linger{
		l_onoff:  1,
		l_linger: 250,
	}

	n, err = Setsockopt(socket, UDT_LINGER, linger)
	if err != nil {
		t.Errorf("Unable to set option :%s %d", err, n)
	}

	optvalupdated, err = Getsockopt(socket, UDT_LINGER)
	if err != nil {
		t.Logf("err %d", err)
	}
	resval := optvalupdated.(Linger)
	if resval.l_onoff != linger.l_onoff {
		t.Errorf("Return value should be %d but returned :%d", linger.l_onoff, resval.l_onoff)
	}
	if resval.l_linger != linger.l_linger {
		t.Errorf("Return value should be %d but returned :%d", linger.l_linger, resval.l_linger)
	}

	optvalupdated, err = Getsockopt(socket, UDT_STATE)
	if err != nil {
		t.Logf("err %d", err)
	}
	if optvalupdated.(uint16) != 1 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optvalupdated, err = Getsockopt(socket, UDT_EVENT)
	if err != nil {
		t.Logf("err %d", err)
	}
	if optvalupdated.(uint16) != 0 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optvalupdated, err = Getsockopt(socket, UDT_SNDDATA)
	if err != nil {
		t.Logf("err %d", err)
	}
	if optvalupdated.(uint16) != 0 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

	optvalupdated, err = Getsockopt(socket, UDT_RCVDATA)
	if err != nil {
		t.Logf("err %d", err)
	}
	if optvalupdated.(uint16) != 0 {
		t.Errorf("Return value should be %d but returned :%d", optval1, optvalupdated)
	}

}

func TestSocknames(t *testing.T) {
	s, err := startServer(PORT9005, "ip4", true)
	if err != nil {
		t.Errorf("Unable to startserver")
	}
	defer Close(s)

	sockaddr, err := Getsockname(s)
	if err != nil {
		t.Errorf("Unable to get sock name %s", err)
	}
	if sockaddr.sa_data != "0.0.0.0" {
		t.Errorf("Unable to get sock name")
	}

	sc, err := startClient("ip4", "localhost", PORT9005, true)

	if err != nil {
		t.Errorf("Unable to start client")
	}
	defer Close(sc)

	sockaddr, err = Getpeername(sc)
	if err != nil {
		t.Errorf("Unable to get sock peer name %s", err)
	}
	if sockaddr.sa_data != "127.0.0.1" {
		t.Errorf("Unable to get peer sock name")
	}

}

func TestPerfmon(t *testing.T) {

	s, err := startServer(PORT9006, "ip4", true)
	if err != nil {
		t.Errorf("Unable to start server")
	}
	defer Close(s)

	dataSize := 100000

	go sendPerfmonData(t, "ip4", "localhost", PORT9006, true, dataSize)

	ns, err := Accept(s)
	if err != nil {
		t.Errorf("Unable to accept request on socket")
	}
	defer Close(ns)

	data := make([]byte, 10)
	rSize := 10

	for rSize < dataSize {
		n, err := Recv(ns, &data[0], 10)
		if err != nil {
			t.Errorf("Unable to receive data %s %d", err, n)
		}
		rSize = rSize + n

	}

}

func TestEpollClientEvent(t *testing.T) {
	s, err := startServer(PORT9008, "ip4", true)
	if err != nil {
		t.Errorf("Unable to start server")
	}
	defer Close(s)

	message := "Hello from Kamlesh"
	go sendDataUsingEpoll(t, "ip4", "localhost", PORT9008, true, message)

	ns, err := Accept(s)
	if err != nil {
		t.Errorf("Unable to accept request on socket")
	}
	defer Close(ns)

	data := make([]byte, 100)

	n, err := Recv(ns, &data[0], len(data))

	if err != nil {
		t.Errorf("Unable to receive data %s %d", err, n)
	}
	recvMsg := strings.TrimSpace(string(data[:n]))

	if !(bytes.Equal(([]byte)(message), ([]byte)(recvMsg))) {
		t.Errorf("Unable to verify the message")
	}
}

func TestEpollServerEvent(t *testing.T) {
	s, err := startServer(PORT9007, "ip4", true)
	if err != nil {
		t.Errorf("Unable to start server")
	}
	defer Close(s)

	message := "Hello from Kamlesh"
	go sendData(t, "ip4", "localhost", PORT9007, true, message)

	ns, err := Accept(s)
	if err != nil {
		t.Errorf("Unable to accept request on socket")
	}
	defer Close(ns)

	eid, err := EpollCreate()
	if err != nil {
		t.Errorf("Unable to create Epoll id")
	}

	n, err := EpollAddUsock(eid, ns, UDT_EPOLL_IN|UDT_EPOLL_OUT)
	if err != nil {
		t.Errorf("Unable to add socket for polling")
	}

	ursocks := CreateSockets(1)
	uwsocks := CreateSockets(1)
	sysrsocks := CreateSysSockets(1)
	syswsocks := CreateSysSockets(1)

	num := 0
	num, err = EpollWait2(eid, ursocks.socks, uwsocks.socks, int64(-1), sysrsocks.socks, syswsocks.socks)
	if err != nil {
		t.Errorf("Unable to wait on polling")
	}

	t.Logf("# sockets %d", num)
	t.Logf("# retval last %d\n", n)

	socket := &Socket{
		sock: uwsocks.socks[0],
	}

	socksstate, err := Getsockstate(socket)
	if err != nil {
		t.Errorf("Unable to get socket state  %s %d", err, socksstate)
	}
	t.Logf("Socket state %d\n", socksstate)

	data := make([]byte, 100)
	n, err = Recv(socket, &data[0], len(data))
	if err != nil {
		t.Errorf("Unable to receive data %s %d", err, n)
	}
	recvMsg := strings.TrimSpace(string(data[:n]))
	if !(bytes.Equal(([]byte)(message), ([]byte)(recvMsg))) {
		t.Errorf("Unable to verify the message")
	}
	t.Logf("message received from epoll = %s", recvMsg)

}

func startServer(portno int, network string, isStream bool) (socket *Socket, err error) {

	socket, err = CreateSocket(network, isStream)
	if err != nil {
		return nil, fmt.Errorf("Unable to create socket :%s", err)
	}
	n, err := Bind(socket, portno)
	if err != nil {
		return nil, fmt.Errorf("Unable to bind socket :%d", n)
	}
	n, err = Listen(socket, 4)
	if err != nil {
		return nil, fmt.Errorf("Unable to listen socket :%d", n)
	}

	return
}

func sendPerfmonData(t *testing.T, network string, host string, portno int,
	isStream bool, dataSize int) {

	s, err := startClient(network, host, portno, isStream)

	if err != nil {
		t.Errorf("Unable to start client")
	}
	defer Close(s)

	data := make([]byte, dataSize)
	sDataSize := 10

	go monitor(t, s)

	for sDataSize < dataSize {
		n, err := Send(s, &data[0], sDataSize)

		if err != nil {
			t.Errorf("Unable to send data %s %d", err, n)
		}

		sDataSize = sDataSize + n

	}

}

func monitor(t *testing.T, socket *Socket) {

	t.Logf("SendRate(Mb/s)\tRTT(ms)\tCWnd\tPktSndPeriod(us)\tRecvACK\tRecvNAK")

	for {
		time.Sleep(10 * time.Millisecond)
		traceinfo, err := Perfmon(socket, false)
		if err != nil {
			t.Logf("Unable to get perform data %s", err)
			break
		}
		t.Logf("%f\t%f\t%d\t%f\t%d\t%d",
			traceinfo.mbpsSendRate, traceinfo.msRTT, traceinfo.pktCongestionWindow,
			traceinfo.usPktSndPeriod, traceinfo.pktRecvACK, traceinfo.pktRecvNAK)
	}
}

func sendDataUsingEpoll(t *testing.T, network string, host string, portno int, isStream bool, message string) {

	s, err := startClient(network, host, portno, isStream)

	if err != nil {
		t.Errorf("Unable to start client")
	}

	eid, err := EpollCreate()
	if err != nil {
		t.Errorf("Unable to create Epoll id")
	}

	n, err := EpollAddUsock(eid, s, UDT_EPOLL_OUT)
	if err != nil {
		t.Errorf("Unable to add socket for polling")
	}

	ursocks := CreateSockets(1)
	uwsocks := CreateSockets(1)
	sysrsocks := CreateSysSockets(1)
	syswsocks := CreateSysSockets(1)

	n, err = EpollWait2(eid, ursocks.socks, uwsocks.socks, int64(-1), sysrsocks.socks, syswsocks.socks)
	if err != nil {
		t.Errorf("Unable to wait on polling")
	}

	socket := &Socket{
		sock: uwsocks.socks[0],
	}

	socksstate, err := Getsockstate(socket)
	if err != nil {
		t.Errorf("Unable to get socket state  %s %d", err, socksstate)
	}
	t.Logf("Socket state %d\n", socksstate)

	defer Close(s)
	msg := ([]byte)(message)

	n, err = Send(socket, &msg[0], len(msg))

	if err != nil {
		t.Errorf("Unable to send data %s %d", err, n)
	}

	return
}

func sendData(t *testing.T, network string, host string, portno int, isStream bool, message string) {

	s, err := startClient(network, host, portno, isStream)

	if err != nil {
		t.Errorf("Unable to start client")
	}
	defer Close(s)
	msg := ([]byte)(message)

	n, err := Send(s, &msg[0], len(msg))

	if err != nil {
		t.Errorf("Unable to send data %s %d", err, n)
	}

	return
}

func sendFile(t *testing.T, network string, host string, portno int, isStream bool,
	filepath string, size int64) {

	s, err := startClient(network, host, portno, isStream)

	if err != nil {
		t.Errorf("Unable to start client")
	}
	defer Close(s)

	//data := ([]byte)(filepath)
	var offset int64 = 0

	n, err := Sendfile(s, filepath, &offset, size)

	if err != nil {
		t.Errorf("Unable to send file %s %d", err, n)
	}

	return
}

func sendMsg(t *testing.T, network string, host string, portno int, isStream bool, message string) {

	s, err := startClient(network, host, portno, isStream)

	if err != nil {
		t.Errorf("Unable to start client")
	}
	defer Close(s)
	msg := ([]byte)(message)

	n, err := SendMsg(s, &msg[0], len(msg), -1, false)

	if err != nil {
		t.Errorf("Unable to send data %s %d", err, n)
	}

	return
}

func startClient(network string, host string, portno int, isStream bool) (socket *Socket, err error) {

	socket, err = CreateSocket(network, isStream)
	if err != nil {
		return nil, fmt.Errorf("Unable to create socket :%s", err)
	}

	n, err := Connect(socket, host, portno)

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to the socket :%s %d", err, n)
	}

	return
}

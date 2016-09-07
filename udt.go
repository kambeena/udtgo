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
   Kamlesh Sharma, last updated 09/04/2016
*****************************************************************************/

package udtgo

// #cgo LDFLAGS: /usr/local/lib/libudt.so
//
// #include "udtc.h"
// #include <arpa/inet.h>
// #include <string.h>
// #include <stdlib.h>
// #include <sys/socket.h>
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"unsafe"
)

type Socket struct {
	sock C.UDTSOCKET
	af   C.int
}

type Sockets struct {
	socks []C.UDTSOCKET
}

type SysSockets struct {
	socks []C.SYSSOCKET
}

type Sockaddr struct {
	sa_family int
	sa_data   string
}

type Linger struct {
	l_onoff  int
	l_linger int
}

type Traceinfo struct {
	msTimeStamp        int64 // time since the UDT entity is started, in milliseconds
	pktSentTotal       int64 // total number of sent data packets, including retransmissions
	pktRecvTotal       int64 // total number of received packets
	pktSndLossTotal    int   // total number of lost packets (sender side)
	pktRcvLossTotal    int   // total number of lost packets (receiver side)
	pktRetransTotal    int   // total number of retransmitted packets
	pktSentACKTotal    int   // total number of sent ACK packets
	pktRecvACKTotal    int   // total number of received ACK packets
	pktSentNAKTotal    int   // total number of sent NAK packets
	pktRecvNAKTotal    int   // total number of received NAK packets
	usSndDurationTotal int64 // total time duration when UDT is sending data (idle time exclusive)

	// local measurements
	pktSent       int64   // number of sent data packets, including retransmissions
	pktRecv       int64   // number of received packets
	pktSndLoss    int     // number of lost packets (sender side)
	pktRcvLoss    int     // number of lost packets (receiver side)
	pktRetrans    int     // number of retransmitted packets
	pktSentACK    int     // number of sent ACK packets
	pktRecvACK    int     // number of received ACK packets
	pktSentNAK    int     // number of sent NAK packets
	pktRecvNAK    int     // number of received NAK packets
	mbpsSendRate  float64 // sending rate in Mb/s
	mbpsRecvRate  float64 // receiving rate in Mb/s
	usSndDuration int64   // busy sending time (i.e., idle time exclusive)

	// instant measurements
	usPktSndPeriod      float64 // packet sending period, in microseconds
	pktFlowWindow       int     // flow window size, in number of packets
	pktCongestionWindow int     // congestion window size, in number of packets
	pktFlightSize       int     // number of packets on flight
	msRTT               float64 // RTT, in milliseconds
	mbpsBandwidth       float64 // estimated bandwidth, in Mb/s
	byteAvailSndBuf     int     // available UDT sender buffer size
	byteAvailRcvBuf     int     // available UDT receiver buffer size
}

const (
	INVLID int = iota
	INIT
	OPENED
	LISTENING
	CONNECTING
	CONNECTED
	BROKEN
	CLOSING
	CLOSED
	NONEXIST
)

const (
	UDT_EPOLL_IN  = 0x1
	UDT_EPOLL_OUT = 0x4
	UDT_EPOLL_ERR = 0x8
)

const (
	UDT_MSS        string = "UDT_MSS"
	UDT_SNDSYN     string = "UDT_SNDSYN"
	UDT_RCVSYN     string = "UDT_RCVSYN"
	UDT_CC         string = "UDT_CC"
	UDT_FC         string = "UDT_FC"
	UDT_SNDBUF     string = "UDT_SNDBUF"
	UDT_RCVBUF     string = "UDT_RCVBUF"
	UDP_SNDBUF     string = "UDP_SNDBUF"
	UDP_RCVBUF     string = "UDP_RCVBUF"
	UDT_LINGER     string = "UDT_LINGER"
	UDT_RENDEZVOUS string = "UDT_RENDEZVOUS"
	UDT_SNDTIMEO   string = "UDT_SNDTIMEO"
	UDT_RCVTIMEO   string = "UDT_RCVTIMEO"
	UDT_REUSEADDR  string = "UDT_REUSEADDR"
	UDT_MAXBW      string = "UDT_MAXBW"
	UDT_STATE      string = "UDT_STATE"
	UDT_EVENT      string = "UDT_EVENT"
	UDT_SNDDATA    string = "UDT_SNDDATA"
	UDT_RCVDATA    string = "UDT_RCVDATA"
)

//Use this function to create udt socket. This function returns
//structure contains udt socket and information about IP
//family AF_INET or AF_INET6. 
//parameter - network - IP family ip4 or ip6
//parameter - isStream - true socket type SOCK_STREAM or SOCK_DGRAM


func CreateSocket(network string, isStream bool) (socket *Socket, err error) {
	var n C.int

	if network == "ip4" {
		n = C.AF_INET
	} else if network == "ip6" {
		n = C.AF_INET6
	} else {
		return nil, fmt.Errorf("network must be either ip4 or ip6")
	}

	var trnType C.int
	if isStream {
		trnType = C.SOCK_STREAM
	} else {
		trnType = C.SOCK_DGRAM
	}

	sock := C.udt_socket(n, trnType, 0)

	if  C.int(sock) == C.int(C.UDT_INVALID_SOCK) {
		return nil, udtErrDesc("Invalid socket")
	}

	socket = &Socket{
		sock: sock,
		af:   n,
	}

	return
}

//Binds socket to the passed port number. If the binding is successful, bind returns 0, otherwise it returns
//error code (http://udt.sourceforge.net/udt4/doc/ecode.htm) and error object with error details.

func Bind(socket *Socket, portno int) (retval int, err error) {

	var serv_addr C.struct_sockaddr_in
	serv_addr.sin_family = C.sa_family_t(socket.af)
	serv_addr.sin_port = C.in_port_t(C.htons(C.uint16_t(portno)))
	serv_addr.sin_addr.s_addr = C.INADDR_ANY
	if _, err := C.memset(unsafe.Pointer(&(serv_addr.sin_zero)), 0, 8); err != nil {
		return -1, fmt.Errorf("Unable to zero sin_zero")
	}

	retval = int(C.udt_bind(socket.sock, (*C.struct_sockaddr)(unsafe.Pointer(&serv_addr)),
		C.int(unsafe.Sizeof(serv_addr))))
	if retval < 0 {
		return -1, udtErrDesc("Unable to bind socket")
	}
	return 
}

//This function turns socket to listening state and makes socket ready to recieve connection
//requests. Pass backlog parameter to configure number of pending connections. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func Listen(socket *Socket, backlog int) (retval int, err error) {

	retval = int(C.udt_listen(socket.sock, C.int(backlog)))

	if retval < 0 {
		return retval, udtErrDesc("Unable to listen")
	}
	return
}

//Retrieves and returns newly accepted socket. If successful,
// this method returns new socket and error object if unable to accept new socket with error details.

func Accept(socket *Socket) (newSocket *Socket, err error) {
	var cli_addr C.struct_sockaddr_in
	var addrlen C.int
	newSock := C.udt_accept(socket.sock, (*C.struct_sockaddr)(unsafe.Pointer(&cli_addr)),
		&addrlen)

	if  C.int(newSock) == C.int(C.UDT_INVALID_SOCK) {
		return nil, udtErrDesc("Unable to accept on socket")
	}

	newSocket = &Socket{
		sock: newSock,
	}

	return
}

//The connect method connects to a server socket (in regular mode) or
// a peer socket (in rendezvous mode) to set up a UDT connection. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.


func Connect(socket *Socket, host string, portno int) (retval int, err error) {

	addrs, err := net.LookupHost(host)

	if err != nil {
		return -1, fmt.Errorf("Unable to connect to the socket: %s", err)
	}

	resolvedHost := addrs[0]

	var serv_addr C.struct_sockaddr_in
	serv_addr.sin_family = C.sa_family_t(socket.af)
	serv_addr.sin_port = C.in_port_t(C.htons(C.uint16_t(portno)))
	C.inet_pton(socket.af, C.CString(resolvedHost), unsafe.Pointer(&serv_addr.sin_addr))

	retval = int(C.udt_connect(socket.sock, (*C.struct_sockaddr)(unsafe.Pointer(&serv_addr)),
		C.int(unsafe.Sizeof(serv_addr))))
	if retval < 0 {
		return retval, udtErrDesc("Unable to connect to the socket")
	}
	return
}

//Retrives socket status. If successful, this method returns status, otherwise it
// returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.


func Getsockstate(socket *Socket) (status int, err error) {
	status = int(C.udt_getsockstate(socket.sock))
	if status < 0 {
		return status, udtErrDesc("Unable to get socket status")
	}
	return
}

//Sends out certain amount of data from an application buffer. If successful, this method returns size of the data send, otherwise it
//returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
//and error object with error details.

func Send(socket *Socket, data *byte, length int) (retval int, err error) {

	retval = int(C.udt_send(socket.sock, (*C.char)(unsafe.Pointer(data)), C.int(length), C.int(0)))
	if retval < 0 {
		return retval, udtErrDesc("Unable to send data")
	}
	return
}

//This method reads certain amount of data into a local memory buffer. If successful, this method returns size of the data received otherwise it
//returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
//and error object with error details.

func Recv(socket *Socket, data *byte, length int) (retval int, err error) {

	retval = int(C.udt_recv(socket.sock, (*C.char)(unsafe.Pointer(data)), C.int(length), C.int(0)))
	if retval < 0 {
		return retval, udtErrDesc("Unable to recive data")
	}
	return
}

//The sendmsg method sends a message to the peer side. The input parameters contains socket, data, data length, message
//ttl (optional) (time to live) in milliseconds (default ttl is -1, which means infinite) and flag (optional) indicating if the message
// should be delivered in order (default is negative). If successful, this method returns size of the message sent otherwise it
//returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
//and error object with error details.

func SendMsg(socket *Socket, data *byte, length int, ttl int, inorder bool) (retval int, err error) {
	var cInorder C.int
	if inorder {
		cInorder = 1
	} else {
		cInorder = 0
	}
	retval = int(C.udt_sendmsg(socket.sock, (*C.char)(unsafe.Pointer(data)),
		C.int(length), C.int(ttl), cInorder))
	if retval < 0 {
		return retval, udtErrDesc("Unable to send message")
	}
	return
}

//The recvmsg method receives a valid message. If successful, this method returns size of the message sent otherwise it
//returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
//and error object with error details.

func RecvMsg(socket *Socket, data *byte, length int) (retval int, err error) {

	retval = int(C.udt_recvmsg(socket.sock, (*C.char)(unsafe.Pointer(data)), C.int(length)))
	if retval < 0 {
		return retval, udtErrDesc("Unable to receive message")
	}
	return
}

//This method send local file. On success, sendfile returns the actual size of data that has been sent
//otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)

func Sendfile(socket *Socket, filepath string, offset *int64, size int64) (retval int64, err error) {
	
	return Sendfile2(socket, filepath, offset, size, 7320000)
}

//This method send local file using defined block size as input parameter. On success, sendfile returns the actual size of data that has been sent
//otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
//and error object with error details.

func Sendfile2(socket *Socket, filepath string, offset *int64, 
		size int64, block int) (retval int64, err error) {

	retval = int64(C.udt_sendfile2(socket.sock, C.CString(filepath),
		(*C.int64_t)(unsafe.Pointer(offset)), C.int64_t(size), C.int(block)))
	if retval < 0 {
		return retval, udtErrDesc("Unable to send file ")
	}
	return
}

//The recvfile method reads certain amount of data into a local file. This method usages block size of 366000.
//On success, recvfile returns the actual size of received data otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)

func Recvfile(socket *Socket, filepath string, offset *int64, size int64) (retval int64, err error) {

	return Recvfile2(socket, filepath, offset, size, 366000)
}

//The recvfile method reads certain amount of data into a local file. This method usages block size provided as input parameter.
//On success, recvfile returns the actual size of received data otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
//and error object with error details.

func Recvfile2(socket *Socket, filepath string, offset *int64, 
					size int64, block int) (retval int64, err error) {

	retval = int64(C.udt_recvfile2(socket.sock, C.CString(filepath),
		(*C.int64_t)(unsafe.Pointer(offset)), C.int64_t(size), C.int(block)))
	if retval < 0 {
		return retval, udtErrDesc("Unable to receive file")
	}
	return
}

//The method reads UDT socket options. If successful, returns requested option value otherwise
//returns error object with error details.

func Getsockopt(socket *Socket, option string) (value interface{}, err error) {

	var data []byte
	var optlen C.int
	var retval int = 0
	data = make([]byte, 100)

	switch option {
	case UDT_MSS:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_MSS,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_SNDSYN:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_SNDSYN,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_RCVSYN:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_RCVSYN,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_CC:
		{
			return -1, fmt.Errorf("UDT CCC Not implemented.")
		}
	case UDT_FC:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_FC,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_SNDBUF:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_SNDBUF,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_RCVBUF:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_RCVBUF,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDP_SNDBUF:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDP_SNDBUF,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDP_RCVBUF:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDP_RCVBUF,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_LINGER:
		{
			var clinger C.struct_linger
			var clinger_len = C.int(unsafe.Sizeof(clinger))
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_LINGER,
				unsafe.Pointer(&clinger), &clinger_len))
			if retval < 0 {
				return retval, udtErrDesc("Unable to get option")
			}
			linger := Linger{
				l_onoff:  int(clinger.l_onoff),
				l_linger: int(clinger.l_linger),
			}

			return linger, nil
		}

	case UDT_RENDEZVOUS:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_RENDEZVOUS,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_SNDTIMEO:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_SNDTIMEO,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_RCVTIMEO:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_RCVTIMEO,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_REUSEADDR:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_REUSEADDR,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_MAXBW:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_MAXBW,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_STATE:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_STATE,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_EVENT:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_EVENT,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_SNDDATA:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_SNDDATA,
				unsafe.Pointer(&data[0]), &optlen))
		}
	case UDT_RCVDATA:
		{
			retval = int(C.udt_getsockopt(socket.sock, C.int(0), C.UDT_UDT_RCVDATA,
				unsafe.Pointer(&data[0]), &optlen))
		}
	default:
		{
			return -1, fmt.Errorf("Invalid option %s", option)
		}
	}

	if retval < 0 {
		return retval, udtErrDesc("Unable to get option")
	}

	optlengo := int(optlen)
	switch optlengo {
	case 1:
		{
			value, _ = binary.Uvarint(data[:optlengo])
		}
	case 4:
		{
			value = binary.LittleEndian.Uint16(data[:optlengo])
		}
	case 8:
		{
			value = binary.LittleEndian.Uint64(data[:optlengo])
		}
	}

	//fmt.Printf("Got data %d length %d option %s\n", value, optlengo,option)

	return
}

//This method sets requested UDT socket option. If successful, returns requested option value otherwise
//returns error object with error details.

func Setsockopt(socket *Socket, option string, value interface{}) (retval int, err error) {
	var data []byte
	if option != UDT_LINGER {
		data, err = getBytes(value)
		if err != nil {
			return -1, fmt.Errorf("Unable to convert interface to byte array %s", err)
		}
	}
	switch option {
	case UDT_MSS:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint16 {
				return -1, fmt.Errorf("Requires Uint16 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_MSS,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_SNDSYN:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint64 {
				return -1, fmt.Errorf("Requires Uint64 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_SNDSYN,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}
	case UDT_RCVSYN:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint64 {
				return -1, fmt.Errorf("Requires Uint64 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_RCVSYN,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_CC:
		{
			return -1, fmt.Errorf("UDT CCC Not implemented.")
		}

	case UDT_FC:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint16 {
				return -1, fmt.Errorf("Requires Uint16 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_FC,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_SNDBUF:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint16 {
				return -1, fmt.Errorf("Requires Uint16 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_SNDBUF,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_RCVBUF:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint16 {
				return -1, fmt.Errorf("Requires Uint16 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_RCVBUF,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDP_SNDBUF:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint16 {
				return -1, fmt.Errorf("Requires Uint16 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDP_SNDBUF,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}
	case UDP_RCVBUF:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint16 {
				return -1, fmt.Errorf("Requires Uint16 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDP_RCVBUF,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_LINGER:
		{
			var clingerin Linger = value.(Linger)
			var clinger C.struct_linger
			clinger.l_onoff = C.int(clingerin.l_onoff)
			clinger.l_linger = C.int(clingerin.l_linger)
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_LINGER,
				unsafe.Pointer(&clinger), C.int(unsafe.Sizeof(clinger))))
		}

	case UDT_RENDEZVOUS:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint64 {
				return -1, fmt.Errorf("Requires Uint64 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_RENDEZVOUS,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_REUSEADDR:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint64 {
				return -1, fmt.Errorf("Requires Uint64 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_REUSEADDR,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	case UDT_MAXBW:
		{
			if reflect.TypeOf(value).Kind() != reflect.Uint64 {
				return -1, fmt.Errorf("Requires Uint64 type")
			}
			retval = int(C.udt_setsockopt(socket.sock, C.int(0), C.UDT_UDT_MAXBW,
				unsafe.Pointer(&data[0]), C.int(len(data))))
		}

	default:
		{
			return -1, fmt.Errorf("Invalid option %s", option)
		}
	}

	if retval < 0 {
		return retval, udtErrDesc("Unable set option")
	}

	return
}

//This method retrieves the address informtion of the peer side of a connected UDT socket. If successful returns
//peer socket address otherwise returns error object with error details.

func Getpeername(socket *Socket) (sockaddr Sockaddr, err error) {

	var sockaddr_in C.struct_sockaddr_in
	var namelen C.int

	retval := int(C.udt_getpeername(socket.sock,
		(*C.struct_sockaddr)(unsafe.Pointer(&sockaddr_in)), &namelen))
	if retval < 0 {
		return sockaddr, udtErrDesc("Unable to get socket peername")
	}

	sockaddr = Sockaddr{
		sa_family: int(sockaddr_in.sin_family),
		sa_data:   C.GoString(C.inet_ntoa(sockaddr_in.sin_addr)),
	}

	return

}

//This method retrieves the address informtion of the UDT socket. If successful returns
//socket address otherwise returns error object with error details.

func Getsockname(socket *Socket) (sockaddr Sockaddr, err error) {

	var sockaddr_in C.struct_sockaddr_in
	var namelen C.int

	retval := int(C.udt_getsockname(socket.sock,
		(*C.struct_sockaddr)(unsafe.Pointer(&sockaddr_in)), &namelen))

	if retval < 0 {
		return sockaddr, udtErrDesc("Unable to get socket name")
	}

	sockaddr = Sockaddr{
		sa_family: int(sockaddr_in.sin_family),
		sa_data:   C.GoString(C.inet_ntoa(sockaddr_in.sin_addr)),
	}

	return
}

//This method retrieves the internal protocol parameters and performance trace. If successful returns
// Traceinfo struct otherwise returns error object with error details.

func Perfmon(socket *Socket, clear bool) (traceinfo Traceinfo, err error) {

	var cClear C.int = 0
	if clear {
		cClear = 1
	}

	var udtTraceinfo C.UDT_TRACEINFO

	retval := int(C.udt_perfmon(socket.sock,
		(*C.UDT_TRACEINFO)(unsafe.Pointer(&udtTraceinfo)), cClear))
	if retval < 0 {
		return traceinfo, udtErrDesc("Unable to get trace info")
	}

	traceinfo = Traceinfo{
		msTimeStamp:         int64(udtTraceinfo.msTimeStamp),
		pktSentTotal:        int64(udtTraceinfo.pktSentTotal),
		pktRecvTotal:        int64(udtTraceinfo.pktRecvTotal),
		pktSndLossTotal:     int(udtTraceinfo.pktSndLossTotal),
		pktRcvLossTotal:     int(udtTraceinfo.pktRcvLossTotal),
		pktRetransTotal:     int(udtTraceinfo.pktRetransTotal),
		pktSentACKTotal:     int(udtTraceinfo.pktSentACKTotal),
		pktRecvACKTotal:     int(udtTraceinfo.pktRecvACKTotal),
		pktSentNAKTotal:     int(udtTraceinfo.pktSentNAKTotal),
		pktRecvNAKTotal:     int(udtTraceinfo.pktRecvNAKTotal),
		usSndDurationTotal:  int64(udtTraceinfo.usSndDurationTotal),
		pktSent:             int64(udtTraceinfo.pktSent),
		pktRecv:             int64(udtTraceinfo.pktRecv),
		pktSndLoss:          int(udtTraceinfo.pktSndLoss),
		pktRcvLoss:          int(udtTraceinfo.pktRcvLoss),
		pktRetrans:          int(udtTraceinfo.pktRetrans),
		pktSentACK:          int(udtTraceinfo.pktSentACK),
		pktRecvACK:          int(udtTraceinfo.pktRecvACK),
		pktSentNAK:          int(udtTraceinfo.pktSentNAK),
		pktRecvNAK:          int(udtTraceinfo.pktRecvNAK),
		mbpsSendRate:        float64(udtTraceinfo.mbpsSendRate),
		mbpsRecvRate:        float64(udtTraceinfo.mbpsRecvRate),
		usSndDuration:       int64(udtTraceinfo.usSndDuration),
		usPktSndPeriod:      float64(udtTraceinfo.usPktSndPeriod),
		pktFlowWindow:       int(udtTraceinfo.pktFlowWindow),
		pktCongestionWindow: int(udtTraceinfo.pktCongestionWindow),
		pktFlightSize:       int(udtTraceinfo.pktFlightSize),
		msRTT:               float64(udtTraceinfo.msRTT),
		mbpsBandwidth:       float64(udtTraceinfo.mbpsBandwidth),
		byteAvailSndBuf:     int(udtTraceinfo.byteAvailSndBuf),
		byteAvailRcvBuf:     int(udtTraceinfo.byteAvailRcvBuf),
	}

	return

}

//This methods creates epoll. If successful returns epoll id
//otherwise returns error object with error details.

func EpollCreate() (eid int, err error) {
	eid = int(C.udt_epoll_create())
	if eid < 0 {
		return eid, udtErrDesc("Unable create new epoll ID")
	}
	return
}

// This method binds epoll to provided UDT socket and watches provided event. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.


func EpollAddUsock(eid int, socket *Socket, events int) (retval int, err error) {

	retval = int(C.udt_epoll_add_usock(C.int(eid), socket.sock,
		(*C.int)(unsafe.Pointer(&events))))
	if retval < 0 {
		return retval, udtErrDesc("Unable to add UDT socket for epoll")
	}
	return

}

// This method binds epoll to provided system socket and watches provided event. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.


func EpollAddSsock(eid int, socket C.SYSSOCKET, events int) (retval int, err error) {

	retval = int(C.udt_epoll_add_ssock(C.int(eid), socket,
		(*C.int)(unsafe.Pointer(&events))))
	if retval < 0 {
		return retval, udtErrDesc("Unable to add sys socket for epoll")
	}
	return

}

// This method removes epoll from UDT socket. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func EpollRemoveUsock(eid int, socket *Socket) (retval int, err error) {

	retval = int(C.udt_epoll_remove_usock(C.int(eid), socket.sock))
	if retval < 0 {
		return retval, udtErrDesc("Unable to remove UDT socket")
	}
	return

}

// This method removes epoll from sysmtem socket. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func EpollRemoveSsock(eid int, socket C.SYSSOCKET) (retval int, err error) {

	retval = int(C.udt_epoll_remove_ssock(C.int(eid), socket))
	if retval < 0 {
		return retval, udtErrDesc("Unable to remove sys socket")
	}
	return

}

//This method add wait on give read and write UDT and ystem sockets. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func EpollWait2(eid int, readfds []C.UDTSOCKET, writefds []C.UDTSOCKET, msTimeOut int64,
	lrfds []C.SYSSOCKET, lwfds []C.SYSSOCKET) (retval int, err error) {
	rnum := C.int(len(readfds))
	wnum := C.int(len(writefds))
	lrnum := C.int(len(lrfds))
	lwnum := C.int(len(lwfds))

	retval = int(C.udt_epoll_wait2(C.int(eid), (*C.UDTSOCKET)(unsafe.Pointer(&readfds[0])),
		&rnum, (*C.UDTSOCKET)(unsafe.Pointer(&writefds[0])), &wnum, C.int64_t(msTimeOut),
		(*C.SYSSOCKET)(unsafe.Pointer(&lrfds[0])), &lrnum,
		(*C.SYSSOCKET)(unsafe.Pointer(&lwfds[0])), &lwnum))
	if retval < 0 {
		return retval, udtErrDesc("Unable to epoll wait")
	}
	return
}

//This method release epoll. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func EpollRelease(eid int) (retval int, err error) {
	retval = int(C.udt_epoll_release(C.int(eid)))
	if retval < 0 {
		return retval, udtErrDesc("Unable to release socket")
	}
	return
}

//This method closes requested UDT socket. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func Close(socket *Socket) (retval int, err error) {
	retval = int(C.udt_close(socket.sock))
	if retval < 0 {
		return retval, udtErrDesc("Unable to close socket")
	}
	return
}

//This method starts UDT system. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.

func Startup() (retval int, err error) {
	retval = int(C.udt_startup())
	if retval != 0 {
		return retval, udtErrDesc("Unable to startup UDT library")
	}
	return
}

//This method cleans UDT system. If successful,
// this method returns 0, otherwise it returns error code (http://udt.sourceforge.net/udt4/doc/ecode.htm)
// and error object with error details.


func Cleanup() (retval int, err error) {
	retval = int(C.udt_cleanup())
	if retval != 0 {
		return retval, udtErrDesc("Unable to execute cleanup")
	}
	return
}

//This method clears last error.

func Clearlasterror() {
	C.udt_clearlasterror()
}

//This method retrieves recent error from UDT sysytem.

func udtErrDesc(appMsg string) (err error) {
	return fmt.Errorf("%s - UDT Error-%d:%s ", appMsg,
		C.GoString(C.udt_getlasterror_desc()), int(C.udt_getlasterror_code()))
}

//Utility method converts boolean to int.

func boolToInt(boolvalue bool) (boolint int) {
	if boolvalue {
		return 1
	}

	return 0
}

func getBytes(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//This method creates UDT socket.

func CreateSockets(size int) (sockets Sockets) {
	sockets = Sockets{
		socks: make([]C.UDTSOCKET, size),
	}
	return
}

//This method creates system socket.

func CreateSysSockets(size int) (sockets SysSockets) {
	sockets = SysSockets{
		socks: make([]C.SYSSOCKET, size),
	}
	return
}

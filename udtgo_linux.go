package udtgo

import (
	"syscall"
	"unsafe"
)

type Socketlen uint

//Converts Sockaddr to system level RawSockaddrAny which matches c strucure.
func SockaddrToRawSockAny(sa syscall.Sockaddr) (*syscall.RawSockaddrAny, Socketlen, error) {
	if sa == nil {
		return nil, 0, syscall.EINVAL
	}

	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		if sa.Port < 0 || sa.Port > 65535 {
			return nil, 0, syscall.EINVAL
		}
		var raw syscall.RawSockaddrInet4
		raw.Family = syscall.AF_INET
		pport := (*[2]byte)(unsafe.Pointer(&raw.Port))
		pport[0] = byte(sa.Port >> 8)
		pport[1] = byte(sa.Port)
		for i := 0; i < len(sa.Addr); i++ {
			raw.Addr[i] = sa.Addr[i]
		}
		return (*syscall.RawSockaddrAny)(unsafe.Pointer(&raw)), syscall.SizeofSockaddrInet4, nil

	case *syscall.SockaddrInet6:
		if sa.Port < 0 || sa.Port > 65535 {
			return nil, 0, syscall.EINVAL
		}
		var raw syscall.RawSockaddrInet6
		raw.Family = syscall.AF_INET6
		pport := (*[2]byte)(unsafe.Pointer(&raw.Port))
		pport[0] = byte(sa.Port >> 8)
		pport[1] = byte(sa.Port)
		raw.Scope_id = sa.ZoneId
		for i := 0; i < len(sa.Addr); i++ {
			raw.Addr[i] = sa.Addr[i]
		}
		return (*syscall.RawSockaddrAny)(unsafe.Pointer(&raw)), syscall.SizeofSockaddrInet6, nil

	}
	return nil, 0, syscall.EAFNOSUPPORT
}

//Parses address string from RawSockaddrAny strucure

func parseAddr(rsa *syscall.RawSockaddrAny) (string, error) {
	switch rsa.Addr.Family {


	case syscall.AF_INET:
		prsa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))

		sa := new(syscall.SockaddrInet4)
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = prsa.Addr[i]
		}
		return ip4String(sa.Addr), nil

	case syscall.AF_INET6:
		prsa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(syscall.SockaddrInet6)
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = prsa.Addr[i]
		}

		return  ip6String(sa.Addr), nil
	}
	return "", syscall.EAFNOSUPPORT
}

//Converts [4]byte ipv4 address to string

func ip4String(p [4]byte) string {


	if len(p) == 0 {
		return "<nil>"
	}


	return uitoa(uint(p[0])) + "." +
			uitoa(uint(p[1])) + "." +
			uitoa(uint(p[2])) + "." +
			uitoa(uint(p[3]))

}

//Converts [16]byte ipv6 address to string

func ip6String(p [16]byte) string {

	IPv6len := 16



	if len(p) == 0 {
		return "<nil>"
	}


	// Find longest run of zeros.
	e0 := -1
	e1 := -1
	for i := 0; i < IPv6len; i += 2 {
		j := i
		for j < IPv6len && p[j] == 0 && p[j+1] == 0 {
			j += 2
		}
		if j > i && j-i > e1-e0 {
			e0 = i
			e1 = j
			i = j
		}
	}
	// The symbol "::" MUST NOT be used to shorten just one 16 bit 0 field.
	if e1-e0 <= 2 {
		e0 = -1
		e1 = -1
	}

	const maxLen = len("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
	b := make([]byte, 0, maxLen)

	// Print with possible :: in place of run of zeros
	for i := 0; i < IPv6len; i += 2 {
		if i == e0 {
			b = append(b, ':', ':')
			i = e1
			if i >= IPv6len {
				break
			}
		} else if i > 0 {
			b = append(b, ':')
		}
		b = appendHex(b, (uint32(p[i])<<8)|uint32(p[i+1]))
	}
	return string(b)
}



func uitoa(val uint) string {
	var buf [32]byte
	i := len(buf) - 1
	for val >= 10 {
		buf[i] = byte(val%10 + '0')
		i--
		val /= 10
	}
	buf[i] = byte(val + '0')
	return string(buf[i:])
}


func appendHex(dst []byte, i uint32) []byte {
	hexDigit := "0123456789abcdef"
	if i == 0 {
		return append(dst, '0')
	}
	for j := 7; j >= 0; j-- {
		v := i >> uint(j*4)
		if v > 0 {
			dst = append(dst, hexDigit[v&0xf])
		}
	}
	return dst
}
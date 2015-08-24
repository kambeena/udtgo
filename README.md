# udtgo

This is a cgo wrapper for UDT (http://udt.sourceforge.net/). Compile udt4c and 
install in /usr/local/lib directory for linux and system directory for windows. The name of complied 
library should have name libudt.so or change #cgo LDFLAGS in udt.go. Here you will find a compiled library for linux amd64 architecture in libs folder. Please copy 
the file to /local/lib and rename the file to libudt.so.

For compiling udt use udt source code in udt4c and follow the instructions at 
http://udt.sourceforge.net/udt4/index.htm.



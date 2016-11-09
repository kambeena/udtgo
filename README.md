# udtgo

This is a cgo wrapper for UDT (http://udt.sourceforge.net/). Compile udt4c and install in /usr/local/lib directory for linux. 
Currently no support for windows system. The name of the complied library should be libudt.so or change #cgo LDFLAGS in udt.go. 
With source code you will find a compiled library for linux amd64 architecture in libs folder. Copy the file to /usr/local/lib and 
rename the file to libudt.so.

For compiling udt use udt source code in udt4c and follow the instructions at http://udt.sourceforge.net/udt4/index.htm.

All UDT functionality is ported udtgo except ability to set User-defined Congestion Control Algorithm. File upload examples are in the examples directory.

This cgo wrapper for UDT ((http://udt.sourceforge.net/) is available under BSD license.



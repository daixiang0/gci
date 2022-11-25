package main

import (
	"fmt"

	"github.com/daixiang0/gci"
	g "github.com/golang"
)

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"

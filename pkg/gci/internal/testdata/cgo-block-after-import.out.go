package main

// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"

import (
	"fmt"

	g "github.com/golang"

	"github.com/daixiang0/gci"
)

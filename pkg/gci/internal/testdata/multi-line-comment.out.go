package proc

import (
	"context" // in-line comment
	"fmt"
	"os"
	//nolint:depguard // A multi-line comment explaining why in
	// this one case it's OK to use os/exec even though depguard
	// is configured to force us to use dlib/exec instead.
	"os/exec"

	"github.com/local/dlib/dexec"
	"golang.org/x/sys/unix"
)

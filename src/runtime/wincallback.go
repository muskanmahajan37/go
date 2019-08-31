// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Generate Windows callback assembly file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

const maxCallback = 2000

func genasm386Amd64() {
	var buf bytes.Buffer

	buf.WriteString(`// Code generated by wincallback.go using 'go generate'. DO NOT EDIT.

// +build 386 amd64
// runtime·callbackasm is called by external code to
// execute Go implemented callback function. It is not
// called from the start, instead runtime·compilecallback
// always returns address into runtime·callbackasm offset
// appropriately so different callbacks start with different
// CALL instruction in runtime·callbackasm. This determines
// which Go callback function is executed later on.

TEXT runtime·callbackasm(SB),7,$0
`)
	for i := 0; i < maxCallback; i++ {
		buf.WriteString("\tCALL\truntime·callbackasm1(SB)\n")
	}

	filename := fmt.Sprintf("zcallback_windows.s")
	err := ioutil.WriteFile(filename, buf.Bytes(), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wincallback: %s\n", err)
		os.Exit(2)
	}
}

func genasmArm() {
	var buf bytes.Buffer

	buf.WriteString(`// Code generated by wincallback.go using 'go generate'. DO NOT EDIT.

// External code calls into callbackasm at an offset corresponding
// to the callback index. Callbackasm is a table of MOV and B instructions.
// The MOV instruction loads R12 with the callback index, and the
// B instruction branches to callbackasm1.
// callbackasm1 takes the callback index from R12 and
// indexes into an array that stores information about each callback.
// It then calls the Go implementation for that callback.
#include "textflag.h"

TEXT runtime·callbackasm(SB),NOSPLIT|NOFRAME,$0
`)
	for i := 0; i < maxCallback; i++ {
		buf.WriteString(fmt.Sprintf("\tMOVW\t$%d, R12\n", i))
		buf.WriteString("\tB\truntime·callbackasm1(SB)\n")
	}

	err := ioutil.WriteFile("zcallback_windows_arm.s", buf.Bytes(), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wincallback: %s\n", err)
		os.Exit(2)
	}
}

//todo(ragav): check for correctness
func genasmArm64() {
	var buf bytes.Buffer

	buf.WriteString(`// Code generated by wincallback.go using 'go generate'. DO NOT EDIT.

// External code calls into callbackasm at an offset corresponding
// to the callback index. Callbackasm is a table of MOV and B instructions.
// The MOV instruction loads R16 with the callback index, and the
// B instruction branches to callbackasm1.
// callbackasm1 takes the callback index from R16 and
// indexes into an array that stores information about each callback.
// It then calls the Go implementation for that callback.
#include "textflag.h"`

TEXT runtime·callbackasm(SB),NOSPLIT|NOFRAME,$0
`)
	for i := 0; i < maxCallback; i++ {
		buf.WriteString(fmt.Sprintf("\tMOVD\t$%d, R16\n", i))
		buf.WriteString("\tB\truntime·callbackasm1(SB)\n")
	}

	err := ioutil.WriteFile("zcallback_windows_arm64.s", buf.Bytes(), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wincallback: %s\n", err)
		os.Exit(2)
	}
}

func gengo() {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`// Code generated by wincallback.go using 'go generate'. DO NOT EDIT.

package runtime

const cb_max = %d // maximum number of windows callbacks allowed
`, maxCallback))
	err := ioutil.WriteFile("zcallback_windows.go", buf.Bytes(), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wincallback: %s\n", err)
		os.Exit(2)
	}
}

func main() {
	genasm386Amd64()
	genasmArm()
	genasmArm64()
	gengo()
}

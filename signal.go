// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"sync"
)

var interrupted = make(chan struct{})
var subProcExited = make(chan struct{})

func processSignals() {
	sig := make(chan os.Signal)
	signal.Notify(sig, signalsToIgnore...)
	go func() {
		<-sig
	}()
}

var onceProcessSignals sync.Once

func startSigHandlers() {
	onceProcessSignals.Do(processSignals)
}

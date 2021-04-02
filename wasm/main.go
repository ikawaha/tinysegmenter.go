package main

import (
	"syscall/js"

	"github.com/ikawaha/tinysegmenter.go"
)

const exportSegmentFuncName = "waSegment"

func segment(_ js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}
	var ret []interface{}
	ss := tinysegmenter.Segment(args[0].String())
	for _, v := range ss {
		ret = append(ret, v)
	}
	return ret
}

func registerCallbacks() {
	js.Global().Set(exportSegmentFuncName, js.FuncOf(segment))
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}

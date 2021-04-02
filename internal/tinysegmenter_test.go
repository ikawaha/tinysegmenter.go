package internal

import (
	_ "embed"
	"reflect"
	"testing"
)

//go:embed testdata/timemachineu8j.txt
var sampletext string

func TestSegment(t *testing.T) {
	ary := Segment("私の名前は中野です")
	expect := []string{"私", "の", "名前", "は", "中野", "です"}
	if !reflect.DeepEqual(ary, expect) {
		t.Errorf("got %+v, expected %v", ary, expect)
	}
}

func BenchmarkSegment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Segment("私の名前は中野です")
	}
}

func BenchmarkSegmentLargeText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Segment(sampletext)
	}
}

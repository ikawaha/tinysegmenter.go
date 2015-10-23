package tinysegmenter

import (
	"io/ioutil"
	"reflect"
	"testing"
)

var sampletext string

func init() {
	b, err := ioutil.ReadFile("./timemachineu8j.txt")
	if err != nil {
		panic(err)
	}
	sampletext = string(b)
}

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

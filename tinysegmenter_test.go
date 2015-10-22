package tinysegmenter

import (
	"testing"
)

func BenchmarkSegment(b *testing.B) {
	s := NewSegmenter()
	for i := 0; i < b.N; i++ {
		s.Segment("私の名前は中野です")
	}
}

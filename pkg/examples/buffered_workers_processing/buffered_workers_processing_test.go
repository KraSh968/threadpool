package bufferedworkersprocessing

import "testing"

func TestBaseExample(t *testing.T) {
	RunBaseExample()
}

func BenchmarkBaseExample(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	RunBaseExample()
}

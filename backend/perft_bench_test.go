package backend

import (
	"testing"
)

func BenchmarkPerftD1(b *testing.B) {
	pos := InitialPosition
	for i := 0; i < b.N; i++ {
		Perft(&pos, 1)
	}
}

func BenchmarkPerftD2(b *testing.B) {
	pos := InitialPosition
	for i := 0; i < b.N; i++ {
		Perft(&pos, 2)
	}
}

func BenchmarkPerftD3(b *testing.B) {
	pos := InitialPosition
	for i := 0; i < b.N; i++ {
		Perft(&pos, 3)
	}
}

func BenchmarkPerftD4(b *testing.B) {
	pos := InitialPosition
	for i := 0; i < b.N; i++ {
		Perft(&pos, 4)
	}
}

func BenchmarkPerftD5(b *testing.B) {
	pos := InitialPosition
	for i := 0; i < b.N; i++ {
		Perft(&pos, 5)
	}
}

func BenchmarkPerftD6(b *testing.B) {
	pos := InitialPosition
	for i := 0; i < b.N; i++ {
		Perft(&pos, 6)
	}
}

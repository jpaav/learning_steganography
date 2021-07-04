package secrets

import (
	"bytes"
	"testing"
)

func TestWriter_WriteCrumb(t *testing.T) {
	const expectedBytes = "1010101010101010101010101010101010101010101010101010101010101010"
	bufferBytes := make([]byte, 0, 8)
	buffer := bytes.NewBuffer(bufferBytes)
	secretWriter := NewWriter(buffer, 3)
	for i := 0; i < 32; i++ {
		err := secretWriter.WriteCrumb(0b10)
		if err != nil {
			t.Fatal(err)
		}
	}
	err := secretWriter.Sink()
	if err != nil {
		t.Fatal(err)
	}

	if byteBufferToBinaryString(buffer.Bytes()) != expectedBytes {
		t.Log(byteBufferToBinaryString(buffer.Bytes()))
		t.Log(expectedBytes)
		t.Fatal("written buffer does not have expected value")
	}
}

func benchmarkWriteCrumb(payload []byte, bufferSize int, b *testing.B) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		writeBuffer := make([]byte, 0, 3)
		writer := bytes.NewBuffer(writeBuffer)
		secretWriter := Writer{writer: writer, buffer: make([]byte, bufferSize), bitsIndex: 0}
		for _, crumb := range payload {
			err := secretWriter.WriteCrumb(crumb)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkWriter_WriteCrumb_Short1(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremShort), 1, b)
}

func BenchmarkWriter_WriteCrumb_Short3(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremShort), 3, b)
}

func BenchmarkWriter_WriteCrumb_Short32(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremShort), 32, b)
}

func BenchmarkWriter_WriteCrumb_Short1000(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremShort), 1000, b)
}

func BenchmarkWriter_WriteCrumb_Long1(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremLong), 1, b)
}

func BenchmarkWriter_WriteCrumb_Long3(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremLong), 3, b)
}

func BenchmarkWriter_WriteCrumb_Long32(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremLong), 32, b)
}

func BenchmarkWriter_WriteCrumb_Long1000(b *testing.B) {
	benchmarkWriteCrumb([]byte(loremLong), 1000, b)
}

func BenchmarkWriter_WriteCrumb_ExtraLong1(b *testing.B) {
	benchmarkWriteCrumb(make([]byte, 10000), 1, b)
}

func BenchmarkWriter_WriteCrumb_ExtraLong3(b *testing.B) {
	benchmarkWriteCrumb(make([]byte, 10000), 3, b)
}

func BenchmarkWriter_WriteCrumb_ExtraLong32(b *testing.B) {
	benchmarkWriteCrumb(make([]byte, 10000), 32, b)
}

func BenchmarkWriter_WriteCrumb_ExtraLong1000(b *testing.B) {
	benchmarkWriteCrumb(make([]byte, 10000), 1000, b)
}

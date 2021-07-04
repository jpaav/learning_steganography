package secrets

import (
	"fmt"
	"strings"
	"testing"
)

const loremShort = "Lorem ipsum"
const loremShortAligned = "Lorem ipsum."
const loremLong = "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea"

func byteBufferToBinaryString(buffer []byte) string {
	var result string
	for _, b := range buffer {
		result = fmt.Sprintf("%s%08b", result, b)
	}
	return result
}

func readStringAsCrumbs(payload string, bufferSize int, t *testing.T) {
	t.Helper()
	// Construct payload
	var binaryPayload = byteBufferToBinaryString([]byte(payload))

	var crumb, curByte byte
	var err error
	var binaryString string
	var resultBuffer []byte
	reader := strings.NewReader(payload)
	secretReader := Reader{reader: reader, buffer: make([]byte, bufferSize), bitsLeft: 0}
	index := 0
	for {
		crumb, err = secretReader.ReadCrumb()
		if err != nil {
			resultBuffer = append(resultBuffer, curByte)
			break
		}
		// Fill the result buffer
		if index*2%8 == 0 {
			if index != 0 {
				resultBuffer = append(resultBuffer, curByte)
			}
			curByte = crumb
		} else {
			curByte <<= 2
			curByte |= crumb
		}
		binaryString = fmt.Sprintf("%s%02b", binaryString, crumb)

		index += 1
	}
	if string(resultBuffer) != payload {
		t.Log(string(resultBuffer))
		t.Log(payload)
		t.Error("output string not the same as input string")
	}
	if binaryPayload != binaryString {
		t.Log(binaryPayload)
		t.Log(binaryString)
		t.Error("binaryPayload not the same as binaryString")
	}
}

func TestReader_ReadCrumb(t *testing.T) {
	readStringAsCrumbs(loremShort, 3, t)
	readStringAsCrumbs(loremShortAligned, 3, t)
	readStringAsCrumbs(loremLong, 3, t)

	readStringAsCrumbs(loremShort, 5, t)
	readStringAsCrumbs(loremShortAligned, 7, t)
	readStringAsCrumbs(loremLong, 1, t)
}

func benchmarkReadCrumb(payload string, bufferSize int, b *testing.B) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(payload)
		secretReader := Reader{reader: reader, buffer: make([]byte, bufferSize), bitsLeft: 0}
		for {
			_, err := secretReader.ReadCrumb()
			if err != nil {
				break
			}
		}
	}
}

func BenchmarkReader_ReadCrumb_Short1(b *testing.B) {
	benchmarkReadCrumb(loremShort, 1, b)
}

func BenchmarkReader_ReadCrumb_Short3(b *testing.B) {
	benchmarkReadCrumb(loremShort, 3, b)
}

func BenchmarkReader_ReadCrumb_Short5(b *testing.B) {
	benchmarkReadCrumb(loremShort, 5, b)
}

func BenchmarkReader_ReadCrumb_Short32(b *testing.B) {
	benchmarkReadCrumb(loremShort, 32, b)
}

func BenchmarkReader_ReadCrumb_Long1(b *testing.B) {
	benchmarkReadCrumb(loremLong, 1, b)
}

func BenchmarkReader_ReadCrumb_Long3(b *testing.B) {
	benchmarkReadCrumb(loremLong, 3, b)
}

func BenchmarkReader_ReadCrumb_Long5(b *testing.B) {
	benchmarkReadCrumb(loremLong, 5, b)
}

func BenchmarkReader_ReadCrumb_Long32(b *testing.B) {
	benchmarkReadCrumb(loremLong, 32, b)
}

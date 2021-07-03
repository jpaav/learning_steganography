package secrets

import (
	"errors"
	"io"
)

const CrumbBitMask uint8 = 0b11
const NegatedCrumbBitMask uint8 = 0b11111100

const CrumbBitCount = 2
const ByteBitCount = 8

type Reader struct {
	reader    io.Reader // reader reads the secret file.
	buffer    []byte    // buffer helps helps us incrementally read the secret file.
	bytesLeft int       // bytesLeft tracks how many bytes are available in the buffer. When the buffer is partially filled, it will differ from len(buffer).
	bitsLeft  int       // bitsLeft tracks how many bits are available in the buffer.
}

func (r *Reader) fillBuffer() error {
	count, err := io.ReadFull(r.reader, r.buffer)
	r.bytesLeft = count
	r.bitsLeft = count * ByteBitCount
	return err
}

// ReadCrumb reads the crumb and returns it in the least significant two bits.
// If the next crumb cannot be read, the appropriate error is returned from the io package.
func (r *Reader) ReadCrumb() (byte, error) {
	if r.bitsLeft < CrumbBitCount {
		err := r.fillBuffer()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return 0, err
			}
			// If we get io.ErrUnexpectedEOF we still might have enough data to continue
		}
	}
	// Check that we read enough data to return a crumb
	if r.bitsLeft >= CrumbBitCount {
		// Get the byte out of the buffer
		curBit := r.bytesLeft*ByteBitCount - r.bitsLeft
		curByte := r.buffer[curBit/ByteBitCount]
		// Get the right bits out of the byte
		// curBit = 0: shift right 6
		// curBit = 2: shift right 4
		// curBit = 4: shift right 2
		// curBit = 6: shift right 0
		curByte >>= ByteBitCount - (CrumbBitCount + curBit%ByteBitCount)
		// Grab the last two bits
		curByte &= CrumbBitMask
		r.bitsLeft -= CrumbBitCount
		return curByte, nil
	}
	return 0, io.ErrUnexpectedEOF
}

func NewReader(reader io.Reader, bufferSize int) *Reader {
	return &Reader{reader: reader, buffer: make([]byte, bufferSize), bitsLeft: 0}
}

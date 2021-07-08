package secrets

import "io"

// Writer facilitates writing less than a full byte at a time.
type Writer struct {
	writer    io.Writer // writer writes the secrets
	buffer    []byte    // buffer helps helps us incrementally write the secret file.
	bitsIndex int       // bitsIndex tracks the current bit index being written in the buffer.
}

// sinkBuffer is an internal method for writing a certain number of bytes from the buffer to the writer.
// If writeCount specifies to write more bytes than are in the buffer, io.ErrShortWrite is returned.
func (w *Writer) sinkBuffer(writeCount int) error {
	if writeCount > len(w.buffer) {
		return io.ErrShortWrite
	}
	_, err := w.writer.Write(w.buffer[:writeCount])
	w.bitsIndex = 0
	return err
}

// WriteCrumb writes a crumb to the buffer.
// If the buffer is full, it will write it and reset itself.
func (w *Writer) WriteCrumb(crumb byte) error {
	w.buffer[w.bitsIndex/8] <<= 2
	w.buffer[w.bitsIndex/8] |= crumb & CrumbBitMask
	w.bitsIndex += 2
	if w.bitsIndex == len(w.buffer)*ByteBitCount {
		err := w.sinkBuffer(len(w.buffer))
		if err != nil {
			return err
		}
		w.bitsIndex = 0
	}
	return nil
}

// Sink writes all the bytes that have been filled.
func (w *Writer) Sink() error {
	return w.sinkBuffer(w.bitsIndex / 8)
}

// NewWriter creates a Writer wrapping the given writer and with a buffer of given size.
func NewWriter(writer io.Writer, bufferSize int) *Writer {
	return &Writer{writer: writer, buffer: make([]byte, bufferSize), bitsIndex: 0}
}

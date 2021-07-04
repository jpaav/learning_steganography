package secrets

import "io"

type Writer struct {
	writer    io.Writer // writer writes the secrets
	buffer    []byte    // buffer helps helps us incrementally write the secret file.
	bitsIndex int       // bitsIndex tracks the current bit index being written in the buffer.
}

func (w *Writer) sinkBuffer(writeCount int) error {
	_, err := w.writer.Write(w.buffer[:writeCount])
	w.bitsIndex = 0
	return err
}

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

func (w *Writer) Sink() error {
	return w.sinkBuffer(w.bitsIndex / 8)
}

func NewWriter(writer io.Writer, bufferSize int) *Writer {
	return &Writer{writer: writer, buffer: make([]byte, bufferSize), bitsIndex: 0}
}

package procker

import "io"

// PrefixedWriter implements prefixed output for an io.Writer object.
type PrefixedWriter struct {
	Prefix string
	writer io.Writer
	inline bool
}

// NewPrefixedWriter creates a PrefixedWriter
func NewPrefixedWriter(w io.Writer, prefix string) io.Writer {
	return &PrefixedWriter{Prefix: prefix, writer: w}
}

// Writes a Prefix string before writing to the underlying writer.
func (w *PrefixedWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if !w.inline {
			io.WriteString(w.writer, w.Prefix)
		}
		w.writer.Write([]byte{b})
		w.inline = b != '\n'
	}
	return len(p), nil
}

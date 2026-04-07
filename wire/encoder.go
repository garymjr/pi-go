package wire

import "io"

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder { return &Encoder{w: w} }

func (e *Encoder) Encode(v any) error {
	line, err := EncodeLine(v)
	if err != nil {
		return err
	}
	_, err = e.w.Write(line)
	return err
}

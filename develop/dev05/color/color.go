package color

import (
	"fmt"
	"io"
)

const (
	colorDefault = "\033[0m"
	colorRed     = "\033[31m"
	colorBlue    = "\033[34m"
)

// Wrapper for color settings
type Wrapper struct {
	needColor bool
}

// NewWrapper used to create wrapper for color settings
func NewWrapper(need bool) Wrapper {
	return Wrapper{need}
}

// Reset used to set color to default
func (cw Wrapper) Reset(w io.Writer) {
	if cw.needColor {
		Reset(w)
	}
}

// SetRed used to set color to red
func (cw Wrapper) SetRed(w io.Writer) {
	if cw.needColor {
		SetRed(w)
	}
}

// SetBlue used to set color to blue
func (cw Wrapper) SetBlue(w io.Writer) {
	if cw.needColor {
		SetBlue(w)
	}
}

// Reset used to set color to default
func Reset(w io.Writer) {
	fmt.Fprint(w, colorDefault)
}

// SetRed used to set color to red
func SetRed(w io.Writer) {
	fmt.Fprint(w, colorRed)
}

// SetBlue used to set color to blue
func SetBlue(w io.Writer) {
	fmt.Fprint(w, colorBlue)
}

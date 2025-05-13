package miln

import (
	"errors"

	"gonum.org/v1/plot/plotter"
)

type MilnEquasion struct {
	Kx          float64 `json:"kx"`
	Ky          float64 `json:"ky"`
	C           float64 `json:"c"`
	X0          float64 `json:"x0"`
	Y0          float64 `json:"y0"`
	RightBorder float64 `json:"right_border"`
	H           float64 `json:"h"`
}

func (m MilnEquasion) Validate() error {
	if m.H <= 0 {
		return errors.New("invalid H")
	}
	if m.RightBorder < m.X0 {
		return errors.New("invalid RightBorder")
	}
	return nil
}

type Res struct {
	Picture string      `json:"picture"`
	XYs     plotter.XYs `json:"xy"`
}

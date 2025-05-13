package main

import (
	miln "MethodMilna/internal"
	"MethodMilna/internal/middleware"
	"bytes"
	"encoding/base64"
	"log"
	"math"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func RungeKutta(x, y, h float64, derrivative func(x, y float64) float64) float64 {
	k0 := h * derrivative(x, y)
	k1 := h * derrivative(x+h/2, y+k0/2)
	k2 := h * derrivative(x+h/2, y+k1/2)
	k3 := h * derrivative(x+h, y+k2)
	return y + (k0+2*k1+2*k2+k3)/6
}

func milne(kx, ky, c, h, x0, y0, border float64) ([]float64, []float64) {
	y := make([]float64, 4)
	derrive := func(x, y float64) float64 {
		return kx*x + y*ky + c
	}
	f := make([]float64, 4)
	n := int(math.Ceil((border - x0) / h))
	h = math.Abs(h)
	x := make([]float64, n+1)
	x[0] = x0
	for i := 1; i <= n; i++ {
		x[i] = x[i-1] + h
	}
	y[0] = y0
	f[0] = kx*x0 + ky*y0 + c
	for i := 1; i < len(y) && i < len(x); i++ {
		x[i] = x[i-1] + h
		y[i] = RungeKutta(x[i-1], y[i-1], h, derrive)
		f[i] = derrive(x[i], y[i])
	}

	for i := len(y); i < len(x); i++ {
		y = append(y, y[i-4]+h*4/3*(2*f[i-3]-f[i-2]+2*f[i-1]))
		f = append(f, derrive(x[i], y[i]))
	}

	return x, y
}

func Index(c *fiber.Ctx) error {
	return c.SendFile("./index.html")
}

func convert(x, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))
	for i := 0; i < len(x); i++ {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	return pts
}

func b2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func h(c *fiber.Ctx) error {
	eq := &miln.MilnEquasion{}
	if err := c.BodyParser(eq); err != nil {
		log.Println(err)
		return err
	}
	if err := eq.Validate(); err != nil {
		log.Println(err)
		return err
	}
	p := plot.New()
	xy := convert(milne(eq.Kx, eq.Ky, eq.C, eq.H, eq.X0, eq.Y0, eq.RightBorder))
	err := plotutil.AddLines(p, xy)
	if err != nil {
		log.Println(err)
		return err
	}

	writer, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
	if err != nil {
		log.Println(err)
		return err
	}
	buf := bytes.Buffer{}
	_, err = writer.WriteTo(&buf)
	if err != nil {
		log.Println(err)
		return err
	}
	encoded := make([]byte, base64.RawURLEncoding.EncodedLen(buf.Len()))
	base64.RawStdEncoding.Encode(encoded, buf.Bytes())

	return c.JSON(&miln.Res{
		Picture: b2s(encoded),
		XYs:     xy,
	})
}

func main() {

	// p.X.Label.Text = "X"
	// p.Y.Label.Text = "Y"

	app := fiber.New()
	app.Use(middleware.Logger())
	app.Get("/", Index)
	app.Get("./index.html", Index)
	app.Post("/solve", h)
	app.Listen(":8080")
}

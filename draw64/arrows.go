package draw64

import (
	"image"
	"image/color"
	"math"

	data "github.com/seeder-research/uMagNUS/data64"
	raster "github.com/seeder-research/uMagNUS/freetype/raster"
)

func drawArrows(img *image.RGBA, arr [3][][][]float64, sub int) {
	c := NewCanvas(img)

	Na := data.SizeOf(arr[0]) // number of arrows
	h := Na[Y]                // orignal image height
	Na[X] = imax(Na[X]/sub, 1)
	Na[Y] = imax(Na[Y]/sub, 1)
	Na[Z] = 1
	small := data.Downsample(arr[:], Na)
	S := float64(sub)

	for iy := 0; iy < Na[Y]; iy++ {
		Ay := float64(h) - (float64(iy)+0.5)*S
		for ix := 0; ix < Na[X]; ix++ {
			Ax := (float64(ix) + 0.5) * S
			mx := small[X][0][iy][ix]
			my := small[Y][0][iy][ix]
			mz := small[Z][0][iy][ix]
			c.Arrow(Ax, Ay, mx, my, mz, float64(sub))

		}
	}

	c.rasterizer.Rasterize(c.RGBAPainter)
	c.rasterizer.Clear()
}

// A Canvas is used to draw on.
type Canvas struct {
	*image.RGBA
	*raster.RGBAPainter
	rasterizer *raster.Rasterizer
}

// Make a new canvas of size w x h.
func NewCanvas(img *image.RGBA) *Canvas {
	c := new(Canvas)
	c.RGBA = img
	c.RGBAPainter = raster.NewRGBAPainter(c.RGBA)
	c.rasterizer = raster.NewRasterizer(img.Bounds().Max.X, img.Bounds().Max.Y)
	c.rasterizer.UseNonZeroWinding = true
	c.SetColor(color.RGBA{0, 0, 0, 100})
	return c
}

func (c *Canvas) Arrow(x, y, mx, my, mz, size float64) {

	arrlen := 0.4 * size
	arrw := 0.2 * size

	norm := float64(math.Sqrt(float64(mx*mx + my*my + mz*mz)))
	if norm == 0 {
		return
	}
	if norm > 1 {
		norm = 1
	}

	theta := math.Atan2(float64(my), float64(mx))
	cos := float64(math.Cos(theta))
	sin := float64(math.Sin(theta))
	r1 := arrlen * norm * float64(math.Cos(math.Asin(float64(mz))))
	r2 := arrw * norm

	pt1 := pt((r1*cos)+x, -(r1*sin)+y)
	pt2 := pt((r2*sin-r1*cos)+x, -(-r2*cos-r1*sin)+y)
	pt3 := pt((-r2*sin-r1*cos)+x, -(r2*cos-r1*sin)+y)

	var path raster.Path
	path.Start(pt1)
	path.Add1(pt2)
	path.Add1(pt3)
	path.Add1(pt1)

	c.rasterizer.AddPath(path)
}

func pt(x, y float64) raster.Point {
	return raster.Point{fix32(x), fix32(y)}
}

func fix32(x float64) raster.Fix32 {
	return raster.Fix32(int(x * (1 << 8)))
}

func imax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

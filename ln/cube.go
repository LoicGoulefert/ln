package ln

import (
	"math"
	"math/rand"
)

type Cube struct {
	Min Vector
	Max Vector
	Box Box
}

func NewCube(min, max Vector) *Cube {
	box := Box{min, max}
	return &Cube{min, max, box}
}

func (c *Cube) Compile() {
}

func (c *Cube) BoundingBox() Box {
	return c.Box
}

func (c *Cube) Contains(v Vector, f float64) bool {
	if v.X < c.Min.X-f || v.X > c.Max.X+f {
		return false
	}
	if v.Y < c.Min.Y-f || v.Y > c.Max.Y+f {
		return false
	}
	if v.Z < c.Min.Z-f || v.Z > c.Max.Z+f {
		return false
	}
	return true
}

func (c *Cube) Intersect(r Ray) Hit {
	n := c.Min.Sub(r.Origin).Div(r.Direction)
	f := c.Max.Sub(r.Origin).Div(r.Direction)
	n, f = n.Min(f), n.Max(f)
	t0 := math.Max(math.Max(n.X, n.Y), n.Z)
	t1 := math.Min(math.Min(f.X, f.Y), f.Z)
	if t0 < 1e-3 && t1 > 1e-3 {
		return Hit{c, t1}
	}
	if t0 >= 1e-3 && t0 < t1 {
		return Hit{c, t0}
	}
	return NoHit
}

func (c *Cube) Paths() Paths {
	x1, y1, z1 := c.Min.X, c.Min.Y, c.Min.Z
	x2, y2, z2 := c.Max.X, c.Max.Y, c.Max.Z
	paths := Paths{
		{{x1, y1, z1}, {x1, y1, z2}},
		{{x1, y1, z1}, {x1, y2, z1}},
		{{x1, y1, z1}, {x2, y1, z1}},
		{{x1, y1, z2}, {x1, y2, z2}},
		{{x1, y1, z2}, {x2, y1, z2}},
		{{x1, y2, z1}, {x1, y2, z2}},
		{{x1, y2, z1}, {x2, y2, z1}},
		{{x1, y2, z2}, {x2, y2, z2}},
		{{x2, y1, z1}, {x2, y1, z2}},
		{{x2, y1, z1}, {x2, y2, z1}},
		{{x2, y1, z2}, {x2, y2, z2}},
		{{x2, y2, z1}, {x2, y2, z2}},
	}
	return paths
}

type StripedCube struct {
	Cube
	StripesX, StripesY, StripesZ int
	PercentX, PercentY, PercentZ float64
}

func NewStripedCube(min, max Vector, stripesX, stripesY, stripesZ int, px, py, pz float64) *StripedCube {
	cube := NewCube(min, max)
	return &StripedCube{*cube, stripesX, stripesY, stripesZ, px, py, pz}
}

func (c *StripedCube) Paths() Paths {
	var paths Paths
	x1, y1, z1 := c.Min.X, c.Min.Y, c.Min.Z
	x2, y2, z2 := c.Max.X, c.Max.Y, c.Max.Z
	xLen := x2 - x1
	yLen := y2 - y1
	zLen := z2 - z1

	// along x
	for i := 1; i < c.StripesX; i++ {
		p := Remap(float64(i), 0, float64(c.StripesX), 0, 1)
		ry := rand.Float64() * (y2 - y1) * c.PercentY
		rz := rand.Float64() * (z2 - z1) * c.PercentZ
		px := p*xLen + x1

		// y1, z1
		paths = append(paths, Path{{px, y1 + ry, z1}, {px, y1, z1}, {px, y1, z1 + rz}})
		// y1, z2
		paths = append(paths, Path{{px, y1, z2 - rz}, {px, y1, z2}, {px, y1 + ry, z2}})
		// y2, z2
		paths = append(paths, Path{{px, y2 - ry, z2}, {px, y2, z2}, {px, y2, z2 - rz}})
		// y2, z1
		paths = append(paths, Path{{px, y2, z1 + rz}, {px, y2, z1}, {px, y2 - ry, z1}})
	}

	// along y
	for i := 1; i < c.StripesY; i++ {
		p := Remap(float64(i), 0, float64(c.StripesY), 0, 1)
		rx := rand.Float64() * (x2 - x1) * c.PercentX
		rz := rand.Float64() * (z2 - z1) * c.PercentZ
		py := p*yLen + y1

		paths = append(paths, Path{{x1 + rx, py, z1}, {x1, py, z1}, {x1, py, z1 + rz}})
		paths = append(paths, Path{{x1, py, z2 - rz}, {x1, py, z2}, {x1 + rx, py, z2}})
		paths = append(paths, Path{{x2 - rx, py, z2}, {x2, py, z2}, {x2, py, z2 - rz}})
		paths = append(paths, Path{{x2, py, z1 + rz}, {x2, py, z1}, {x2 - rx, py, z1}})
	}

	// along z
	for i := 1; i < c.StripesZ; i++ {
		p := Remap(float64(i), 0, float64(c.StripesZ), 0, 1)
		rx := rand.Float64() * (x2 - x1) * c.PercentX
		ry := rand.Float64() * (y2 - y1) * c.PercentY
		pz := p*zLen + z1

		paths = append(paths, Path{{x1 + rx, y1, pz}, {x1, y1, pz}, {x1, y1 + ry, pz}})
		paths = append(paths, Path{{x1, y2 - ry, pz}, {x1, y2, pz}, {x1 + rx, y2, pz}})
		paths = append(paths, Path{{x2 - rx, y2, pz}, {x2, y2, pz}, {x2, y2 - ry, pz}})
		paths = append(paths, Path{{x2, y1 + ry, pz}, {x2, y1, pz}, {x2 - rx, y1, pz}})
	}

	// Add original cube paths
	normalCube := NewCube(c.Min, c.Max)
	paths = append(paths, normalCube.Paths()...)

	return paths
}

package ln

type Shard struct {
	mesh                     *Mesh
	p1, p2, p3, p4, ph1, ph2 Vector
}

func NewShard(a, b, h1, h2 float64) *Shard {
	p1 := Vector{a / 2, -b / 2, 0}
	p2 := Vector{a / 2, b / 2, 0}
	p3 := Vector{-a / 2, b / 2, 0}
	p4 := Vector{-a / 2, -b / 2, 0}
	ph1 := Vector{0, 0, h1}
	ph2 := Vector{0, 0, -h2}

	// Top pyramid 4 faces
	t1 := NewTriangle(p1, p2, ph1)
	t2 := NewTriangle(p2, p3, ph1)
	t3 := NewTriangle(p3, p4, ph1)
	t4 := NewTriangle(p4, p1, ph1)

	// Bottom pyramid 4 faces
	t5 := NewTriangle(p1, p2, ph2)
	t6 := NewTriangle(p2, p3, ph2)
	t7 := NewTriangle(p3, p4, ph2)
	t8 := NewTriangle(p4, p1, ph2)

	triangles := []*Triangle{t1, t2, t3, t4, t5, t6, t7, t8}
	mesh := NewMesh(triangles)

	return &Shard{mesh, p1, p2, p3, p4, ph1, ph2}
}

func (s *Shard) BoundingBox() Box {
	return s.mesh.Box
}

func (s *Shard) Compile() {
	s.mesh.Compile()
}

func (s *Shard) Contains(v Vector, f float64) bool {
	return false
}

func (s *Shard) Intersect(r Ray) Hit {
	return s.mesh.Intersect(r)
}

func (s *Shard) Paths() Paths {
	paths := s.mesh.Paths()
	return paths
}

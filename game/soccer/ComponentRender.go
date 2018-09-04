package soccer

type Render struct {
	type_         string
	static        bool
	DebugPoints   [][2]float64
	DebugSegments [][2][2]float64
}

func (r Render) GetType() string {
	return r.type_
}

package message

const (
	// KindConnected is sent when user connects
	KindConnected = iota + 1
	// KindStroke message specifies a drawn stroke by a user
	KindStroke
	// KindClear message is sent when a user clears the screen
	KindClear
	KindUndo
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Stroke struct {
	Points   []Point `json:"points"`
	Color    string  `json:"color"`
	OwnerID  int     `json:"ownerID"`
	StrokeID int     `json:"strokeID"`
	ShapeID  int     `json:"shapeID"`
}

type Message struct {
	Kind   int    `json:"kind,omitempty"`
	Stroke Stroke `json:"stroke,omitempty"`
}

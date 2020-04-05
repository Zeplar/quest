package message

const (
	// KindConnected is sent when user connects
	KindConnected = iota + 1
	// KindUserJoined is sent when someone else joins
	KindUserJoined
	// KindUserLeft is sent when someone leaves
	KindUserLeft
	// KindStroke message specifies a drawn stroke by a user
	KindStroke
	// KindClear message is sent when a user clears the screen
	KindClear
)

type Point struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

type User struct {
	ID    string `json:"id,omitempty"`
	Color string `json:"color,omitempty"`
}

type Message struct {
	Kind   int     `json:"kind,omitempty"`
	User   User    `json:"user,omitempty"`
	Users  []User  `json:"users,omitempty"`
	Points []Point `json:"points,omitempty"`
	Finish bool    `json:"finish,omitempty"`
}

func NewConnected(color string, users []User) *Message {
	return &Message{
		Kind:  KindConnected,
		User:  User{Color: color},
		Users: users,
	}
}

func NewUserJoined(userID string, color string) *Message {
	return &Message{
		Kind: KindUserJoined,
		User: User{ID: userID, Color: color},
	}
}

func NewUserLeft(userID string) *Message {
	return &Message{
		Kind: KindUserLeft,
		User: User{ID: userID},
	}
}

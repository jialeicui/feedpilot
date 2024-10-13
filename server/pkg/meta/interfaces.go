package meta

type Meta interface {
}

type Media interface {
	GetID() MediaID
	GetType() MediaType
	GetContent() []byte
}

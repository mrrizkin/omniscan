package encoder

type TextEncoding interface {
	Decode(raw string) string
}

type NoOpEncoder struct{}

func (e *NoOpEncoder) Decode(raw string) string {
	return raw
}

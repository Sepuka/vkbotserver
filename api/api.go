package api

const (
	defaultOutput = `ok`
)

func Response() []byte {
	return []byte(defaultOutput)
}

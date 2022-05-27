package transport

type server struct {
}

func NewServer() (*server, error) {
	return &server{}, nil
}

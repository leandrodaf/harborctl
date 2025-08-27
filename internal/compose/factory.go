package compose

// NewDefaultService cria um serviço usando o gerador com micro-interfaces
func NewDefaultService() Service {
	return NewService(NewMicroGenerator())
}

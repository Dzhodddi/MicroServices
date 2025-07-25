package shared_errors

import "fmt"

var (
	ValidationError        = fmt.Errorf("validation error")
	BadRequestPayload      = fmt.Errorf("bad request payload")
	ServerError            = fmt.Errorf("server error")
	ViolatePK              = fmt.Errorf("user with this credantionals already exists")
	EstablishingConnection = fmt.Errorf("error during connection")
	QueueDeclareError      = fmt.Errorf("error declaring queue")
	NotFoundError          = fmt.Errorf("resource not found")
	InvalidToken           = fmt.Errorf("invalid token")
)

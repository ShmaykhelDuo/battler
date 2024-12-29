package api

import "fmt"

type Kind struct {
	ID          int
	Description string
}

var KindInternal = Kind{
	ID:          0,
	Description: "internal",
}

var KindInvalidRequest = Kind{
	ID:          1,
	Description: "invalid request",
}

var KindInvalidArgument = Kind{
	ID:          2,
	Description: "invalid argument",
}

var KindNotFound = Kind{
	ID:          3,
	Description: "not found",
}

var KindAlreadyExists = Kind{
	ID:          4,
	Description: "already exists",
}

var KindInvalidCredentials = Kind{
	ID:          5,
	Description: "invalid username or password",
}

var KindUnauthenticated = Kind{
	ID:          6,
	Description: "unauthenticated",
}

type Error struct {
	Kind    Kind
	Message string
}

func (e Error) Error() string {
	if e.Message == "" {
		return string(e.Kind.Description)
	}

	return fmt.Sprintf("%s: %s", e.Kind.Description, e.Message)
}

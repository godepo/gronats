package gronats

import (
	"errors"

	"github.com/google/uuid"
)

func UnexpectError() error {
	return errors.New(uuid.NewString())
}

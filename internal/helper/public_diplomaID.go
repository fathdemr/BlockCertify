package helper

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

func GenerateDiplomaPublicIDFromUUID(id uuid.UUID) string {
	short := strings.ToUpper(strings.ReplaceAll(id.String()[12:], "-", ""))
	return fmt.Sprintf("BC-%s", short)
}

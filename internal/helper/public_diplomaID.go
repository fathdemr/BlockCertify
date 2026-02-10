package helper

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

func GenerateDiplomaPublicIDFromUUID(id uuid.UUID) string {
	short := strings.ToUpper(strings.ReplaceAll(id.String(), "-", ""))
	last12 := short[len(short)-12:]
	return fmt.Sprintf("BC-%s", strings.ToUpper(last12))
}

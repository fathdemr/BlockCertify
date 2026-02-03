package helper

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateDiplomaPublicIDFromUUID(id uuid.UUID) string {
	short := strings.ToUpper(strings.ReplaceAll(id.String()[:8], "-", ""))
	return fmt.Sprintf("BC-%s", short)
}

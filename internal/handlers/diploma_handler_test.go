package handlers

import (
	"BlockCertify/internal/dto"
	"testing"
	"time"
)

func validMeta() dto.DiplomaMetadataRequest {
	return dto.DiplomaMetadataRequest{
		FirstName:      "Fatih",
		LastName:       "Demir",
		Email:          "fatih@example.com",
		University:     "Test University",
		Faculty:        "Engineering",
		Department:     "Computer Engineering",
		GraduationYear: time.Now().Year(),
		StudentNumber:  "12345",
		Nationality:    "TR",
	}
}

func TestValidateUploadMetadataValid(t *testing.T) {
	if err := validateUploadMetadata(validMeta()); err != nil {
		t.Errorf("valid metadata rejected: %v", err)
	}
}

func TestValidateUploadMetadataMissingFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*dto.DiplomaMetadataRequest)
	}{
		{"missing firstName", func(m *dto.DiplomaMetadataRequest) { m.FirstName = " " }},
		{"missing lastName", func(m *dto.DiplomaMetadataRequest) { m.LastName = "" }},
		{"missing email", func(m *dto.DiplomaMetadataRequest) { m.Email = "" }},
		{"missing university", func(m *dto.DiplomaMetadataRequest) { m.University = "" }},
		{"missing department", func(m *dto.DiplomaMetadataRequest) { m.Department = "" }},
		{"year too old", func(m *dto.DiplomaMetadataRequest) { m.GraduationYear = 1900 }},
		{"year in future", func(m *dto.DiplomaMetadataRequest) { m.GraduationYear = time.Now().Year() + 5 }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			meta := validMeta()
			tc.mutate(&meta)
			if err := validateUploadMetadata(meta); err == nil {
				t.Errorf("%s: invalid metadata accepted", tc.name)
			}
		})
	}
}

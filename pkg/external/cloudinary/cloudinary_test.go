package cloudinary

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGenerateSignature(t *testing.T) {
	t.Parallel()

	timestamp, _ := time.Parse("20060102150405", "20060102150405")
	signature := generateSignature("PUBLICID", timestamp.String(), false, "SECRETKEY")

	expected := "f96dab52fb7291dcd162d5bd7bd47f33454afe3a"
	if signature != expected {
		t.Errorf("Error generating Signature, got %s, expected: %s", signature, expected)
	}
}

func TestGeneratePublicID(t *testing.T) {
	t.Parallel()

	game := "CALL_OF_DUTY"
	teamName := "Valiance&Co"
	publicID := generatePublicIDForLogo(game, teamName)
	expected := "estars/logos/CALL_OF_DUTY/VALIANCE-CO"
	if publicID != expected {
		t.Errorf("Error generating Public ID, got %s, expected: %s", publicID, expected)
	}

	game = "LEAGUE_OF_LEGENDS"
	teamName = "Test Team 123"
	publicID = generatePublicIDForLogo(game, teamName)
	expected = "estars/logos/LEAGUE_OF_LEGENDS/TEST_TEAM_123"
	if publicID != expected {
		t.Errorf("Error generating Public ID, got %s, expected: %s", publicID, expected)
	}
}

func setup() {}

func teardown() {}

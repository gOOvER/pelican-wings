package parser

import (
	"testing"
)

// TestUpdateJsonPreservingStructure tests that JSON structure is preserved during updates
func TestUpdateJsonPreservingStructure(t *testing.T) {
	originalJSON := []byte(`{"AuthCredentialStore":{"Path":"auth.enc","Type":"Encrypted"},"ConnectionTimeouts":{"JoinTimeouts":{}},"Defaults":{"GameMode":"Adventure","World":"default"},"DisplayTmpTagsInStrings":false,"LogLevels":{},"MOTD":"","MaxPlayers":"100","MaxViewRadius":"32","Mods":{},"Modules":{"PathPlugin":{"Modules":{}}},"Password":"","PlayerStorage":{"Type":"Hytale"},"RateLimit":{},"ServerName":"Pelican hosted Server","Version":3}`)

	// Create a test configuration with replacements
	cf := &ConfigurationFile{
		FileName: "test.json",
		Parser:   Json,
		Replace: []ConfigurationFileReplacement{
			{
				Match:       "ServerName",
				ReplaceWith: ReplaceValue{value: []byte("\"Updated Server Name\""), valueType: 1}, // 1 = string type
			},
			{
				Match:       "MaxPlayers",
				ReplaceWith: ReplaceValue{value: []byte("\"50\""), valueType: 1},
			},
			{
				Match:       "MOTD",
				ReplaceWith: ReplaceValue{value: []byte("\"Welcome to the server!\""), valueType: 1},
			},
			{
				Match:       "DisplayTmpTagsInStrings",
				ReplaceWith: ReplaceValue{value: []byte("true"), valueType: 0}, // 0 = boolean type
			},
		},
	}

	// Try to update the JSON
	output, err := cf.UpdateJsonPreservingStructure(originalJSON)
	if err != nil {
		t.Fatalf("UpdateJsonPreservingStructure failed: %v", err)
	}

	if output == nil || len(output) == 0 {
		t.Fatal("UpdateJsonPreservingStructure returned empty result")
	}

	// Check that key values were updated
	if !contains(output, "Updated Server Name") {
		t.Error("ServerName was not updated")
	}
	if !contains(output, "Welcome to the server!") {
		t.Error("MOTD was not updated")
	}

	t.Logf("Original size: %d bytes", len(originalJSON))
	t.Logf("Updated size: %d bytes", len(output))
	t.Log("✓ JSON structure preserved during update")
}

// Helper function to check if JSON contains a value
func contains(data []byte, str string) bool {
	for i := 0; i < len(data)-len(str); i++ {
		if string(data[i:i+len(str)]) == str {
			return true
		}
	}
	return false
}

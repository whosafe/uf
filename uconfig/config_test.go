package uconfig

import (
	"testing"
)

// Mock items
type TestConfig struct {
	Name    string
	Version string
	Port    int
}

func (t *TestConfig) UnmarshalYAML(key string, value *Node) error {
	switch key {
	case "name":
		t.Name = value.Value
	case "version":
		t.Version = value.Value
	case "port":
		// simple int conversation logic for test, real usage uses uconv
		if value.Value == "8080" {
			t.Port = 8080
		}
	}
	return nil
}

func TestParseConfig(t *testing.T) {
	yamlData := []byte(`
app:
  name: "MyApp"
  version: v1
  port: 8080
list:
  - a
  - b
unknown:
  foo: bar
`)

	var appCfg TestConfig

	// Register
	Register("app", appCfg.UnmarshalYAML)

	var listData []string
	Register("list", func(key string, value *Node) error {
		return value.Iter(func(i int, v *Node) error {
			listData = append(listData, v.Value)
			return nil
		})
	})

	var unknownKey string
	var unknownVal string
	RegisterUnknown(func(key string, value *Node) error {
		unknownKey = key
		if value.Kind == MappingNode {
			unknownVal = "map" // simplified check
		}
		return nil
	})

	if err := ParseConfig(yamlData); err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	// Verify App
	if appCfg.Name != "MyApp" {
		t.Errorf("expected Name MyApp, got %s", appCfg.Name)
	}
	if appCfg.Version != "v1" {
		t.Errorf("expected Version v1, got %s", appCfg.Version)
	}
	if appCfg.Port != 8080 {
		t.Errorf("expected Port 8080, got %d", appCfg.Port)
	}

	// Verify List
	if len(listData) != 2 || listData[0] != "a" || listData[1] != "b" {
		t.Errorf("expected list [a, b], got %v", listData)
	}

	// Verify Unknown
	if unknownKey != "unknown" || unknownVal != "map" {
		t.Errorf("expected unknown key 'unknown' and val 'map', got %s, %s", unknownKey, unknownVal)
	}
}

func TestCallback(t *testing.T) {
	// Re-parse to set rootNode for this test
	yamlData := []byte(`
app:
  name: "CallbackApp"
`)
	if err := ParseConfig(yamlData); err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	var name string
	// Manually trigger callback for "app" section
	// Note: since "app" is a map, callback is called for each child
	err := Callback("app", func(key string, value *Node) error {
		if key == "name" {
			name = value.Value
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Callback failed: %v", err)
	}

	if name != "CallbackApp" {
		t.Errorf("expected name CallbackApp, got %s", name)
	}
}

func TestDecode(t *testing.T) {
	// Construct a node manually to test Decode
	root := &Node{
		Kind: MappingNode,
		Children: map[string]*Node{
			"name": {Kind: ScalarNode, Value: "Test"},
		},
	}

	var cfg TestConfig
	if err := root.Decode(&cfg); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if cfg.Name != "Test" {
		t.Errorf("expected Name Test, got %s", cfg.Name)
	}
}

func TestFailures(t *testing.T) {
	// Test wrong root type
	// This requires tricking ParseConfig logic or using invalid YAML that Parse accepts but is not map?
	// Our Parse always returns Map as root?
	// Yes, `root := &Node{Kind: MappingNode}` in parser.go.
	// So root is always Map unless parsing fails.

	// Test Decode on non-map
	scalar := &Node{Kind: ScalarNode, Value: "val"}
	var cfg TestConfig
	if err := scalar.Decode(&cfg); err == nil {
		t.Error("expected error decoding scalar to struct, got nil")
	}

	// Test Iter on non-list
	if err := scalar.Iter(func(i int, v *Node) error { return nil }); err == nil {
		t.Error("expected error iterating scalar, got nil")
	}
}

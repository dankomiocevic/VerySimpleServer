package config

import (
	"testing"
)

func TestConfigureSlot(t *testing.T) {
	config := make(map[interface{}]interface{})
	slot_one := make(map[string]interface{})

	slot_one["kind"] = "simple_memory"
	config["slot_000"] = slot_one

	response := ConfigureSlots(config)

	if response[0] == nil {
		t.Fatalf("slot zero not configured: %s", response)
	}
}

func TestConfigureTimeoutSlot(t *testing.T) {
	config := make(map[interface{}]interface{})
	slot_one := make(map[string]interface{})

	slot_one["kind"] = "timeout_memory"
	slot_one["timeout"] = 50
	config["slot_000"] = slot_one

	response := ConfigureSlots(config)

	if response[0] == nil {
		t.Fatalf("slot zero not configured: %s", response)
	}
}

func TestNotConfigureSlot(t *testing.T) {
	config := make(map[interface{}]interface{})
	slot_one := make(map[string]interface{})

	slot_one["kind"] = "simple_memory"
	config["slot_000"] = slot_one

	response := ConfigureSlots(config)

	for i := 1; i < 1000; i++ {
		if response[i] != nil {
			t.Fatalf("slot %d should not be configured: %s", i, response)
		}
	}
}

func TestConfigureUnknownType(t *testing.T) {
	config := make(map[interface{}]interface{})
	slot_one := make(map[string]interface{})

	slot_one["kind"] = "unknown"
	config["slot_000"] = slot_one

	response := ConfigureSlots(config)

	if response[0] != nil {
		t.Fatalf("slot zero should not be configured: %s", response)
	}
}

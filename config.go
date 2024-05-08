package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func loadConfig(filename string) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})

	// Open config file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func configureSlots(conf map[interface{}]interface{}) [1000]Slot {
	slots := [1000]Slot{}

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("slot_%03d", i)
		value, ok := conf[key]
		if ok {
			// Assert this is a map
			valueMap, ok := value.(map[string]interface{})
			if ok {
				slot, _ := getSlot(valueMap)
				slots[i] = slot
			}
		}
	}

	return slots
}

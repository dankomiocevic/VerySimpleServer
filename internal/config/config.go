package config

import (
	"fmt"
	"os"

	"github.com/dankomiocevic/VerySimpleServer/internal/slots"
	"gopkg.in/yaml.v3"
)

func LoadConfig(filename string) (map[interface{}]interface{}, error) {
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

func ConfigureSlots(conf map[interface{}]interface{}) [1000]slots.Slot {
	slotsArray := [1000]slots.Slot{}

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("slot_%03d", i)
		value, ok := conf[key]
		if ok {
			// Assert this is a map
			valueMap, ok := value.(map[string]interface{})
			if ok {
				slot, _ := slots.GetSlot(valueMap)
				slotsArray[i] = slot
			}
		}
	}

	return slotsArray
}

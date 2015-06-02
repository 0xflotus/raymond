package raymond

import "reflect"

// Cf. private variables at: http://handlebarsjs.com/block_helpers.html

// A private data frame
type DataFrame struct {
	parent *DataFrame
	data   map[string]interface{}
}

// Instanciate a new private data frame
func NewDataFrame() *DataFrame {
	return &DataFrame{
		data: make(map[string]interface{}),
	}
}

// Returns a new private data frame, with parent set to self
func (p *DataFrame) Copy() *DataFrame {
	result := NewDataFrame()

	for k, v := range p.data {
		result.data[k] = v
	}

	result.parent = p

	return result
}

func (p *DataFrame) NewIterDataFrame(length int, i int, key interface{}) *DataFrame {
	result := p.Copy()

	result.Set("index", i)
	result.Set("key", key)
	result.Set("first", i == 0)
	result.Set("last", i == length-1)

	return result
}

// Set a data value
func (p *DataFrame) Set(key string, val interface{}) {
	p.data[key] = val
}

// Get a data value
func (p *DataFrame) Get(key string) interface{} {
	return p.Find([]string{key})
}

// Get a deep data value
func (p *DataFrame) Find(parts []string) interface{} {
	data := p.data

	for i, part := range parts {
		val := data[part]
		if val == nil {
			return nil
		}

		if i == len(parts)-1 {
			// found
			return val
		}

		valValue := reflect.ValueOf(val)
		if valValue.Kind() != reflect.Map {
			// not found
			return nil
		}

		// continue
		data = mapStringInterface(valValue)
	}

	// not found
	return nil
}

// converts any `map` to `map[string]interface{}`
func mapStringInterface(value reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range value.MapKeys() {
		result[strValue(key)] = value.MapIndex(key).Interface()
	}

	return result
}

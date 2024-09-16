package serial

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

type OrderedMap struct {
	Values []*OrderedValue
}

type OrderedValue struct {
	Key   string
	Value any
}

type OrderedArray []any

// MarshalJSON implements the json.Marshaler interface.
func (om OrderedMap) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteRune('{')
	for i, val := range om.Values {
		// write key
		b, e := json.Marshal(val.Key)
		if e != nil {
			return nil, e
		}
		buf.Write(b)
		// write delimiter
		buf.WriteRune(':')
		// write value
		b, e = json.Marshal(val.Value)
		if e != nil {
			return nil, e
		}
		buf.Write(b)
		// write delimiter
		if i+1 < len(om.Values) {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()

	t, err := d.Token()
	if err == io.EOF {
		return nil
	}

	if t != json.Delim('{') {
		return errors.New("unexpected start of object")
	}
	return om.unmarshalEmbededObject(d)
}

// Get return an OrderedValue from OrderedMap.
func (om *OrderedMap) Get(key string) (*OrderedValue, bool) {
	for i := range om.Values {
		if om.Values[i].Key == key {
			return om.Values[i], true
		}
	}

	return nil, false
}

// Set change or add an Object to OrderedMap.
func (om *OrderedMap) Set(key string, value any) {

	for i := range om.Values {
		if om.Values[i].Key == key {
			om.Values[i].Value = value
			return
		}
	}
	om.Values = append(om.Values, &OrderedValue{key, value})
}

// SetValue change or add an OrderedValue to OrderedMap.
func (om *OrderedMap) SetValue(v *OrderedValue) {
	for i := range om.Values {
		if om.Values[i].Key == v.Key {
			om.Values[i] = v
			return
		}
	}
	om.Values = append(om.Values, v)
}

// Delete remove an Object from OrderedMap.
func (om *OrderedMap) Delete(key string) {
	for i, val := range om.Values {
		if val.Key == key {
			om.Values = append(om.Values[:i], om.Values[i+1:]...)
		}
	}
}

// DeleteValue remove an OrderedValue from OrderedMap.
func (om *OrderedMap) DeleteValue(v *OrderedValue) {
	for i, val := range om.Values {
		if val.Key == v.Key {
			om.Values = append(om.Values[:i], om.Values[i+1:]...)
		}
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (arr OrderedArray) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteRune('[')
	for i, val := range arr {
		// write key
		b, e := json.Marshal(val)
		if e != nil {
			return nil, e
		}
		buf.Write(b)
		// write delimiter
		if i+1 < len(arr) {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune(']')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (arr *OrderedArray) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	t, err := d.Token()
	if err == io.EOF {
		return nil
	}
	if t != json.Delim('[') {
		return errors.New("unexpected start of array")
	}
	return arr.unmarshalEmbededArray(d)
}

func (om *OrderedMap) unmarshalEmbededObject(d *json.Decoder) error {
	for d.More() {
		kToken, err := d.Token()
		if err == io.EOF || (err == nil && kToken == json.Delim('}')) {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		vToken, err := d.Token()
		if err == io.EOF {
			// log.Print("unexpected EOF")
			return errors.New("unexpected EOF")
		}

		var val any
		switch vToken {
		case json.Delim('{'):
			var obj OrderedMap
			if err = obj.unmarshalEmbededObject(d); err != nil {
				return err
			}
			val = obj
		case json.Delim('['):
			var arr OrderedArray
			err = arr.unmarshalEmbededArray(d)
			val = arr
		default:
			val = vToken
		}

		if err != nil {
			return err
		}

		om.Values = append(om.Values, &OrderedValue{kToken.(string), val})
	}

	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim('}') {
		return errors.New("unexpected EOF")
	}
	return err
}

func (arr *OrderedArray) unmarshalEmbededArray(d *json.Decoder) error {
	for d.More() {
		token, err := d.Token()
		if err == io.EOF || (err == nil && token == json.Delim(']')) {
			return errors.New("unexpected EOF")
		}
		var val any
		switch token {
		case json.Delim('{'):
			var obj OrderedMap
			if err = obj.unmarshalEmbededObject(d); err != nil {
				return err
			}
			val = obj
		case json.Delim('['):
			var arr OrderedArray
			err = arr.unmarshalEmbededArray(d)
			val = arr
		default:
			val = token
		}
		if err != nil {
			return err
		}
		*arr = append(*arr, val)
	}

	kToken, err := d.Token()
	if err == io.EOF || kToken != json.Delim(']') {
		return errors.New("unexpected EOF")
	}

	if *arr == nil {
		*arr = make([]any, 0)
	}

	return nil
}

// MarshalYAML implements the yaml.MarshalYAML interface.
func (om OrderedMap) MarshalYAML() (any, error) {
	node := yaml.Node{
		Kind: yaml.MappingNode,
	}

	for _, val := range om.Values {
		key, value := val.Key, val.Value
		keyNode := &yaml.Node{}

		// serialize key to yaml, then deserialize it back into the node
		// this is a hack to get the correct tag for the key
		if err := keyNode.Encode(key); err != nil {
			return nil, err
		}

		valueNode := &yaml.Node{}
		if err := valueNode.Encode(value); err != nil {
			return nil, err
		}

		node.Content = append(node.Content, keyNode, valueNode)
	}

	return &node, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (om *OrderedMap) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("pipeline must contain YAML mapping, has %v", value.Kind)
	}

	if om.Values == nil {
		om.Values = make([]*OrderedValue, 0)
	}

	key := ""
	for _, node := range value.Content {
		if len(key) == 0 {
			if node.Kind != yaml.ScalarNode {
				continue
			}
			if err := node.Decode(&key); err != nil {
				return err
			}
		} else {
			switch node.Kind {
			case yaml.MappingNode:
				var obj OrderedMap
				if err := node.Decode(&obj); err != nil {
					return err
				}
				om.Values = append(om.Values, &OrderedValue{key, obj})
			case yaml.SequenceNode:
				var arr OrderedArray
				if err := node.Decode(&arr); err != nil {
					return err
				}
				om.Values = append(om.Values, &OrderedValue{key, arr})
			case yaml.AliasNode, yaml.ScalarNode:
				var ins any
				if err := node.Decode(&ins); err != nil {
					return err
				}
				om.Values = append(om.Values, &OrderedValue{key, ins})
			default:
				continue
			}
			key = ""
		}
	}
	return nil
}

// MarshalYAML implements the yaml.MarshalYAML interface.
func (arr OrderedArray) MarshalYAML() (any, error) {
	node := yaml.Node{
		Kind: yaml.SequenceNode,
	}

	for _, val := range arr {
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(val); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, valueNode)
	}

	return &node, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (arr *OrderedArray) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("pipeline must contain YAML sequence, has %v", value.Kind)
	}

	for _, node := range value.Content {
		switch node.Kind {
		case yaml.MappingNode:
			var obj OrderedMap
			if err := node.Decode(&obj); err != nil {
				return err
			}
			*arr = append(*arr, obj)
		case yaml.SequenceNode:
			var r OrderedArray
			if err := node.Decode(&r); err != nil {
				return err
			}
			*arr = append(*arr, r)
		case yaml.AliasNode, yaml.ScalarNode:
			var ins any
			if err := node.Decode(&ins); err != nil {
				return err
			}
			*arr = append(*arr, ins)
		default:
			continue
		}
	}
	return nil
}

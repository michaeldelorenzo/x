// Package jsonblob provides a JSON blob type that can be used for marshaling/unmarshaling arbitrary json data from/to
// various mediums
package jsonblob

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/michaeldelorenzo/x/pkg/utils/maputils"
	"go.mongodb.org/mongo-driver/bson"
)

// Blob is a useful type for marshaling/unmarshaling arbitrary json to/from various transport & storage mediums
type Blob map[string]interface{}

// Value implements interface https://pkg.go.dev/database/sql/driver#Valuer.
// Value marshals the Blob serialized as JSON for use in jsonb sql columns.
func (b Blob) Value() (driver.Value, error) {
	bytes, err := json.Marshal(b)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal blob to json: %w", err)
	}

	return bytes, nil
}

// Scan implements the interface https://pkg.go.dev/database/sql#Scanner
// Scan unmarshals a raw value into a jsonblob.Blob
func (b *Blob) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	if b == nil {
		return errors.New("cannot scan into nil *Blob")
	}

	var source []byte

	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	case map[string]interface{}:
		*b = v
		return nil
	default:
		return fmt.Errorf("recieved unsupported type %T when scanning into Blob", src)

	}

	err := json.Unmarshal(source, b)
	if err != nil {
		return err
	}

	return nil
}

// MarshalBSON implements interface https://pkg.go.dev/go.mongodb.org/mongo-driver/bson#Marshaler.
// MarshalBSON returns the Blob serialized as BSON such that the keys are in alphabetical order to remain consistent
// with JSON serialization.
func (b Blob) MarshalBSON() ([]byte, error) {
	return bson.Marshal(BuildOrderedBSONMap(map[string]interface{}(b)))
}

// BuildOrderedBSONMap organizes the input into BSON, deeply ordering object keys alphabetically via bson.D
func BuildOrderedBSONMap(v interface{}) interface{} {
	switch input := v.(type) {
	case map[string]interface{}: // For a map, sort its keys recursively.
		keys := maputils.Keys(input)
		sort.Strings(keys)

		orderedBSONObject := make(bson.D, len(keys))
		for i, k := range keys {
			orderedBSONObject[i] = bson.E{Key: k, Value: BuildOrderedBSONMap(input[k])}
		}

		return orderedBSONObject
	case []interface{}: // For an array, sort each element recursively.
		orderedBSONArray := make(bson.A, len(input))
		for i, v := range input {
			orderedBSONArray[i] = BuildOrderedBSONMap(v)
		}
		return orderedBSONArray
	default: // For any other type, return it as is.
		return input
	}
}

package jsonblob_test

import (
	"encoding/json"
	"testing"

	"github.com/michaeldelorenzo/x/v3/pkg/types/jsonblob"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBlob(t *testing.T) {
	t.Run("marshals and unmarshals blob for json appropriately", func(t *testing.T) {
		for i := 0; i < 20; i += 1 {
			// given
			original := jsonblob.Blob{
				"c": "doerayme",
				"a": 45.67,
				"b": map[string]interface{}{"yes": "hey", "no": "way"},
				"1": []interface{}{
					map[string]interface{}{"z": "oh", "x": "yes"},
					map[string]interface{}{"y": true},
				},
			}
			expected := `{"1":[{"x":"yes","z":"oh"},{"y":true}],"a":45.67,"b":{"no":"way","yes":"hey"},"c":"doerayme"}`
			output := jsonblob.Blob{}

			marshaled, err := json.Marshal(original)

			// then
			require.NoError(t, err)
			require.Equal(t, expected, string(marshaled))
			err = json.Unmarshal(marshaled, &output)
			require.NoError(t, err)
			require.EqualValues(t, original, output)
		}
	})

	t.Run("marshals and unmarshals blob for sql appropriately", func(t *testing.T) {
		for i := 0; i < 20; i += 1 {
			// given
			original := jsonblob.Blob{
				"c": "doerayme",
				"a": 45.67,
				"b": map[string]interface{}{"yes": "hey", "no": "way"},
				"1": []interface{}{
					map[string]interface{}{"z": "oh", "x": "yes"},
					map[string]interface{}{"y": true},
				},
			}
			expected := `{"1":[{"x":"yes","z":"oh"},{"y":true}],"a":45.67,"b":{"no":"way","yes":"hey"},"c":"doerayme"}`
			output := jsonblob.Blob{}

			marshaled, err := original.Value()

			// then
			require.NoError(t, err)
			require.Equal(t, expected, string(marshaled.([]byte)))
			err = output.Scan(marshaled)
			require.NoError(t, err)
			require.EqualValues(t, original, output)
		}
	})

	t.Run("marshals and unmarshals blob for bson appropriately", func(t *testing.T) {
		for i := 0; i < 20; i += 1 {
			// given
			original := jsonblob.Blob{
				"c": "doerayme",
				"a": 45.67,
				"b": map[string]interface{}{"yes": "hey", "no": "way"},
				"1": []interface{}{
					map[string]interface{}{"z": "oh", "x": "yes"},
					map[string]interface{}{"y": true},
				},
			}

			expected, err := bson.Marshal(bson.D{
				{Key: "1", Value: bson.A{
					bson.D{
						{Key: "x", Value: "yes"},
						{Key: "z", Value: "oh"},
					},
					bson.D{{Key: "y", Value: true}},
				}},
				{Key: "a", Value: 45.67},
				{Key: "b", Value: bson.D{{Key: "no", Value: "way"}, {Key: "yes", Value: "hey"}}},
				{Key: "c", Value: "doerayme"},
			})
			require.NoError(t, err)

			marshaled, err := bson.Marshal(original)

			// then
			require.NoError(t, err)
			require.Equal(t, expected, marshaled)
		}
	})

	t.Run("supports scanning valid map into Blob", func(t *testing.T) {
		// given
		original := map[string]interface{}{"abc": "doerayme", "123": 45.67}
		output := jsonblob.Blob{}

		// when
		err := output.Scan(original)

		// then
		require.NoError(t, err)
		require.EqualValues(t, original, output)
	})

	t.Run("supports scanning valid string into Blob", func(t *testing.T) {
		// given
		original := `{"abc": "doerayme", "123": 45.67}`
		expected := jsonblob.Blob{"abc": "doerayme", "123": 45.67}
		output := jsonblob.Blob{}

		// when
		err := output.Scan(original)

		// then
		require.NoError(t, err)
		require.EqualValues(t, expected, output)
	})

	t.Run("leaves Blob unchanged when scanning in nil src", func(t *testing.T) {
		// given
		output := jsonblob.Blob{}

		// when
		err := output.Scan(nil)

		// then
		require.NoError(t, err)
		require.Empty(t, output)
	})

	t.Run("returns error when scanning invalid string into Blob", func(t *testing.T) {
		// given
		original := `"abc": "doerayme", "123": 45.67`
		output := jsonblob.Blob{}

		// when
		err := output.Scan(original)

		// then
		require.Error(t, err)
	})

	t.Run("returns error when scanning invalid type into Blob", func(t *testing.T) {
		// given
		original := 55
		output := jsonblob.Blob{}

		// when
		err := output.Scan(original)

		// then
		require.Error(t, err)
	})

	t.Run("returns error without panicking when scanning into nil Blob", func(t *testing.T) {
		// given
		original := jsonblob.Blob{"abc": "doerayme", "123": 45.67}
		var output *jsonblob.Blob

		// when
		err := output.Scan(original)

		// then
		require.Error(t, err)
	})
}

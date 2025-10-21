package xapm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/michaeldelorenzo/x/pkg/utils/maputils"
)

// xapmtxContextKey is an internal key type for storing transactions in context.Context
type xapmtxContextKey string

// TxKey avoids context name collisions due to its protected nature.
const TxKey = xapmtxContextKey("xapm_current_tx")

var (
	// Apm is the global APM provider instance
	Apm Provider
)

// SetProvider sets the global APM provider
func SetProvider(provider Provider) {
	Apm = provider
}

// GetProvider returns the global APM provider
func GetProvider() Provider {
	return Apm
}

// SendCustomEvent initializes and transmits a CustomEvent to the APM provider.
func SendCustomEvent(reporter EventReporter) {
	if Apm == nil {
		log.Println("WARN: APM provider not initialized, skipping custom event")
		return
	}

	params := eventStructToFlatMap(reporter)
	Apm.RecordCustomEvent(reporter.GetEventType(), params)
}

// StartTransaction is a convenience method for starting a transaction
func StartTransaction(name string) Transaction {
	if Apm == nil {
		log.Println("WARN: APM provider not initialized, returning no-op transaction")
		return &noopTransaction{}
	}
	return Apm.StartTransaction(name)
}

// CtxFromTx embeds the provided transaction in the provided context and returns
// the updated value.
func CtxFromTx(ctx context.Context, tx Transaction) context.Context {
	// Use the transaction's own Context method to embed provider-specific data
	ctx = tx.Context(ctx)
	// Also store our own transaction reference
	return context.WithValue(ctx, TxKey, tx)
}

// TxFromCtx retrieves a Transaction from the provided context.
// Returns nil if a transaction is not found in the context.
func TxFromCtx(ctx context.Context) Transaction {
	t := ctx.Value(TxKey)
	switch tx := t.(type) {
	case Transaction:
		return tx
	default:
		return nil
	}
}

// sanitizeQueryParams converts any arbitrary typed params into strings since DB segments only support
// `booleans`, `numbers`, and `strings`
func sanitizeQueryParams(params map[string]interface{}) map[string]interface{} {
	return maputils.Reduce(
		params,
		make(map[string]interface{}, len(params)),
		func(acc map[string]interface{}, k string, v interface{}) map[string]interface{} {
			// attempt to JSON marshal the input value for best readability,
			// and if that fails, fallback to fmt.Sprintf("%v")
			outBytes, err := json.Marshal(v)
			out := string(outBytes)
			if err != nil {
				out = fmt.Sprintf("%v", v)
			}

			acc[k] = out
			return acc
		},
	)
}

// eventStructToFlatMap flattens event structs to comply with APM standards for reporting events.
func eventStructToFlatMap(data EventReporter) map[string]interface{} {
	result := make(map[string]interface{})

	// First we automatically dereference pointers if needed
	val := reflect.Indirect(reflect.ValueOf(data))

	// If the data we are flattening isn't a struct, return nil
	if val.Kind() != reflect.Struct {
		return nil
	}

	var processStruct func(val reflect.Value, prefix string)
	processStruct = func(val reflect.Value, prefix string) {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			typeField := val.Type().Field(i)

			// Skip unexported fields by checking if the field is accessible (exported).
			if !field.CanInterface() {
				continue
			}

			key := typeField.Name
			if prefix != "" {
				key = prefix + "." + key // Prefixing for nested struct fields.
			}

			// Check if the field is a pointer and dereference it if it points to a supported type.
			if field.Kind() == reflect.Ptr {
				// Check if the pointer is nil; if so, skip or handle as desired.
				if field.IsNil() {
					continue
				}

				// Dereference the pointer for further type checks.
				field = field.Elem()
			}

			// Now, we handle the dereferenced field based on its (possibly dereferenced) type.
			switch field.Kind() {
			case reflect.Struct:
				// Recursively handle nested structs with the current field's key as the new prefix.
				processStruct(field, key)
			case reflect.String, reflect.Bool:
				fallthrough
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				// Directly include numeric and boolean types.
				// Do not include empty values
				if !field.IsZero() {
					result[key] = field.Interface()
				}
			default:
				log.Printf("custom event skipping unsupported field: field %s of type %s", key, field.Type())
			}
		}
	}

	// Kickstart the recursive processing with the initial struct value and an empty prefix.
	processStruct(val, "")
	return result
}

// noopTransaction is a no-op implementation used when APM is not initialized
type noopTransaction struct{}

func (n *noopTransaction) End()                                               {}
func (n *noopTransaction) StartSegment(name string) Segment                   { return &noopSegment{} }
func (n *noopTransaction) StartExternalSegment(url string) Segment            { return &noopSegment{} }
func (n *noopTransaction) StartDBSegment(params *DBSegParams) Segment         { return &noopSegment{} }
func (n *noopTransaction) StartPSQLSegment(params *PSQLDBSegParams) Segment   { return &noopSegment{} }
func (n *noopTransaction) StartKafkaSegment(topicName string) Segment         { return &noopSegment{} }
func (n *noopTransaction) SetWebRequestHTTP(req *http.Request)                {}
func (n *noopTransaction) SetWebResponse(w http.ResponseWriter) http.ResponseWriter { return w }
func (n *noopTransaction) NoticeError(err error)                              {}
func (n *noopTransaction) AddAttribute(key string, value interface{})        {}
func (n *noopTransaction) Context(ctx context.Context) context.Context       { return ctx }

// noopSegment is a no-op implementation
type noopSegment struct{}

func (n *noopSegment) End()                                        {}
func (n *noopSegment) AddAttribute(key string, value interface{}) {}

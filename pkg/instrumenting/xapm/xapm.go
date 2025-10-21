package xapm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/michaeldelorenzo/x/pkg/utils/maputils"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/newrelic/sqlparse"
)

// xapmtxContextKey is an internal key type
// for storing transactions in context.Context
type xapmtxContextKey string

// TxKey avoids context name collisions due to its protected nature.
const TxKey = xapmtxContextKey("xapm_current_tx")

// DBSegParams takes common parameters needed to start a newrelic.DatastoreSegment
type DBSegParams struct {
	Host         string
	DatabaseName string
	Collection   string
	Operation    DBOperation
	DatabaseType DBType
	QueryString  string
	QueryParams  map[string]interface{}
}

// PSQLDBSegParams takes common parameters needed to start a postgres newrelic.DatastoreSegment
type PSQLDBSegParams struct {
	Host         string
	DatabaseName string
	QueryString  string
	QueryParams  map[string]interface{}
}

// EventReporter is an interface that allows the receiver model to be processed as a custom event.
// The exported scalar properties of the implementing struct will be flattened and processed as
// attributes to the New Relic custom event.
type EventReporter interface {
	GetEventType() string
}

type customEvent struct {
	eventType string
	params    map[string]interface{}
}

// SendCustomEvent initializes and transmits a CustomEvent to New Relic.
func SendCustomEvent(reporter EventReporter) {
	e := &customEvent{
		eventType: reporter.GetEventType(),
		params:    eventStructToFlatMap(reporter),
	}

	Apm.RecordCustomEvent(e.eventType, e.params)
}

// StartTransaction is a convenience method for starting a newrelic.Transaction
func StartTransaction(name string) *newrelic.Transaction {
	return Apm.StartTransaction(name)
}

// StartExternalSegment is a convenience method for starting a newrelic.ExternalSegment
func StartExternalSegment(tx *newrelic.Transaction, url string) *newrelic.ExternalSegment {
	return &newrelic.ExternalSegment{
		StartTime: tx.StartSegmentNow(),
		URL:       url,
	}
}

// StartSegment is a convenience method for starting a newrelic.Segment
func StartSegment(tx *newrelic.Transaction, name string) *newrelic.Segment {
	return &newrelic.Segment{
		Name:      name,
		StartTime: tx.StartSegmentNow(),
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

// StartDBSegment is a convenience method for starting a newrelic.DatastoreSegment
func StartDBSegment(tx *newrelic.Transaction, dbSegmentInfo *DBSegParams) *newrelic.DatastoreSegment {
	return &newrelic.DatastoreSegment{
		StartTime:          tx.StartSegmentNow(),
		Host:               dbSegmentInfo.Host,
		DatabaseName:       dbSegmentInfo.DatabaseName,
		Collection:         dbSegmentInfo.Collection,
		Operation:          string(dbSegmentInfo.Operation),
		Product:            dbSegmentInfo.DatabaseType,
		ParameterizedQuery: dbSegmentInfo.QueryString,
		QueryParameters:    sanitizeQueryParams(dbSegmentInfo.QueryParams),
	}
}

func StartPSQLSegment(tx *newrelic.Transaction, dbSegmentInfo *PSQLDBSegParams) *newrelic.DatastoreSegment {
	seg := &newrelic.DatastoreSegment{
		Host:               dbSegmentInfo.Host,
		DatabaseName:       dbSegmentInfo.DatabaseName,
		Product:            DBPostgres,
		ParameterizedQuery: dbSegmentInfo.QueryString,
		QueryParameters:    sanitizeQueryParams(dbSegmentInfo.QueryParams),
	}

	// mutates the segment and sets the `operation` and `collection` by analyzing the raw query
	sqlparse.ParseQuery(seg, dbSegmentInfo.QueryString)

	seg.StartTime = tx.StartSegmentNow()

	return seg
}

// StartKafkaSegment is a convenience method for starting a newrelic.MessageProducerSegment
func StartKafkaSegment(tx *newrelic.Transaction, topicName string) *newrelic.MessageProducerSegment {
	return &newrelic.MessageProducerSegment{
		StartTime:       tx.StartSegmentNow(),
		Library:         "kafka",
		DestinationType: newrelic.MessageTopic,
		DestinationName: topicName,
	}
}

// CtxFromTx embeds the provided transaction in the provided context and returns
// the updated value.
func CtxFromTx(ctx context.Context, tx *newrelic.Transaction) context.Context {
	return context.WithValue(newrelic.NewContext(ctx, tx), TxKey, tx)
}

// TxFromCtx retrieves a *newrelic.Transaction from the provided context.
// Returns nil and outputs a warning if a transaction is not found in the context.
func TxFromCtx(ctx context.Context) *newrelic.Transaction {
	t := ctx.Value(TxKey)
	switch tx := t.(type) {
	case *newrelic.Transaction:
		return tx
	default:
		return nil
	}
}

// EventStructToFlatMap flattens event structs to comply with New Relic standards for reporting events.
func eventStructToFlatMap(data EventReporter) map[string]interface{} {
	result := make(map[string]interface{})

	//: First we automatically dereferences pointers if needed
	val := reflect.Indirect(reflect.ValueOf(data))

	//: If the data we are flattening isn't a struct, return nil
	if val.Kind() != reflect.Struct {
		return nil
	}

	var processStruct func(val reflect.Value, prefix string)
	processStruct = func(val reflect.Value, prefix string) {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			typeField := val.Type().Field(i)

			//: Lets skip unexported fields by checking if the field is accessible (exported).
			if !field.CanInterface() {
				continue
			}

			key := typeField.Name
			if prefix != "" {
				key = prefix + "." + key //: Prefixing for nested struct fields.
			}

			//: Next we check if the field is a pointer and dereference it if it points to a supported type.
			if field.Kind() == reflect.Ptr {
				//: Lets check if the pointer is nil; if so, skip or handle as desired.
				if field.IsNil() {
					continue
				}

				//: Dereference the pointer for further type checks.
				field = field.Elem()
			}

			//: Now, we handle the dereferenced field based on its (possibly dereferenced) type.
			switch field.Kind() {
			case reflect.Struct:
				//: Recursively handle nested structs with the current field's key as the new prefix.
				processStruct(field, key)
			case reflect.String, reflect.Bool:
				fallthrough
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				//: Directly include numeric and boolean types.
				//: Do not include empty values
				if !field.IsZero() {
					result[key] = field.Interface()
				}
			default:
				log.Printf("custom event skipping unsupported field: field %s of type %s", key, field.Type())
			}
		}
	}

	//: Kickstart the recursive processing with the initial struct value and an empty prefix.
	processStruct(val, "")
	return result
}

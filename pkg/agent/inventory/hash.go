package inventory

import (
	"reflect"

	"github.com/cespare/xxhash"
	"google.golang.org/protobuf/proto"
)

func hashMessage(msg proto.Message) uint64 {
	cloned := proto.Clone(msg)
	resetUpdatedAtField(cloned)
	b, _ := proto.Marshal(cloned)
	return xxhash.Sum64(b)
}

type HasResetUpdatedAt interface {
	resetUpdatedAtField()
}

// ResetUpdatedAtField vide le champ "UpdatedAt" s'il existe (et est une string)
func resetUpdatedAtField(msg proto.Message) {
	v := reflect.ValueOf(msg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}

	// On travaille sur la valeur pointée
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}

	field := v.FieldByName("UpdatedAt")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString("")
	}
}

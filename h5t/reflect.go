package h5t

import (
	"fmt"
	"reflect"
)
import (
	"github.com/valoox/h5go/core"
)

// The type of an atomic type
func atomic(k reflect.Kind) (Datatype, error) {
	switch k {
	case reflect.Int8:
		return Int8()
	case reflect.Int16:
		return Int16()
	case reflect.Int32:
		return Int32()
	case reflect.Int64:
		return Int64()
	case reflect.Uint8:
		return Uint8()
	case reflect.Uint16:
		return Uint16()
	case reflect.Uint32:
		return Uint32()
	case reflect.Uint64:
		return Uint64()
	case reflect.Float32:
		return Float32()
	case reflect.Float64:
		return Float64()
	case reflect.Complex64:
		return Complex64()
	case reflect.Complex128:
		return Complex128()
	case reflect.Bool:
		return Bool()
	case reflect.String:
		// Returns a varlength UTF8 string
		return String(-1, true)
	default:
		return -1, fmt.Errorf("Non-atomic kind: %s", k)
	}
}

// A fixed-length array (possibly multi-dimensional)
func array(v reflect.Type, loc core.Location) (Datatype, error) {
	dims := make([]uint, 1, 5)
	dims[0] = uint(v.Len())
	elt := v.Elem()
	for elt.Kind() == reflect.Array {
		dims = append(dims, uint(elt.Len()))
		elt = elt.Elem()
	}
	T, err := parse(elt, loc)
	if err != nil {
		return -1, err
	}
	defer T.Close()
	return NDarray(T, dims...)
}

// A variable-length slice
func slice(v reflect.Type, loc core.Location) (Datatype, error) {
	elt := v.Elem()
	T, err := parse(elt, loc)
	if err != nil {
		return -1, err
	}
	defer T.Close()
	return List(T)
}

// A structure
// This will return the compound type comprising of all the
// fields in this structure.
// Names are the raw structure names, except if the `hdf` tag
// is found, in which case this name is used instead
// Special case is 'ignore' for the hdf name, which means the
// field is IGNORED and NOT SERIALIZED.
// The HDF5 type of the fields is discovered via reflection as
// well, EXCEPT if the `hdftype` tag is present. In that case,
// the name of the type provided is loaded and this is used
// instead. Note that if this lookup raises an error, the entire
// type definition fails.
func structure(v reflect.Type, lc core.Location) (Datatype, error) {
	n := v.NumField()
	fields := make([]Field, 0, n)
	var err error
	for i := 0; i < n; i++ {
		fld := v.Field(i)
		fname := fld.Name
		if tag := fld.Tag.Get("hdf"); tag == "ignore" {
			continue
		} else if tag != "" {
			fname = tag
		}
		var ftype Datatype
		if tag := fld.Tag.Get("hdftype"); tag != "" {
			ftype, err = Open(lc, tag, DefaultAccess)
		} else {
			ftype, err = parse(fld.Type, lc)
		}
		if err != nil {
			return -1, err
		}
		defer ftype.Close()
		fields = append(fields, Field{
			Name:   fname,
			Type:   ftype,
			Offset: int(fld.Offset),
		})
	}
	// All completed succesfully: packing and returning
	return Struct(int(v.Size()), fields...)
}

// Parses the reflected value and returns the correpsonding datatype
func parse(T reflect.Type, ctxt core.Location) (Datatype, error) {
	switch K := T.Kind(); K {
	case reflect.Array:
		// A fixed-length array
		return array(T, ctxt)
	case reflect.Slice:
		return slice(T, ctxt)
	case reflect.Map:
		// A map
		// Might be serializable in the future, but seems
		// like too much work for now
		return -1, fmt.Errorf("Maps are not serializable")
	case reflect.Struct:
		// A structure
		return structure(T, ctxt)
	case reflect.Chan, reflect.Func, reflect.Invalid:
		// -> cannot be serialized
		// Returns error immediately
		return -1, fmt.Errorf("%s not HDF serializable", K)
	case reflect.Uintptr, reflect.UnsafePointer:
		return -1, fmt.Errorf("Unsafe ptr not supported")
	case reflect.Ptr, reflect.Interface:
		// Pointers are dereferenced
		return parse(T.Elem(), ctxt)
	default:
		// An atomic type (int, bool, string...)
		return atomic(K)
	}
}

// Parses the type of the provided structure, and returns the
// corresponding HDF5 datatype
// The `ctxt` location is used to lookup types encoded in the file.
// If none is required (i.e. there is no structure with type
// tags referring to saved types), it is safe to pass 'nil' as
// the context.
// If the object is a pointer/slice (or a structure containing
// pointers or slices) which are not explicitely hdf:"ignore"'d,
// the resulting type will assume a 'flattened' version of the
// type, i.e. one where all references are dereferenced, and where
// all the slices refer to their underlying data.
func Parse(obj interface{}, ctxt core.Location) (Datatype, error) {
	if obj == nil {
		return -1, fmt.Errorf("Nothing provided")
	}
	// Gets the reflected value of the object
	val := reflect.ValueOf(obj)
	return parse(val.Type(), ctxt)
}

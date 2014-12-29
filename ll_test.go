package h5go

import (
	"math/rand"
	"testing"
	"unsafe"
)
import (
	"github.com/valoox/h5go/core"
	"github.com/valoox/h5go/h5d"
	"github.com/valoox/h5go/h5l"
	"github.com/valoox/h5go/h5s"
	"github.com/valoox/h5go/h5t"
)

// A simple array of ints
type Int32 []int32

// The shape of the array
func (i Int32) Shape() (h5s.Dataspace, error) {
	return h5s.CreateSimple([]int{len(i)}, nil)
}

// The ID of the type
func (i Int32) Type() (h5t.Datatype, error) { return h5t.Int32() }

// Reads the data
func (i Int32) ReadPtr() unsafe.Pointer { return unsafe.Pointer(&i[0]) }

// Writes the data
func (i Int32) WritePtr() unsafe.Pointer { return unsafe.Pointer(&i[0]) }

type Struct struct {
	data   [2][2]int32
	cst    float64 `hdf:"constant"`
	length uint32  `hdf:"len"` // hdftype:"int32"`
	foo    complex128
	tag    string   // Some tags
	tags   []string `hdf:"ignore"`
}

// Creates the file
func createfile(path, dataname string, data h5d.Buffer, X *Struct) error {
	f, err := Create(path, true)
	if err != nil {
		return err
	}
	defer f.Close()
	T, err := data.Type()
	if err != nil {
		return err
	}
	if err := T.Commit(f, "int32",
		h5l.DefaultCreate,
		h5t.DefaultCreate,
		h5t.DefaultAccess); err != nil {
		return err
	}
	sh, err := data.Shape()
	if err != nil {
		return err
	}
	d, err := h5d.Create(f, core.Path(dataname), T, sh,
		h5l.DefaultCreate, h5d.DefaultCreate,
		h5d.DefaultAccess)
	if err != nil {
		return err
	}
	if err := d.Write(data, h5s.ALL, h5d.DefaultXfer); err != nil {
		return err
	}
	SType, err := h5t.Parse(X, f)
	if err != nil {
		return err
	}
	if err = SType.Commit(f, "structure", h5l.DefaultCreate,
		h5t.DefaultCreate, h5t.DefaultAccess); err != nil {
		return err
	}
	dat := h5d.Wrap(SType, unsafe.Pointer(X))
	scal, err := h5s.CreateScalar()
	if err != nil {
		return err
	}
	sd, err := h5d.Create(f, "bar", SType, scal,
		h5l.DefaultCreate,
		h5d.DefaultCreate,
		h5d.DefaultAccess)
	if err != nil {
		return err
	}
	if err := sd.Write(dat, h5s.ALL, h5d.DefaultXfer); err != nil {
		return err
	}
	return nil
}

// Reads the data back from the file
func readback(path, dataname string, data h5d.Buffer, X *Struct) error {
	f, err := Open(path, true)
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := h5d.Open(f, core.Path(dataname), h5d.DefaultAccess)
	if err != nil {
		return err
	}
	if err := d.Read(data, h5s.ALL, h5d.DefaultXfer); err != nil {
		return err
	}
	T, err := h5t.Open(f, "structure", h5t.DefaultAccess)
	if err != nil {
		return err
	}
	bfr := h5d.Wrap(T, unsafe.Pointer(X))
	ds, err := h5d.Open(f, "bar", h5d.DefaultAccess)
	if err != nil {
		return err
	}
	if err := ds.Read(bfr, h5s.ALL, h5d.DefaultXfer); err != nil {
		return err
	}

	return nil
}

// Performs a complete test by creating a file and writing some
// data to it, using the low-level interfaces
func TestAll(t *testing.T) {
	N := 100000
	t.Logf("Using %v-length array", N)
	ints := make(Int32, N)
	for i := 0; i < N; i++ {
		ints[i] = rand.Int31()
	}
	X := &Struct{
		data: [2][2]int32{
			[2]int32{1, 2},
			[2]int32{3, 4},
		},
		cst:    12.345,
		length: 15,
		tag:    "little bunny",
		tags:   []string{"foo", "bar", "baz"},
		foo:    complex(12.342, 13.456),
	}
	t.Log("Writing...")
	if err := createfile("./test.h5", "foo", ints, X); err != nil {
		t.Fatal(err)
	}
	t.Log("Writing OK")
	t.Log("Reading...")
	// Reads back the data
	out := make(Int32, len(ints))
	Y := new(Struct)
	if err := readback("./test.h5", "foo", out, Y); err != nil {
		t.Fatal(err)
	}
	t.Log("Reading OK")
	if Y.cst != X.cst || Y.length != X.length || Y.foo != X.foo {
		t.Logf("X=%s, Y=%s", *X, *Y)
		t.Fatalf("Wrong object deserialization !")
	}
	for i, x := range ints {
		if x != out[i] {
			t.Fatalf("Expected %v, got %v", x, out[i])
		}
	}
	t.Log("Match !")
}

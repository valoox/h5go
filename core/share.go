package core

import "fmt"

// Represents the ID of an HDF5 object
type Id int

// All HDF5 object have id handles
type Object interface {
	// The handle for this object
	Id() Id
}

// The interface for objects which can be closed
type Closable interface {
	// Closes, reporting any error which arose
	Close() error
}

// Specific Id for a location (group or file)
type Location interface {
	// The Id of the location
	At() Id
}

// Represents a status returned by the library as a Go error.
// If the status is negative, this is an error, and a status object
// (which implements the error interface) is returned. Otherwise,
// this simply returns `nil` as no error has been reported. This
// allows to use it as a wrapper for all C functions like this:
//     if err := Status(int(C.<call>), <context>); err != nil {
//         <handle error err>
//     }
// The arguments will be passed to the context string using
// fmt.Sprintf. This allows more insighful error reporting and
// debugging (e.g. `Status(<code>, "while opening file %s", <name>)`
// will report the file name in the error being returned)
func Status(code int, context string, args ...interface{}) error {
	if code < 0 {
		return &status{
			num:     code,
			context: fmt.Sprintf(context, args...),
		}
	}
	return nil
}

// Represents a status from the C library. It is usually returned
// as a simple integer, an is considered an error only if this
// status is negative.
type status struct {
	num     int    // The status integer
	context string // The context in which this status occured
}

// status implements the error interface, returning the context
// in which the error was raised, and the HDF5 status returned
func (self *status) Error() string {
	return fmt.Sprintf("Error while %s [status: %v]",
		self.context, self.num)
}

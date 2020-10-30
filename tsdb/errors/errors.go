// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"bytes"
	stderrors "errors"
	"fmt"
	"io"
)

// nilOrMultiError type allows combining multiple errors into one.
type nilOrMultiError []error

// NewMulti returns nilOrMultiError with provided errors added if not nil.
func NewMulti(errs ...error) nilOrMultiError { // nolint:golint
	m := nilOrMultiError{}
	m.Add(errs...)
	return m
}

// Add adds single or many errors to the error list. Each error is added only if not nil.
// If the error is a multiError type, the errors inside multiError are added to the main nilOrMultiError.
func (es *nilOrMultiError) Add(errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		if merr, ok := err.(multiError); ok {
			*es = append(*es, merr.errs...)
			continue
		}
		*es = append(*es, err)
	}
}

// Err returns the error list as an error or nil if it is empty.
func (es nilOrMultiError) Err() error {
	if len(es) == 0 {
		return nil
	}
	return multiError{errs: es}
}

// multiError implements the error interface, and it represents nilOrMultiError with at least one error inside it.
// NOTE: This type is useful to make sure that nil is returned when no error is combined in nilOrMultiError for err != nil
// check to work.
type multiError struct {
	errs nilOrMultiError
}

// Error returns a concatenated string of the contained errors.
func (es multiError) Error() string {
	var buf bytes.Buffer

	if len(es.errs) > 1 {
		fmt.Fprintf(&buf, "%d errors: ", len(es.errs))
	}

	for i, err := range es.errs {
		if i != 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}

	return buf.String()
}

// As finds the first error in multiError slice of error chains matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
func (es multiError) As(target interface{}) bool {
	for _, err := range es.errs {
		if stderrors.As(err, target) {
			return true
		}
	}
	return false
}

// Is returns true if any error in multiError's slice of error chains matches the given target or
// if the target is of multiError type.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func (es multiError) Is(target error) bool {
	if _, ok := target.(multiError); ok {
		return true
	}
	for _, err := range es.errs {
		if stderrors.Is(err, target) {
			return true
		}
	}
	return false
}

// CloseAll closes all given closers while recording error in multiError.
func CloseAll(cs []io.Closer) error {
	errs := NewMulti()
	for _, c := range cs {
		errs.Add(c.Close())
	}
	return errs.Err()
}

// Copyright 2020 The Prometheus Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	stderrors "errors"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNilMultiError(t *testing.T) {
	require.NoError(t, NewMulti().Err())
	require.NoError(t, NewMulti(nil, nil, nil).Err())

	e := NewMulti()
	e.Add()
	require.NoError(t, e.Err())

	e = NewMulti(nil, nil, nil)
	e.Add()
	require.NoError(t, e.Err())

	e = NewMulti()
	e.Add(nil, nil, nil)
	require.NoError(t, e.Err())

	e = NewMulti(nil, nil, nil)
	e.Add(nil, nil, nil)
	require.NoError(t, e.Err())
}

func TestNotNilMultiError(t *testing.T) {
	err := stderrors.New("test1")
	require.Error(t, NewMulti(err).Err())
	require.Error(t, NewMulti(nil, err, nil).Err())

	e := NewMulti(err)
	e.Add()
	require.Error(t, e.Err())

	e = NewMulti(nil, nil, nil)
	e.Add(err)
	require.Error(t, e.Err())

	e = NewMulti(err)
	e.Add(nil, nil, nil)
	require.Error(t, e.Err())

	e = NewMulti(nil, nil, nil)
	e.Add(nil, err, nil)
	require.Error(t, e.Err())
}

func TestMultiError_Error(t *testing.T) {
	err := stderrors.New("test1")

	require.Equal(t, "test1", NewMulti(err).Err().Error())
	require.Equal(t, "test1", NewMulti(err, nil).Err().Error())
	require.Equal(t, "4 errors: test1; test1; test2; test3", NewMulti(err, err, stderrors.New("test2"), nil, stderrors.New("test3")).Err().Error())
}

type customErr struct{ error }

type customErr2 struct{ error }

type customErr3 struct{ error }

func TestMultiError_As(t *testing.T) {
	err := customErr{}

	require.True(t, stderrors.As(err, &err))
	require.True(t, stderrors.As(err, &customErr{}))
	require.False(t, stderrors.As(err, &customErr2{}))
	require.False(t, stderrors.As(err, &customErr3{}))

	// This is just to show limitation of std As.
	require.False(t, stderrors.As(&err, &err))
	require.False(t, stderrors.As(&err, &customErr{}))
	require.False(t, stderrors.As(&err, &customErr2{}))
	require.False(t, stderrors.As(&err, &customErr3{}))

	e := NewMulti(err).Err()
	require.True(t, stderrors.As(e, &customErr{}))
	require.True(t, stderrors.As(e, &multiError{}))
	require.False(t, stderrors.As(e, &customErr2{}))
	require.False(t, stderrors.As(e, &customErr3{}))

	e2 := NewMulti(err, customErr3{}).Err()
	require.True(t, stderrors.As(e2, &customErr{}))
	require.True(t, stderrors.As(e2, &multiError{}))
	require.False(t, stderrors.As(e2, &customErr2{}))
	require.True(t, stderrors.As(e2, &customErr3{}))

	// Wrapped.
	e3 := pkgerrors.Wrap(NewMulti(err, customErr3{}).Err(), "wrap")
	require.True(t, stderrors.As(e3, &customErr{}))
	require.True(t, stderrors.As(e3, &multiError{}))
	require.False(t, stderrors.As(e3, &customErr2{}))
	require.True(t, stderrors.As(e3, &customErr3{}))

	// This is just to show limitation of std As.
	e4 := NewMulti(err, &customErr3{}).Err()
	require.False(t, stderrors.As(e4, &customErr2{}))
	require.False(t, stderrors.As(e4, &customErr3{}))
}

func TestMultiError_Is(t *testing.T) {
	err := customErr{}

	require.True(t, stderrors.Is(err, err))
	require.True(t, stderrors.Is(err, customErr{}))
	require.False(t, stderrors.Is(err, customErr2{}))
	require.False(t, stderrors.Is(err, customErr3{}))

	require.True(t, stderrors.Is(&err, &err))
	require.False(t, stderrors.Is(&err, &customErr{}))
	require.False(t, stderrors.Is(&err, &customErr2{}))
	require.False(t, stderrors.Is(&err, &customErr3{}))

	e := NewMulti(err).Err()
	require.True(t, stderrors.Is(e, customErr{}))
	require.True(t, stderrors.Is(e, multiError{}))
	require.False(t, stderrors.Is(e, &multiError{}))
	require.False(t, stderrors.Is(e, customErr2{}))
	require.False(t, stderrors.Is(e, customErr3{}))

	e2 := NewMulti(err, customErr3{}).Err()
	require.True(t, stderrors.Is(e2, customErr{}))
	require.True(t, stderrors.Is(e2, multiError{}))
	require.False(t, stderrors.Is(e2, customErr2{}))
	require.True(t, stderrors.Is(e2, customErr3{}))

	// Wrapped.
	e3 := pkgerrors.Wrap(NewMulti(err, customErr3{}).Err(), "wrap")
	require.True(t, stderrors.Is(e3, customErr{}))
	require.True(t, stderrors.Is(e3, multiError{}))
	require.False(t, stderrors.Is(e3, customErr2{}))
	require.True(t, stderrors.Is(e3, customErr3{}))

	exact := &customErr3{}
	e4 := NewMulti(err, exact).Err()
	require.True(t, stderrors.Is(e4, customErr{}))
	require.True(t, stderrors.Is(e4, multiError{}))
	require.False(t, stderrors.Is(e4, customErr2{}))
	require.False(t, stderrors.Is(e4, &customErr3{}))
	require.True(t, stderrors.Is(e4, exact))
}

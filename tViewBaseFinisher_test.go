// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

import (
	"reflect"
	"testing"

	testify "github.com/stretchr/testify/assert"
)

// todo: Add tests for tViewBaseFinisher.trFinish

// test for tViewBaseFinisher.protectFromPanic and tViewBaseFinisher.invoke
func TestViewBaseFinisherProtectFromPanic(t *testing.T) {

	test := func(t *testing.T, finisher *tViewBaseFinisher, f func(*tViewBaseFinisher), tests func(*testing.T, *tViewBaseFinisher), isRecovered bool) {

		defer func(t *testing.T, isRecovered bool) {
			test := testify.Nil
			if isRecovered {
				test = testify.NotNil
			}
			test(t, recover())
		}(t, isRecovered)

		f(finisher)
		tests(t, finisher)
	}

	flushRecoveredPanics := func(finisher *tViewBaseFinisher) {
		finisher.recoveredPanics = nil
	}

	panicer := func(finisher *tViewBaseFinisher) {
		defer finisher.protectFromPanic()
		flushRecoveredPanics(finisher)
		panic(nil)
	}

	notPanicer := func(finisher *tViewBaseFinisher) {
		defer finisher.protectFromPanic()
		flushRecoveredPanics(finisher)
		_ = nil
	}

	invoked := func(what string, to *string) {
		*to = what
	}

	var s string

	invokerWithPanic := func(finisher *tViewBaseFinisher) {
		f := reflect.ValueOf(invoked)
		args := []reflect.Value{
			reflect.ValueOf("with panic"),
			reflect.ValueOf(nil),
		}
		finisher.invoke(f, args)
	}

	invokerWithoutPanic := func(finisher *tViewBaseFinisher) {
		f := reflect.ValueOf(invoked)
		args := []reflect.Value{
			reflect.ValueOf("without panic"),
			reflect.ValueOf(&s),
		}
		finisher.invoke(f, args)
	}

	testsPanicWas := func(t *testing.T, finisher *tViewBaseFinisher) {
		testify.Len(t, finisher.recoveredPanics, 1)
		testify.NotNil(t, finisher.recoveredPanics[0])
	}

	testsPanicWasnt := func(t *testing.T, finisher *tViewBaseFinisher) {
		testify.Empty(t, finisher.recoveredPanics)
	}

	testsStringEmpty := func(t *testing.T, finisher *tViewBaseFinisher) {
		testsPanicWas(t, finisher)
		testify.Empty(t, s)
	}

	testsStringNotEmpty := func(t *testing.T, finisher *tViewBaseFinisher) {
		testsPanicWasnt(t, finisher)
		testify.NotEmpty(t, s)
		s = ""
	}

	finisherWithPanicGuard := &tViewBaseFinisher{}
	finisherWithoutPanicGuard := &tViewBaseFinisher{}

	finisherWithPanicGuard.flags = cSendCallbackEnablePanicGuard

	test(t, finisherWithPanicGuard, panicer, testsPanicWas, false)
	test(t, finisherWithPanicGuard, notPanicer, testsPanicWasnt, false)
	test(t, finisherWithoutPanicGuard, panicer, testsPanicWasnt, true)
	test(t, finisherWithoutPanicGuard, notPanicer, testsPanicWasnt, false)

	test(t, finisherWithPanicGuard, invokerWithPanic, testsStringEmpty, false)
	test(t, finisherWithPanicGuard, invokerWithoutPanic, testsStringNotEmpty, false)
	test(t, finisherWithoutPanicGuard, invokerWithPanic, testsStringEmpty, true)
	test(t, finisherWithoutPanicGuard, invokerWithoutPanic, testsStringNotEmpty, false)
}

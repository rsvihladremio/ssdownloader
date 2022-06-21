/*
   Copyright 2022 Ryan SVIHLA

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package sendsafely

import (
	"testing"
)

// This is the default happy path test, no errors
func TestCalculateExecutionCalls(t *testing.T) {
	req := calculateExecutionCalls(1)
	if len(req) != 1 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)
	}
}

func TestCalculateExecutionCallsBeforeBoundary(t *testing.T) {
	req := calculateExecutionCalls(24)
	if len(req) != 1 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 24 {
		t.Errorf("expected end segment to be 24 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)
	}
}

func TestCalculateExecutionCallsAtBoundary(t *testing.T) {
	req := calculateExecutionCalls(25)
	if len(req) != 1 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 25 {
		t.Errorf("expected end segment to be 25 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)
	}
}
func TestCalculateExecutionCallsOverTheBoundary(t *testing.T) {
	req := calculateExecutionCalls(26)
	if len(req) != 2 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
	part := req[0]
	if part.EndSegment != 25 {
		t.Errorf("expected end segment to be 25 but was '%v'", part.EndSegment)
	}
	if part.StartSegment != 1 {
		t.Errorf("expected end segment to be 1 but was '%v'", part.StartSegment)

	}

	part2 := req[1]
	if part2.EndSegment != 26 {
		t.Errorf("expected end segment to be 26 but was '%v'", part2.EndSegment)
	}
	if part2.StartSegment != 26 {
		t.Errorf("expected end segment to be 26 but was '%v'", part2.StartSegment)
	}
}

func TestNoParts(t *testing.T) {
	req := calculateExecutionCalls(0)
	if len(req) != 0 {
		t.Fatalf("expected '%v' request parts but got '%v'", 1, len(req))
	}
}

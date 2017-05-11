package macomp

import (
	"reflect"
	"testing"
)

func TestGetConstraints(t *testing.T) {
	tests := []struct {
		query        string
		goldSureSeps IntSet
		goldAmbSeps  IntSet
		goldSeps     IntSet
		start        int
		end          int
	}{
		{
			query:        "|ジャパリ?まん|は美味しい。",
			goldSureSeps: IntSet{0: struct{}{}, 6: struct{}{}},
			goldAmbSeps:  IntSet{4: struct{}{}},
			goldSeps:     IntSet{0: struct{}{}, 4: struct{}{}, 6: struct{}{}},
			start:        0,
			end:          6,
		},
		{
			query:        "ジャパリまんは美味しい。",
			goldSureSeps: IntSet{},
			goldAmbSeps:  IntSet{},
			goldSeps:     IntSet{},
			start:        -1,
			end:          -1,
		},
	}
	for _, test := range tests {
		sysSureSeps, sysAmbSeps, sysSeps, sysStart, sysEnd := GetConstraints(test.query)
		if !reflect.DeepEqual(sysSureSeps, test.goldSureSeps) {
			t.Errorf("Expected %v but got %v", test.goldSureSeps, sysSureSeps)
		}
		if !reflect.DeepEqual(sysAmbSeps, test.goldAmbSeps) {
			t.Errorf("Expected %v but got %v", test.goldAmbSeps, sysAmbSeps)
		}
		if !reflect.DeepEqual(sysSeps, test.goldSeps) {
			t.Errorf("Expected %v but got %v", test.goldSeps, sysSeps)
		}
		if sysStart != test.start {
			t.Errorf("Expected %v but got %v", test.start, sysStart)
		}
		if sysEnd != test.end {
			t.Errorf("Expected %v but got %v", test.end, sysEnd)
		}
	}
}

func TestGetSegments(t *testing.T) {
	tests := []struct {
		surfs    []string
		goldSeps IntSet
	}{
		{
			surfs:    []string{"ジャパリ", "まん", "は", "美味しい"},
			goldSeps: IntSet{0: struct{}{}, 4: struct{}{}, 6: struct{}{}, 7: struct{}{}, 11: struct{}{}},
		},
	}
	for _, test := range tests {
		sysSeps := GetSegments(test.surfs)
		if !reflect.DeepEqual(sysSeps, test.goldSeps) {
			t.Errorf("Expected %v but got %v", test.goldSeps, sysSeps)
		}
	}
}

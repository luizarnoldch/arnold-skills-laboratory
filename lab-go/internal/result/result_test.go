package result_test

import (
	"testing"

	"skills-laboratory/lab-go/internal/result"
)

func TestRoundRateAndCorrect(t *testing.T) {
	if got := result.RoundRate(3, 3); got != 1 {
		t.Fatalf("rate=%v", got)
	}
	if got := result.RoundRate(1, 3); got != 0.33 {
		t.Fatalf("rate=%v", got)
	}
	if !result.IsCorrect(1.0, true, 0.5) {
		t.Fatal("expected correct true positive")
	}
	if !result.IsCorrect(0.33, false, 0.5) {
		t.Fatal("expected correct true negative (rate below majority)")
	}
	if !result.IsCorrect(0.0, false, 0.5) {
		t.Fatal("true negative failed")
	}
	if result.IsCorrect(0.0, true, 0.5) {
		t.Fatal("false negative should be incorrect")
	}
	if result.IsCorrect(1.0, false, 0.5) {
		t.Fatal("false positive should be incorrect")
	}
}

package main

import (
	"fmt"
	"testing"
)

func TestApp(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(fmt.Sprintf("RECOVERED FROM PANIC. ERROR: %s", r))
		}
	}()
	example_app := DefaultApp()
	example_app.ListenAndServe()
}

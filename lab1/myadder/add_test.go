package myadder

import "testing"

func TestAdd(t *testing.T) {
    want := "hello world"
    got := Add("hello", "world")
    if want != got {
        t.Errorf("Error in myadder.Add; Want hello world, Got %s", got)
    }   
}


package bloomfilter

import "testing"

func TestBloomFilter(t *testing.T) {
	f := NewEstimated(10, 0.01)
	f.Add([]byte("wow"))
	contains := f.Test([]byte("wow"))
	if contains != true {
		t.Errorf("f.Test([]byte(%q)) = %t, want %t", "wow", contains, true)
	}

	contains = f.Test([]byte("amazing"))
	if contains != false {
		t.Errorf("f.Test([]byte(%q)) = %t, want %t", "amazing", contains, false)
	}

	f.Add([]byte("amazing"))
	contains = f.Test([]byte("amazing"))
	if contains != true {
		t.Errorf("f.Test([]byte(%q)) = %t, want %t", "amazing", contains, true)
	}
}

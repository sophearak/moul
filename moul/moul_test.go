package moul

import "testing"

func TestGetImageDimension(t *testing.T) {
	w, h := GetImageDimension("../assets/preloader.f75eb900.gif")

	if w != 20 && h != 20 {
		t.Errorf("GetImageDimension isn't working right, got w and h: %d %d, want: %d %d.", w, h, 20, 20)
	}
}

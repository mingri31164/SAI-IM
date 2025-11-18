package bitmap

import "testing"

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(10)

	b.Set("pppp")
	b.Set("222")
	b.Set("pppp")
	b.Set("ccc")
	b.Set("ddd")
	b.Set("eee")

	for _, bit := range b.bits {
		t.Logf("%b,%v\n", bit, bit)
	}
}

package domain2

import "testing"

func TestItem_ID(t *testing.T) {
	for i := 0; i < 100; i++ {
		i1 := Item{"a": "1", "b": "2"}
		i2 := Item{"a": "1", "b": "2"}
		if i1.ID() != i2.ID() {
			t.Error("same content should have same id by default")
		}
	}

	i := Item{"a": "1"}
	i[ItemKeyID] = "id1"
	if i.ID() != "id1" {
		t.Error("id should be set")
	}
}

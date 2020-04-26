package domain2

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type (
	diffItem struct {
		diffType diffmatchpatch.Operation
		text     string
	}
)

// CompareItemList computes the differences between ItemLists.
func CompareItemList(list1, list2 ItemList) *Updates {
	// Get ids and compare only id
	idToItem1 := list1.idToItem()
	idToItem2 := list2.idToItem()
	ids1 := make([]string, 0)
	for i := range idToItem1 {
		ids1 = append(ids1, i)
	}
	ids2 := make([]string, 0)
	for i := range idToItem2 {
		ids2 = append(ids2, i)
	}
	// Currently I ignore item order in compare phase
	sort.Strings(ids1)
	sort.Strings(ids2)

	// Compare IDs
	ids1s := strings.Join(ids1, "\n")
	ids2s := strings.Join(ids2, "\n")
	diffs := diff(ids1s, ids2s)
	addedItems := make([]Item, 0)
	removedItems := make([]Item, 0)
	changedItems := make([]ItemChange, 0)
	for _, d := range diffs {
		switch d.diffType {
		case diffmatchpatch.DiffInsert:
			addedItems = append(addedItems, idToItem2[d.text])
		case diffmatchpatch.DiffDelete:
			removedItems = append(removedItems, idToItem1[d.text])
		case diffmatchpatch.DiffEqual:
			item1 := idToItem1[d.text]
			item2 := idToItem2[d.text]
			changed := compareItem(item1, item2)
			if changed != nil {
				changedItems = append(changedItems, *changed)
			}
		default:
			continue
		}
	}

	return &Updates{addedItems, removedItems, changedItems}
}

// diffItem compare 2 items.
// return null if given items are exactly same.
func compareItem(item1, item2 Item) *ItemChange {
	keys1 := make([]string, len(item1))
	keys2 := make([]string, len(item2))

	// Extract keys
	index := 0
	for k := range item1 {
		keys1[index] = k
		index++
	}
	index = 0
	for k := range item2 {
		keys2[index] = k
		index++
	}
	// Sort keys
	sort.Strings(keys1)
	sort.Strings(keys2)

	// Compare them
	keys1s := strings.Join(keys1, "\n")
	keys2s := strings.Join(keys2, "\n")
	diffs := diff(keys1s, keys2s)
	addedKeys := make(map[string]string, 0)
	removedKeys := make(map[string]string, 0)
	changedKeys := make(map[string]ItemValueChange, 0)
	for _, d := range diffs {
		switch d.diffType {
		case diffmatchpatch.DiffInsert:
			addedKeys[d.text] = item2[d.text]
		case diffmatchpatch.DiffDelete:
			removedKeys[d.text] = item1[d.text]
		case diffmatchpatch.DiffEqual:
			v1 := item1[d.text]
			v2 := item2[d.text]
			if v1 != v2 {
				changedKeys[d.text] = ItemValueChange{
					Old: v1,
					New: v2,
				}
			}
		default:
			continue
		}
	}
	if len(addedKeys) == 0 && len(removedKeys) == 0 && len(changedKeys) == 0 {
		return nil
	}
	return &ItemChange{addedKeys, removedKeys, changedKeys}
}

func diff(v1, v2 string) []diffItem {
	v1 = strings.Trim(v1, " \t\n")
	if len(v1) > 0 {
		v1 = v1 + "\n"
	}
	v2 = strings.Trim(v2, " \t\n")
	if len(v2) > 0 {
		v2 = v2 + "\n"
	}

	d := diffmatchpatch.New()
	a, b, c := d.DiffLinesToChars(v1, v2)
	diffs := d.DiffMain(a, b, false)
	res := make([]diffItem, 0)
	for _, diff := range d.DiffCharsToLines(diffs, c) {
		text := strings.Trim(diff.Text, " \t\n")
		res = append(res, diffItem{diff.Type, text})
	}
	return res
}

func (u *ItemChange) String() string {
	var res = make(map[string]map[string]string)
	res["added"] = u.AddedKeys
	res["removed"] = u.RemovedKeys

	changed := make(map[string]string)
	for k, i := range u.ChangedKeys {
		changed[k] = fmt.Sprintf("%s -> %s", i.Old, i.New)
	}
	res["changed"] = changed

	b, _ := json.Marshal(res)
	return string(b)
}

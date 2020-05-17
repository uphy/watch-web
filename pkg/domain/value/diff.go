package value

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
func CompareItemList(list1, list2 ItemList) Updates {
	list1 = list1.forCompare()
	list2 = list2.forCompare()

	// Get ids and compare only id
	idToItem1 := make(map[string]Item)
	ids1 := make([]string, len(list1))
	for i, item := range list1 {
		id := item.ID()
		idToItem1[id] = item
		ids1[i] = id
	}
	idToItem2 := make(map[string]Item)
	ids2 := make([]string, len(list2))
	for i, item := range list2 {
		id := item.ID()
		idToItem2[id] = item
		ids2[i] = id
	}

	// Compare IDs
	ids1s := strings.Join(ids1, "\n")
	ids2s := strings.Join(ids2, "\n")
	diffs := diff(ids1s, ids2s)
	updates := make(Updates, 0)
	for _, d := range diffs {
		switch d.diffType {
		case diffmatchpatch.DiffInsert:
			updates = append(updates, *updateAdd(idToItem2[d.text]))
		case diffmatchpatch.DiffDelete:
			updates = append(updates, *updateRemove(idToItem1[d.text]))
		case diffmatchpatch.DiffEqual:
			item1 := idToItem1[d.text]
			item2 := idToItem2[d.text]
			changed := compareItem(item1, item2)
			if changed != nil {
				updates = append(updates, *updateChange(changed))
			}
		default:
			continue
		}
	}

	return updates
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
	return &ItemChange{item2, addedKeys, removedKeys, changedKeys}
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
		splitted := strings.Split(text, "\n")
		for _, s := range splitted {
			res = append(res, diffItem{diff.Type, s})
		}
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

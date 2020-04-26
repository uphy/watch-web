package domain2

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	// ItemKeyID is the item's special key which is the identifier of Item.
	// By comparing old ID and new ID, we can detect CHANGED in addition to ADDED/REMOVED in diff.
	ItemKeyID = "__id__"
)

type (
	/*
	 * Design Note:
	 *
	 * I use mainly Slack to notify the web update.
	 * Web content often represented as Tree.
	 * But Slack cannot display tree well and
	 * tree structure is not clear way to understand updated content.
	 *
	 * So I represent content as []map[string]string.
	 * Item is the map[string]string and ItemList is the list of items.
	 * Compare two ItemList, old one and new one and report it to Slack.
	 * We can get
	 * - Which item is added, deleted
	 * - Which item is changed (also need to detect same item between old and new one.  maybe need Item ID)
	 */

	// JobContext is the context of job
	JobContext struct {
		Log *logrus.Entry
	}

	// Item is the free format key-value object
	// In case you watch the search result,
	// each search result item should be represented as Item.
	Item map[string]string

	// ItemList is the ordered list of the Item.
	ItemList []Item

	// Source represents the content source.
	// Source can get the list of items.
	// In case source is the plain text, it represented as `{"source":"xxxxx"}`
	Source interface {
		Fetch(ctx *JobContext) ItemList
	}

	// Transformer transforms the source item list.
	Transformer interface {
		Transform(ctx *JobContext, itemList ItemList) (ItemList, error)
	}

	// Notifier notifies the source content update.
	Notifier interface {
		Notify(ctx *JobContext, updates Updates) error
	}

	// Updates is the result of Diff function.
	Updates struct {
		Added   []Item
		Removed []Item
		Changed []ItemChange
	}

	// ItemChange represents changed item.
	ItemChange struct {
		// AddedKeys exist in new item but not exist in old item.
		AddedKeys map[string]string
		// RemovedKeys don't exist in new item but exist in old item.
		RemovedKeys map[string]string
		// ChangedKeys exist in both old and new one but their value was changed.
		ChangedKeys map[string]ItemValueChange
	}
	// ItemValueChange has a change of the Item value.
	ItemValueChange struct {
		Old string
		New string
	}
)

// Clone returns the deep-copied struct.
func (i ItemList) Clone() ItemList {
	clone := make([]Item, len(i))
	for i, item := range i {
		clone[i] = item.Clone()
	}
	return clone
}

func (i ItemList) idToItem() map[string]Item {
	v := make(map[string]Item)
	for _, item := range i {
		v[item.ID()] = item
	}
	return v
}

func (i ItemList) ids() []string {
	ids := make([]string, len(i))
	for j, item := range i {
		id := item.ID()
		ids[j] = id
	}
	return ids
}

// Clone returns the deep-copied struct.
func (i Item) Clone() Item {
	clone := make(map[string]string, len(i))
	for k, v := range i {
		clone[k] = v
	}
	return clone
}

// NewItem create new Item
func NewItem(m map[string]string) Item {
	item := Item(m).Clone()
	return item
}

// NewItemFromJSON create a new Item from JSON object string.
func NewItemFromJSON(jsonString string) (Item, error) {
	var m map[string]string
	if err := json.Unmarshal([]byte(jsonString), &m); err != nil {
		return nil, fmt.Errorf("cannot unmarshal json string: json=%s, err=%w", jsonString, err)
	}
	return NewItem(m), nil
}

// ID returns the item ID.
// If item ID is set, return it.
// Otherwise, Item ID is determined from item content.
func (i Item) ID() string {
	id, exist := i[ItemKeyID]
	if exist {
		return id
	}
	j := i.JSON()
	sum := md5.Sum([]byte(j))
	return hex.EncodeToString(sum[:])
}

// SetID set id.
func (i Item) SetID(id string) {
	i[ItemKeyID] = id
}

// JSON converts item as JSON.
func (i Item) JSON() string {
	b, _ := json.Marshal(i)
	return string(b)
}

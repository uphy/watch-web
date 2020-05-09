package domain

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

const (
	// ItemKeyID is the item's special key which is the identifier of Item.
	// By comparing old ID and new ID, we can detect CHANGED in addition to ADDED/REMOVED in diff.
	ItemKeyID                   = "__id__"
	UpdateTypeAdd    UpdateType = "add"
	UpdateTypeRemove UpdateType = "remove"
	UpdateTypeChange UpdateType = "change"
)

type (
	// Item is the free format key-value object
	// In case you watch the search result,
	// each search result item should be represented as Item.
	Item map[string]string

	// ItemList is the ordered list of the Item.
	ItemList []Item

	// Updates is the result of Diff function.
	Updates []Update
	// UpdateType represents one of the add, remove, change.
	UpdateType string
	// Update is the item update
	Update struct {
		Type UpdateType
		// Add is an event represents Item added
		Add Item
		// Remove is an event represents Item removed
		Remove Item
		Change *ItemChange
	}
	// ItemChange represents changed item.
	ItemChange struct {
		Item Item `json:"item"`
		// AddedKeys exist in new item but not exist in old item.
		AddedKeys map[string]string `json:"add,omitempty"`
		// RemovedKeys don't exist in new item but exist in old item.
		RemovedKeys map[string]string `json:"remove,omitempty"`
		// ChangedKeys exist in both old and new one but their value was changed.
		ChangedKeys map[string]ItemValueChange `json:"change,omitempty"`
	}
	// ItemValueChange has a change of the Item value.
	ItemValueChange struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
)

func updateAdd(item Item) *Update {
	return &Update{Type: UpdateTypeAdd, Add: item}
}

func updateRemove(item Item) *Update {
	return &Update{Type: UpdateTypeRemove, Remove: item}
}

func updateChange(item *ItemChange) *Update {
	return &Update{Type: UpdateTypeChange, Change: item}
}

func (u *Update) UnmarshalJSON(data []byte) error {
	var m map[UpdateType]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to unmarshal:%w", err)
	}
	if len(m) != 1 {
		return errors.New("Either of 'add', 'remove', 'change' is required")
	}
	for k, v := range m {
		u.Type = k
		switch k {
		case UpdateTypeAdd:
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal:%w", err)
			}
			var item Item
			if err := json.Unmarshal(b, &item); err != nil {
				return fmt.Errorf("failed to unmarshal:%w", err)
			}
			u.Add = item
		case UpdateTypeRemove:
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal:%w", err)
			}
			var item Item
			if err := json.Unmarshal(b, &item); err != nil {
				return fmt.Errorf("failed to unmarshal:%w", err)
			}
			u.Remove = item
		case UpdateTypeChange:
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal:%w", err)
			}
			var change ItemChange
			if err := json.Unmarshal(b, &change); err != nil {
				return fmt.Errorf("failed to unmarshal:%w", err)
			}
			u.Change = &change
		}
		break
	}
	return nil
}

func (u *Update) item() interface{} {
	if u.Add != nil {
		return u.Add
	}
	if u.Remove != nil {
		return u.Remove
	}
	if u.Change != nil {
		return u.Change
	}
	return nil
}

func (u *Update) MarshalJSON() ([]byte, error) {
	m := map[UpdateType]interface{}{
		u.Type: u.item(),
	}
	return json.Marshal(m)
}

func NewItemListFromJSON(j string) (i ItemList, err error) {
	j = strings.Trim(j, " \t\n")
	if j == "" {
		return make(ItemList, 0), nil
	}
	err = json.Unmarshal([]byte(j), &i)
	return
}

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

func (i ItemList) Empty() bool {
	return len(i) == 0
}

// JSON converts item list to JSON.
func (i ItemList) JSON() string {
	b, _ := json.Marshal(i)
	return string(b)
}

// YAML converts item list to YAML.
func (i ItemList) YAML() string {
	b, _ := yaml.Marshal(i)
	return string(b)
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

func NewItemForSource(s string) Item {
	return Item{
		"source": s,
	}
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

func (i Item) Empty() bool {
	keys := len(i)
	if _, idExist := i[ItemKeyID]; idExist {
		keys--
	}
	return keys == 0
}

func (u Updates) Changes() bool {
	return len(u) > 0
}

// YAML converts updates to YAML.
func (u Updates) YAML() string {
	b, _ := yaml.Marshal(u)
	return string(b)
}

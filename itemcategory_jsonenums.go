// generated by jsonenums -type=itemCategory; DO NOT EDIT

package goprove

import (
	"encoding/json"
	"fmt"
)

var (
	_itemCategoryNameToValue = map[string]itemCategory{
		"minimumCriteria": minimumCriteria,
		"goodCitizen":     goodCitizen,
		"extraCredit":     extraCredit,
	}

	_itemCategoryValueToName = map[itemCategory]string{
		minimumCriteria: "minimumCriteria",
		goodCitizen:     "goodCitizen",
		extraCredit:     "extraCredit",
	}
)

func init() {
	var v itemCategory
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_itemCategoryNameToValue = map[string]itemCategory{
			interface{}(minimumCriteria).(fmt.Stringer).String(): minimumCriteria,
			interface{}(goodCitizen).(fmt.Stringer).String():     goodCitizen,
			interface{}(extraCredit).(fmt.Stringer).String():     extraCredit,
		}
	}
}

// MarshalJSON is generated so itemCategory satisfies json.Marshaler.
func (r itemCategory) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _itemCategoryValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid itemCategory: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so itemCategory satisfies json.Unmarshaler.
func (r *itemCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("itemCategory should be a string, got %s", data)
	}
	v, ok := _itemCategoryNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid itemCategory %q", s)
	}
	*r = v
	return nil
}

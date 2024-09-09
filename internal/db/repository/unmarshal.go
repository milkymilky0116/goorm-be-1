package repository

import (
	"encoding/json"
	"fmt"
)

func (r *Role) UnmarshalJSON(b []byte) error {
	var roleStr string
	if err := json.Unmarshal(b, &roleStr); err != nil {
		return err
	}
	switch roleStr {
	case string(RoleStudent):
		*r = RoleStudent
	case string(RoleAdmin):
		*r = RoleAdmin
	default:
		return fmt.Errorf("invalid role: %s", roleStr)
	}
	return nil
}
func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

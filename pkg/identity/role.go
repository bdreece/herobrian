//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Role -trimprefix Role
package identity

import (
	"encoding/gob"
	"errors"
	"fmt"
)

type Role int

const (
	RoleUser Role = iota
	RoleModerator
	RoleAdmin
	RoleSuper
)

func (r Role) Authorize(claims *ClaimSet) error {
	if claims.Role < int64(r) {
		return errors.Join(
			fmt.Errorf("failed to authorize user for role: %q", r),
			ErrUnauthorized,
		)
	}

	return nil
}

func init() {
	gob.Register(new(Role))
}

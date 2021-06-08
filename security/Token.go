package security

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Appliance ExternalID Type
const Appliance = "Appliance"

// Session ExternalID Type
const Session = "Session"

// User ExternalID Type
const User = "User"

// JwtToken represents the parsed Token from Authentication Header
type JwtToken struct {
	// UserID is id of user matchimg the token
	UserID         uuid.UUID   `json:"user,omitempty"`
	UserName       string      `json:"name,omitempty"`
	DisplayName    string      `json:"displayName,omitempty"`
	UserGroupIDs   []uuid.UUID `json:"usergroupIds,omitempty"`
	TenantID       uuid.UUID   `json:"tenant,omitempty"`
	TenantName     string      `json:"tenantName,omitempty"`
	ExternalID     string      `json:"externalId,omitempty"`
	ExternalIDType string      `json:"externalIdType,omitempty"`
	Scopes         []string    `json:"scope,omitempty"`
	Admin          bool        `json:"admin,omitempty"`
	Raw            string      `json:"-"`
	jwt.StandardClaims
}

func (token *JwtToken) isValidForScope(allowedScopes []string) bool {
	permissiveTokenScopes := []string{}
	nonPermissiveTokenScopes := []string{}

	for _, tokenScope := range token.Scopes {
		if strings.HasPrefix(tokenScope, "-") {
			nonPermissiveTokenScopes = append(nonPermissiveTokenScopes, tokenScope[1:])
		} else {
			permissiveTokenScopes = append(permissiveTokenScopes, tokenScope)
		}
	}

	fmt.Println("nonPermissiveTokenScopes", nonPermissiveTokenScopes)
	fmt.Println()
	fmt.Println("permissiveTokenScopes", permissiveTokenScopes)
	fmt.Println()

	if len(nonPermissiveTokenScopes) > 0 {
		if isScopePresent(nonPermissiveTokenScopes, allowedScopes, true, false) {
			return false
		}
	}

	return isScopePresent(permissiveTokenScopes, allowedScopes, false, true)
}

func isScopePresent(scopes []string, scopeToCheck []string, isNegativeScopeCheck bool, shouldAllScopeMatch bool) bool {
	if !isNegativeScopeCheck {
		if ok, _ := inArray("*", scopes); ok {
			return true
		}
	}
	allScopesMatched := shouldAllScopeMatch
	for _, allowedScope := range scopeToCheck {
		if ok, _ := inArray(allowedScope, scopes); !ok {
			scopeParts := strings.Split(allowedScope, ":")
			if ok, _ := inArray(scopeParts[0]+":*", scopes); !ok {
				if ok, _ := inArray("*:"+scopeParts[1], scopes); !ok {
					allScopesMatched = !shouldAllScopeMatch
					break
				}
			}
		}
	}
	return allScopesMatched
}

func inArray(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

package domain

import (
	"regexp"
)

var JoinClassRegex = regexp.MustCompile(`^_Join:[A-Za-z0-9_]+:[A-Za-z0-9_]+`)

var ClassAndFieldRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

var RoleRegex = regexp.MustCompile(`^role:.*`)

var PublicRegex = regexp.MustCompile(`^\*$`)

var AuthenticatedRegex = regexp.MustCompile(`^authenticated$`)

var RequiresAuthenticationRegex = regexp.MustCompile(`^requiresAuthentication$`)

var ProtectedFieldsRegexes = []regexp.Regexp{*RoleRegex, *PublicRegex, *AuthenticatedRegex}

var ClassLevelPermissionsRegexes = []regexp.Regexp{*RoleRegex, *PublicRegex, *RequiresAuthenticationRegex}

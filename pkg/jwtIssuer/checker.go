package jwtIssuer

type field string

const (
	fnil field = ""
	fiss field = "iss"
	fsub field = "sub"
	faud field = "aud"
	fjti field = "jti"
)

type checker struct {
	issuers   map[string]bool
	audiences map[string]bool
	subjects  map[string]bool
	jwtids    map[string]bool
}

func (c checker) ParseFieldName(s string) field {
	switch s {
	case "issuer", "iss":
		return fiss
	case "subject", "sub":
		return fsub
	case "audience", "aud":
		return faud
	case "jwt_id", "jti":
		return fjti
	default:
		return fnil
	}
}

func (c *checker) Add(f field, value string) {
	switch f {
	case fiss:
		c.issuers[value] = true
	case faud:
		c.audiences[value] = true
	case fsub:
		c.subjects[value] = true
	case fjti:
		c.subjects[value] = true
	}
}

func (c *checker) has(f field, value string) bool {
	ok := false
	switch f {
	case fiss:
		_, ok = c.issuers[value]
	case faud:
		_, ok = c.audiences[value]
	case fsub:
		_, ok = c.subjects[value]
	case fjti:
		_, ok = c.subjects[value]
	}
	return ok
}

func (c *checker) HasIssue(iss string) bool {
	return c.has(fiss, iss)
}

func (c *checker) HasAudience(aud string) bool {
	return c.has(faud, aud)
}

func (c *checker) HasSubject(sub string) bool {
	return c.has(fsub, sub)
}

func (c *checker) HasJWTID(jti string) bool {
	return c.has(fjti, jti)
}

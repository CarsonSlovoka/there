package there

import "strings"

type Path struct {
	parts      []PathPart
	ignoreCase bool
}

type PathPart struct {
	value    string
	variable bool
}

func ConstructPath(pathString string, ignoreCase bool) Path {

	parts := make([]PathPart, 0)

	split := splitUrl(pathString)

	for _, s := range split {
		if strings.HasPrefix(s, ":") {
			s = strings.TrimPrefix(s, ":")
			for _, part := range parts {
				if part.variable && part.value == s {
					panic(pathString + " has defined the route param \"" + s + "\" more than once")
				}
			}
			parts = append(parts, PathPart{
				value:    s,
				variable: true,
			})
		} else {
			parts = append(parts, PathPart{
				value:    s,
				variable: false,
			})
		}
	}

	path := Path{
		parts:      parts,
		ignoreCase: ignoreCase,
	}
	return path
}

func (p Path) ToString() string {
	path := "/"
	for i, part := range p.parts {
		if part.variable {
			path += ":"
		}
		path += part.value
		if i != len(p.parts)-1 {
			path += "/"
		}
	}
	return path
}

func (p Path) Equals(toCompare Path) bool {
	if len(p.parts) != len(toCompare.parts) || p.ignoreCase != toCompare.ignoreCase {
		return false
	}
	if len(p.parts) == 0 {
		return true
	}

	ignoreCase := p.ignoreCase

	for i := 0; i < len(p.parts); i++ {
		a := p.parts[i]
		b := toCompare.parts[i]
		if !a.variable && !b.variable {
			if (ignoreCase && strings.ToLower(a.value) != strings.ToLower(b.value)) ||
				(!ignoreCase && a.value != b.value) {
				return false
			}
		} else if (!a.variable && b.variable) || (a.variable && !b.variable) {
			return false
		}
	}

	return true
}

func (p Path) Parse(route string) (map[string]string, bool) {
	params := map[string]string{}

	split := splitUrl(route)

	if len(split) != len(p.parts) {
		return nil, false
	}

	ignoreCase := p.ignoreCase

	for i := 0; i < len(p.parts); i++ {
		a := p.parts[i]
		b := split[i]
		if a.variable {
			params[a.value] = b
		} else {
			if (ignoreCase && strings.ToLower(a.value) != strings.ToLower(b)) ||
				(!ignoreCase && a.value != b) {
				return nil, false
			}
		}
	}

	return params, true
}

func splitUrl(route string) []string {
	for strings.Contains(route, "//") {
		route = strings.ReplaceAll(route, "//", "/")
	}

	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")

	if len(route) == 0 {
		return []string{}
	}

	return strings.Split(route, "/")
}

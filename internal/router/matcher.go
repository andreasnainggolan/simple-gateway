package router

import "strings"

type PathKind int

const (
	KindExact PathKind = iota
	KindParam
	KindWildcard
)

type PathPattern struct {
	Raw   string
	Kind  PathKind
	Parts []string // split by "/"
}

// CompilePath mengubah string path dari YAML menjadi pola internal.
// Kita sengaja tidak pakai regex di V1 agar tetap simple dan aman.
func CompilePath(p string) PathPattern {
	p = strings.TrimSpace(p)
	if p == "" {
		p = "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	// wildcard sederhana hanya di akhir, misal "/assets/*" atau "/*"
	if strings.HasSuffix(p, "/*") || p == "/*" {
		parts := splitParts(p)
		return PathPattern{Raw: p, Kind: KindWildcard, Parts: parts}
	}

	// parameter dengan {} misal "/users/{id}"
	if strings.Contains(p, "{") && strings.Contains(p, "}") {
		parts := splitParts(p)
		return PathPattern{Raw: p, Kind: KindParam, Parts: parts}
	}

	// selain itu exact match
	return PathPattern{Raw: p, Kind: KindExact, Parts: splitParts(p)}
}

func splitParts(p string) []string {
	// "/users/123" => ["users","123"]
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}

// MatchPath mengecek apakah requestPath cocok dengan pattern.
// Return (ok, params). Untuk wildcard params kosong.
func MatchPath(pattern PathPattern, requestPath string) (bool, map[string]string) {
	reqParts := splitParts(requestPath)

	switch pattern.Kind {
	case KindExact:
		if len(pattern.Parts) != len(reqParts) {
			return false, nil
		}
		for i := range pattern.Parts {
			if pattern.Parts[i] != reqParts[i] {
				return false, nil
			}
		}
		return true, nil

	case KindWildcard:
		// "/assets/*" cocok dengan "/assets/a/b/c"
		// aturan: prefix (kecuali "*") harus sama
		if len(pattern.Parts) == 0 {
			// "/*" => match semua
			return true, nil
		}
		// pattern last part is "*" (karena kita compile hanya /*)
		prefix := pattern.Parts[:len(pattern.Parts)-1]
		if len(reqParts) < len(prefix) {
			return false, nil
		}
		for i := range prefix {
			if prefix[i] != reqParts[i] {
				return false, nil
			}
		}
		return true, nil

	case KindParam:
		// jumlah segmen harus sama untuk param pattern
		if len(pattern.Parts) != len(reqParts) {
			return false, nil
		}
		params := map[string]string{}
		for i := range pattern.Parts {
			pp := pattern.Parts[i]
			rp := reqParts[i]

			// {id} berarti "apa saja", tapi kita simpan nilainya
			if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
				name := strings.TrimSuffix(strings.TrimPrefix(pp, "{"), "}")
				name = strings.TrimSpace(name)
				if name == "" {
					return false, nil // invalid param name
				}
				params[name] = rp
				continue
			}

			// segmen biasa harus sama persis
			if pp != rp {
				return false, nil
			}
		}
		return true, params
	default:
		return false, nil
	}
}

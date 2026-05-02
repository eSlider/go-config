package ini

import (
	"strings"

	iniv1 "gopkg.in/ini.v1"
)

func insertPath(root map[string]any, path []string, value string) {
	current := root
	last := len(path) - 1
	for i, raw := range path {
		k := strings.ToLower(raw)
		if i == last {
			if _, exists := current[k]; exists {
				return
			}
			current[k] = value
			return
		}
		next, ok := current[k].(map[string]any)
		if !ok {
			next = make(map[string]any)
			current[k] = next
		}
		current = next
	}
}

func mapFromINIFile(f *iniv1.File) map[string]any {
	out := make(map[string]any)
	for _, name := range f.SectionStrings() {
		sec, err := f.GetSection(name)
		if err != nil {
			continue
		}
		if strings.EqualFold(name, iniv1.DefaultSection) {
			for k, v := range sec.KeysHash() {
				insertPath(out, []string{k}, v)
			}
			continue
		}
		path := strings.Split(name, ".")
		for k, v := range sec.KeysHash() {
			insertPath(out, append(append([]string{}, path...), k), v)
		}
	}
	return out
}

// render/funcs.go
package render

import (
	"fmt"
	"html/template"
)

/*
Extra template helpers that can be used from *.gohtml / *.html files.
Currently we expose only “dict”, which lets you build a map inline:

	{{ template "card" (dict "Title" "Hi" "Body" "Lorem…") }}
*/
var funcs = template.FuncMap{
	"dict": func(values ...any) (map[string]any, error) {
		if len(values)%2 != 0 {
			return nil, fmt.Errorf("dict expects even arg count")
		}
		m := make(map[string]any, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			k, ok := values[i].(string)
			if !ok {
				return nil, fmt.Errorf("dict key %d not string", i)
			}
			m[k] = values[i+1]
		}
		return m, nil
	},
}

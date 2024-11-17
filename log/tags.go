package log

import "strings"

type tags map[string]bool

func (t *tags) GetTag(tag string) (msg string, show bool) {
	if tag == "" {
		show = true
		return
	}

	if len(*t) > 0 {
		if _, ok := (*t)[tag]; !ok {
			return
		}
	}

	msg = "[Tag:" + tag + "] "
	show = true

	return
}

func parseTags(str string) (ts tags) {
	if len(str) == 0 {
		return tags{}
	}

	ts = make(tags)
	arr := strings.Split(str, ",")
	for _, v := range arr {
		ts[v] = true
	}

	return
}

package convert

import "net/url"

func GetUrlParam(s, param string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	m, _ := url.ParseQuery(u.RawQuery)
	if m[param] != nil && len(m[param]) > 0 {
		return m[param][0]
	}
	return ""
}

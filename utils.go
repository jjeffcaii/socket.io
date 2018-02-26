package sio

import "regexp"

var regexNsp = regexp.MustCompile("/[0-9a-zA-Z_-]*")

func isValidNamespace(nsp string) bool {
	return regexNsp.MatchString(nsp)
}

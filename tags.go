/*
Copyright 2020 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package httphelper

// First/Last part tags for classifying files during load.
// Hardcoded Tags: Resource
var TagsFirst = map[string][]string{
	".go":     {"Go"},
	".static": {"Static"},
}
var TagsLast = map[string][]string{
	".htm":  {"HTML"},
	".html": {"HTML"},
	".css":  {"StyleSheet"},
	".js":   {"JavaScript"},
}

// GetFileTags finds the file tags for a file with the given name.
// The returned slice of tags is yours to keep.
func GetFileTags(name string) []string {
	f, l := getExtParts(name)

	var fv, lv []string
	if t, ok := TagsFirst[f]; ok {
		fv = t
	}
	if t, ok := TagsLast[l]; ok {
		lv = t
	}

	rtn := make([]string, 0, len(fv)+len(lv))
	rtn = append(rtn, fv...)
	rtn = append(rtn, lv...)
	return rtn
}

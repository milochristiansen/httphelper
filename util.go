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

// This code is lifted almost unchanged from Rubble 8, so it has a lot of advanced stuff I don't actually need.

import "strings"

// replacePrefix checks to see if the string begins with the prefix, and replaces it if it does.
// This is generally used to "fix" AXIS paths before AXIS is available.
func replacePrefix(s, prefix, replace string) string {
	if strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
		return replace + s
	}
	return s
}

// Strip the extension from a file name.
// If a file has multiple extensions strip only the last.
func stripExt(name string) string {
	i := len(name) - 1
	for i >= 0 {
		if name[i] == '.' {
			return name[:i]
		}
		i--
	}
	return name
}

// replaceExt replaces the file extension with n, but only if it matches o.
func replaceExt(name, o, n string) string {
	if strings.HasSuffix(name, o) {
		return strings.TrimSuffix(name, o) + n
	}
	return name
}

// replaceExtAdv replaces a two part extension.
// If you set one of the parts of the old extension to ".%" this will assume any extension may go in that place.
// If n is set to "" then it will use the old last part extension.
func replaceExtAdv(name, o, n string) string {
	of, ol := getExtParts(o)
	f, l := getExtParts(name)
	if of == ".%" {
		of = f
	} else if of != f {
		return name
	}
	if ol == ".%" {
		ol = l
	} else if ol != l {
		return name
	}
	if n == "" {
		n = ol
	}
	return strings.TrimSuffix(name, of+ol) + n
}

// getExt returns the extension from a file name.
func getExt(name string) string {
	// Find the last part of the extension
	i := len(name) - 1
	for i >= 0 {
		if name[i] == '.' {
			return name[i:]
		}
		i--
	}
	return ""
}

// getExtParts returns the extension from a file name.
// Unlike GetExt, GetExtParts returns the first and last part of a two part extension separately. `"abc.x.y"` would
// return `".x", ".y"` and `"abc.d"` would return `"", ".d"`.
func getExtParts(name string) (first string, last string) {
	// Find the last part of the extension
	i := len(name) - 1
	j := 0
	for i >= 0 {
		if name[i] == '.' {
			last = name[i:]
			j = i
			i--
			break
		}
		i--
	}
	// Then look for the first part.
	for i >= 0 {
		if name[i] == '.' {
			return name[i : i+(j-i)], last
		}
		i--
	}
	return "", last
}

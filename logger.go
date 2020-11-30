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

// Loggers passed into this package need to satisfy this interface.
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type logger struct {
	i, w, e Logger
}

func (l *logger) InfoPrint(v ...interface{}) {
	if l.i != nil {
		l.i.Print(v...)
	}
}

func (l *logger) WarnPrint(v ...interface{}) {
	if l.w != nil {
		l.w.Print(v...)
	}
}

func (l *logger) ErrPrint(v ...interface{}) {
	if l.e != nil {
		l.e.Print(v...)
	}
}

func (l *logger) InfoPrintf(format string, v ...interface{}) {
	if l.i != nil {
		l.i.Printf(format, v...)
	}
}

func (l *logger) WarnPrintf(format string, v ...interface{}) {
	if l.w != nil {
		l.w.Printf(format, v...)
	}
}

func (l *logger) ErrPrintf(format string, v ...interface{}) {
	if l.e != nil {
		l.e.Printf(format, v...)
	}
}

func (l *logger) InfoPrintln(v ...interface{}) {
	if l.i != nil {
		l.i.Println(v...)
	}
}

func (l *logger) WarnPrintln(v ...interface{}) {
	if l.w != nil {
		l.w.Println(v...)
	}
}

func (l *logger) ErrPrintln(v ...interface{}) {
	if l.e != nil {
		l.e.Println(v...)
	}
}

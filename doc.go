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

// HTTP Helper: Make servers easy again!
//
// HTTP Helper is a simple library to make it easier to create simple web servers that need to provide a mix of
// static, template, and generated content. This is not suitable for large websites, as all content is loaded and
// served from memory.
//
// When the server is initialized all data files are loaded and classified. Once classification is done and all
// handlers are assigned their requested resources, anything left over is served as static content.
package httphelper

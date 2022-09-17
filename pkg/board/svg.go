/*
 * wars of warp - an implementation of warpwar
 *
 * Copyright (c) 2022 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package board

import "fmt"

// svg is the container for our board
type svg struct {
	id      string
	viewBox struct {
		minX, minY    int
		width, height int
	}
	polygons []*polygon
}

func (s svg) String() string {
	t := "<svg"
	if s.id != "" {
		t += fmt.Sprintf(" id=%q", s.id)
	}
	t += fmt.Sprintf(` width="%d" height="%d"`, s.viewBox.width+40, s.viewBox.height+40)
	t += fmt.Sprintf(` viewBox="%d %d %d %d"`, s.viewBox.minX, s.viewBox.minY, s.viewBox.width+40, s.viewBox.height+40)
	t += ` xmlns="http://www.w3.org/2000/svg">`
	for _, p := range s.polygons {
		t += fmt.Sprintf("\n%s", p.String())
	}
	return t + "\n</svg>"
}

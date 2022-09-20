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

import (
	"fmt"
)

// svg is the container for our board
type svg struct {
	id      string
	viewBox struct {
		minX, minY    int
		width, height int
	}
	hexes    []*polygon
	polygons []*polygon
	lines    [][4]float64
}

func (s svg) String() string {
	fontSize := 14
	t := "<svg"
	if s.id != "" {
		t += fmt.Sprintf(" id=%q", s.id)
	}
	t += fmt.Sprintf(` width="%d" height="%d"`, s.viewBox.width+40, s.viewBox.height+40)
	t += fmt.Sprintf(` viewBox="%d %d %d %d"`, s.viewBox.minX, s.viewBox.minY, s.viewBox.width+40, s.viewBox.height+40)
	t += ` xmlns="http://www.w3.org/2000/svg">`
	for _, h := range s.hexes {
		if len(h.points) == 0 {
			continue
		}
		p := fmt.Sprintf(`<polygon style="fill: %s; stroke: %s; stroke-width: %s;"`, h.style.fill, h.style.stroke, h.style.strokeWidth)
		p += fmt.Sprintf(` points="`)
		for i, pt := range h.points {
			if i > 0 {
				p += " "
			}
			p += pt.String()
		}
		p += `"`
		p += "></polygon>\n"
		p += fmt.Sprintf(`<text x="%f" y="%f" text-anchor="middle" fill="grey" font-size="14" font-weight="bold">%s</text>`, h.cx, h.cy, fmt.Sprintf("%02d%02d", h.col, h.row))
		t += p
	}
	for _, l := range s.lines {
		x1, y1, x2, y2 := l[0], l[1], l[2], l[3]
		t += fmt.Sprintf(`<line x1="%f" y1="%f" x2="%f" y2="%f" stroke-width="2" stroke="black"/>`, x1, y1, x2, y2)
	}
	for _, p := range s.polygons {
		ps := fmt.Sprintf(`<circle cx="%f" cy="%f" r="%f" style="fill: %s; stroke: %s; stroke-width: %s" />`, p.cx, p.cy, p.radius*0.88, p.style.fill, p.style.stroke, p.style.strokeWidth) + "\n"
		t += ps
		//// todo: put in a rounded rectangle behind the text
		//rbHeight, rbWidth := p.radius, p.radius*1.8
		//s += fmt.Sprintf(`<rect x="%f" y="%f" height="%f" width="%f" rx="%f" ry="%f" fill="white" />`, p.cx-rbWidth/2.0, p.cy-rbHeight/2.0, rbHeight, rbWidth, rbHeight/2.0, rbHeight/2.0)
		//t += fmt.Sprintf(`<text x="%f" y="%f" text-anchor="middle" fill="black" font-size="%d" font-weight="bold">%s</text>`, p.cx, p.cy, fontSize, p.label)
		yOffset := float64(fontSize) * 0.6
		for i, text := range p.text {
			if i == 0 {
				t += fmt.Sprintf(`<text x="%f" y="%f" text-anchor="middle" fill="black" font-size="%d" font-weight="bold">%s</text>`, p.cx, p.cy-yOffset, fontSize, text)
			} else {
				t += fmt.Sprintf(`<text x="%f" y="%f" text-anchor="middle" fill="black" font-size="%d" font-weight="bold">%s</text>`, p.cx, p.cy+yOffset*3, fontSize+2, text)
			}
		}
	}
	return t + "\n</svg>"
}

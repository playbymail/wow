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
	"math"
)

// polygon is the actual hex on the board
type polygon struct {
	x, y, radius float64
	label        string
	style        struct {
		fill        string
		stroke      string
		strokeWidth string
	}
	points []point
}

func (p polygon) hexPoints() (points []point) {
	for theta := 0.0; theta < math.Pi*2.0; theta += math.Pi / 3.0 {
		points = append(points, point{x: p.x + p.radius*math.Sin(theta), y: p.y + p.radius*math.Cos(theta)})
	}
	return points
}

func (p polygon) String() string {
	s := fmt.Sprintf(`<polygon style="fill: %s; stroke: %s; stroke-width: %s;"`, p.style.fill, p.style.stroke, p.style.strokeWidth)
	if len(p.points) != 0 {
		s += fmt.Sprintf(` points="`)
		for i, pt := range p.points {
			if i != 0 {
				s += " "
			}
			s += pt.String()
		}
		s += `"`
	}
	s += "></polygon>\n"
	s += fmt.Sprintf(`<text x="%f" y="%f" text-anchor="middle" fill="grey" font-size="12">%s</text>`, p.x, p.y, p.label)
	return s
}

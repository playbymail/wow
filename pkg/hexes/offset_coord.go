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

package hexes

type OFFSET int

const (
	EVEN OFFSET = 1
	ODD  OFFSET = -1
)

type OffsetCoord struct {
	col, row int
}

func NewOffsetCoord(col, row int) OffsetCoord {
	return OffsetCoord{col: col, row: row}
}

func (a OffsetCoord) Equals(b OffsetCoord) bool {
	return a.col == b.col && a.row == b.row
}

func qoffset_from_cube(offset OFFSET, h Hex) OffsetCoord {
	col := h.q
	row := h.r + (h.q+int(offset)*(h.q&1))/2

	return NewOffsetCoord(col, row)
}

func qoffset_to_cube(offset OFFSET, h OffsetCoord) Hex {
	q := h.col
	r := h.row - (h.col+int(offset)*(h.col&1))/2
	s := -q - r

	return NewHex(q, r, s)
}

func roffset_from_cube(offset OFFSET, h Hex) OffsetCoord {
	col := h.q + (h.r+int(offset)*(h.r&1))/2
	row := h.r

	return NewOffsetCoord(col, row)
}

func roffset_to_cube(offset OFFSET, h OffsetCoord) Hex {
	q := h.col - ((h.row + int(offset)*(h.row&1)) / 2)
	r := h.row
	s := -q - r

	return NewHex(q, r, s)
}

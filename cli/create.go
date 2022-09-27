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

package cli

import (
	"github.com/mdhender/wow/pkg/board"
	"github.com/spf13/cobra"
	"os"
)

// cmdCreateMap creates a map
var cmdCreateMap = &cobra.Command{
	Use:   "create",
	Short: "create a new map",
	Run: func(cmd *cobra.Command, args []string) {
		type hex struct {
			name      string
			col, row  int
			econValue int // non-zero only if hasStar
			warps     []string
		}

		// hexes is an untested guess as to the layout of the original map
		hexes := []hex{
			{"Adab", 6, 6, 0, []string{"Erech", "Khafa", "Byblos"}},
			{"Akkad", 7, 16, 3, []string{"Kish"}},
			{"Assur", 12, 10, 2, []string{"Nippur", "Lagash"}},
			{"Babylon", 6, 18, 4, []string{"Sumer"}},
			{"Byblos", 2, 6, 3, []string{"Adab"}},
			{"Calah", 8, 4, 1, []string{"Nippur"}},
			{"Elam", 7, 12, 5, []string{"Lagash"}},
			{"Erech", 4, 4, 3, []string{"Ur", "Adab"}},
			{"Eridu", 12, 16, 1, []string{"Kish", "Ugarit"}},
			{"Girsu", 8, 13, 1, []string{"Umma"}},
			{"Jarmo", 11, 12, 3, []string{"Kish"}},
			{"Isin", 1, 15, 1, []string{"Nineveh"}},
			{"Khafa", 7, 9, 2, []string{"Adab"}},
			{"Kish", 10, 15, 0, []string{"Jarmo", "Eridu"}},
			{"Lagash", 9, 11, 1, []string{"Assur"}},
			{"Larsu", 11, 2, 2, []string{"Susa"}},
			{"Mari", 6, 10, 1, []string{"Ubaid", "Umma"}},
			{"Mosul", 3, 1, 2, []string{"Sippur"}},
			{"Nineveh", 3, 19, 2, []string{"Isin"}},
			{"Nippur", 10, 7, 1, []string{"Calah", "Susa", "Assur", "Lagash"}},
			{"Sippur", 2, 4, 1, []string{"Mosul"}},
			{"Sumarra", 2, 12, 2, []string{"Ubaid", "Umma"}},
			{"Sumer", 4, 16, 0, []string{"Umma", "Babylon"}},
			{"Susa", 12, 5, 0, []string{"Larsu", "Nippur"}},
			{"Ubaid", 3, 8, 5, []string{"Mari", "Sumarra"}},
			{"Ugarit", 11, 20, 2, []string{"Eridu"}},
			{"Umma", 5, 14, 2, []string{"Sumarra", "Mari", "Girsu", "Sumer"}},
			{"Ur", 7, 2, 4, []string{"Erech"}},
		}

		// max row and col determine the size of the board
		maxRow, maxCol := 0, 0
		for _, h := range hexes {
			if h.row > maxRow {
				maxRow = h.row
			}
			if h.col > maxCol {
				maxCol = h.col
			}
		}

		// create the board, add all the stars, then add the wormholes
		gb := board.NewBoard(maxRow, maxCol)
		for _, h := range hexes {
			gb.AddStar(h.name, h.row, h.col, h.econValue)
		}
		for _, h := range hexes {
			for _, target := range h.warps {
				gb.AddWormHole(h.name, target)
			}
		}

		mono := true
		// save the board as an SVG file
		cobra.CheckErr(os.WriteFile("svg-mono-test.svg", gb.AsSVG(mono), 0644))
		cobra.CheckErr(os.WriteFile("svg-test.svg", gb.AsSVG(!mono), 0644))

		// save the board as an HTML file
		cobra.CheckErr(os.WriteFile("svg-mono-test.html", gb.AsHTML(mono), 0644))
		cobra.CheckErr(os.WriteFile("svg-test.html", gb.AsHTML(!mono), 0644))
	},
}

func init() {
	cmdBase.AddCommand(cmdCreateMap)
}

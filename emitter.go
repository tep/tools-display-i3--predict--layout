// Copyright © 2018 Timothy E. Peoples <eng@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"

	"go.i3wm.org/i3"
)

const (
	dfltPeerColor    = "#6F6"
	dfltInsideColor  = "#F66"
	dfltOutsideColor = "#FF6"

	dfltStackedPeer   = "" // \uF8CB
	dfltTabbedPeer    = "" // \uF465
	dfltSplithInside  = "ﲒ" // \uFC92
	dfltSplithOutside = "ﲖ" // \uFC96
	dfltSplitvInside  = "ﲐ" // \uFC90
	dfltSplitvOutside = "ﲔ" // \uFC94

	dfltNumberIcons = ""
)

type emitter struct {
	fontNum int

	colors struct {
		peer    string
		inside  string
		outside string
	}

	icons struct {
		stackedPeer   icon
		tabbedPeer    icon
		splithInside  icon
		splithOutside icon
		splitvInside  icon
		splitvOutside icon
	}
}

func (e *emitter) emit(s *stack) {
	var out string

	ni := []rune(dfltNumberIcons)

	if cnt := s.peers() + 1; cnt < 10 {
		out += string(ni[cnt])
	} else {
		out += string(ni[10])
	}

	lo, pl := s.nextLayout()

	if i := e.icon(lo, pl); i != "" {
		out += fmt.Sprintf(" %s ", i)
	}

	if c := e.color(pl); c != "" {
		out = fmt.Sprintf("%%{F%s}%s%%{F-}", c, out)
	}

	if e.fontNum != 0 {
		out = fmt.Sprintf("%%{T%d}%s%%{T-}", e.fontNum, out)
	}

	fmt.Println(out)
}

func (e *emitter) icon(lo i3.Layout, pl placement) string {
	switch {
	case lo == i3.Stacked && pl == peer:
		return e.icons.stackedPeer.String()
	case lo == i3.Tabbed && pl == peer:
		return e.icons.tabbedPeer.String()
	case lo == i3.SplitH && pl == inside:
		return e.icons.splithInside.String()
	case lo == i3.SplitH && pl == outside:
		return e.icons.splithOutside.String()
	case lo == i3.SplitV && pl == inside:
		return e.icons.splitvInside.String()
	case lo == i3.SplitV && pl == outside:
		return e.icons.splitvOutside.String()
	default:
		return ""
	}
}

func (e *emitter) color(pl placement) string {
	switch pl {
	case peer:
		return e.colors.peer
	case inside:
		return e.colors.inside
	case outside:
		return e.colors.outside
	default:
		return ""
	}
}

type icon string

func (i *icon) ptr() *string {
	return (*string)(i)
}

func (i icon) String() string {
	if len(i) == 0 {
		return ""
	}

	return string([]rune(i)[0])
}

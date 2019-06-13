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
	"strings"

	"go.i3wm.org/i3"
)

type stack struct {
	output string
	window int64
	nodes  []*i3.Node
}

func windowStack(window int64) (*stack, error) {
	t, err := i3.GetTree()
	if err != nil {
		return nil, err
	}

	find := t.Root.FindFocused
	if window != 0 {
		find = t.Root.FindChild
	}

	s := new(stack)

	find(func(n *i3.Node) bool {
		if n.Type == i3.OutputNode {
			s.output = n.Name
		}

		if n.Type == i3.WorkspaceNode || len(s.nodes) > 0 {
			s.nodes = append([]*i3.Node{n}, s.nodes...)
		}

		if window != 0 {
			if window == n.Window {
				s.window = n.Window
				return true
			}
			return false
		}

		return n.Focused
	})

	return s, nil
}

func (s *stack) String() string {
	ss := make([]string, len(s.nodes))
	for i, n := range s.nodes {
		ss[i] = fmt.Sprintf("%s:%s", n.Type, n.Layout)
	}
	return strings.Join(ss, " -> ")
}

func (s *stack) ancestorLayout(i int) i3.Layout {
	if li := len(s.nodes) - 1; i <= li {
		return s.nodes[i].Layout
	}
	return ""
}

func (s *stack) isFloating() bool {
	if len(s.nodes) < 2 {
		return false
	}

	return s.nodes[1].Type == i3.FloatingCon
}

func (s *stack) peers() int {
	if len(s.nodes) < 2 {
		return 0
	}
	return len(s.nodes[1].Nodes) - 1
}

func (s *stack) nextLayout() (i3.Layout, placement) {
	var (
		pl = s.ancestorLayout(1)
		np = outside
	)

	switch pl {
	case i3.Stacked, i3.Tabbed:
		np = peer

	case i3.SplitH, i3.SplitV:
		if gpl := s.ancestorLayout(2); s.peers() == 0 || gpl == i3.Stacked || gpl == i3.Tabbed {
			np = inside
		}

	default:
		if cl := s.ancestorLayout(0); cl != "" {
			return cl, np
		}
		return "", ""
	}

	return pl, np
}

func (s stack) nextLayoutString() string {
	if l, p := s.nextLayout(); l != "" && p != "" {
		return fmt.Sprintf("%s-%s", l, p)
	}
	return ""
}

type placement string

const (
	peer    placement = "peer"
	inside  placement = "inside"
	outside placement = "outside"
)

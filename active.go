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

import "go.i3wm.org/i3"

func activeOutputNode(output string) (*i3.Node, error) {
	ws, err := i3.GetWorkspaces()
	if err != nil {
		return nil, err
	}

	t, err := i3.GetTree()
	if err != nil {
		return nil, err
	}

	for _, w := range ws {
		if w.Output == output && w.Visible {
			return activeNode(
				t.Root.FindChild(func(n *i3.Node) bool {
					return n.Type == i3.WorkspaceNode && n.Name == w.Name
				})), nil
		}
	}

	return nil, nil
}

func activeNode(in *i3.Node) *i3.Node {
	if in == nil || len(in.Focus) == 0 {
		return in
	}

	for _, n := range in.Nodes {
		if n.ID == in.Focus[0] {
			return activeNode(n)
		}
	}

	return nil
}

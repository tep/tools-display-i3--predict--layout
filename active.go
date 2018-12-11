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

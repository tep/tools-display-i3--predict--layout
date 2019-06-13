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
	"os"
	"strings"

	"github.com/spf13/pflag"

	"go.i3wm.org/i3"

	"toolman.org/base/log/v2"
	"toolman.org/base/toolman/v2"
)

var (
	output = pflag.StringP("output", "o", "", "Only process events for the given output")
	oenv   = pflag.String("output-env", "", "Process events only for output specified in given environment variable")
)

var emitr emitter

func init() {
	pflag.IntVar(&emitr.fontNum, "font-num", 0, "Polybar font number")

	pflag.StringVar(&emitr.colors.peer, "peer-color", dfltPeerColor, "Color for Stacked or Tabbed Peers")
	pflag.StringVar(&emitr.colors.inside, "inside-color", dfltInsideColor, "Color for Inside SplitH or SplitV")
	pflag.StringVar(&emitr.colors.outside, "outside-color", dfltOutsideColor, "Color for Outside SplitH or SplitV")

	pflag.StringVar(emitr.icons.stackedPeer.ptr(), "stacked-peer-icon", dfltStackedPeer, "Icon for Stacked Peer Prediction")
	pflag.StringVar(emitr.icons.tabbedPeer.ptr(), "tabbed-peer-icon", dfltTabbedPeer, "Icon for Tabbed Peer Prediction")

	pflag.StringVar(emitr.icons.splithInside.ptr(), "splith-inside-icon", dfltSplithInside, "Icon for SplitH Inside Prediction")
	pflag.StringVar(emitr.icons.splithOutside.ptr(), "splith-outside-icon", dfltSplithOutside, "Icon for SplitH Outside Prediction")

	pflag.StringVar(emitr.icons.splitvInside.ptr(), "splitv-inside-icon", dfltSplitvInside, "Icon for SplitV Inside Prediction")
	pflag.StringVar(emitr.icons.splitvOutside.ptr(), "splitv-outside-icon", dfltSplitvOutside, "Icon for SplitV Outside Prediction")
}

func main() {
	toolman.Init(toolman.Quiet())

	var oname string

	switch {
	case *output != "" && *oenv != "":
		log.Fatal("Only one of --output or --output-env may be specified")
	case *output != "":
		oname = *output
	case *oenv != "":
		var ok bool
		if oname, ok = os.LookupEnv(*oenv); !ok {
			log.Fatalf("No such environment variable: %s", *oenv)
		}
	}

	if err := run(oname); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
}

func run(outputName string) error {
	rcvr := i3.Subscribe(i3.WindowEventType, i3.BindingEventType)
	defer rcvr.Close()

	var win int64

	if outputName != "" {
		n, err := activeOutputNode(outputName)
		if err != nil {
			return err
		}
		win = n.Window
		log.Infof("Initial active window for output %q: %d", outputName, win)
	}

	if err := predictLayout(outputName, win); err != nil {
		log.Error(err)
	}

	for rcvr.Next() {
		if err := eventProc(outputName, rcvr.Event()); err != nil {
			log.Error(err)
		}
	}

	return nil
}

func eventProc(outputName string, event i3.Event) error {
	// TODO(tep): Deal with pointer on desktop
	switch evt := event.(type) {
	case *i3.WindowEvent:
		if evt.Change != "focus" {
			return nil
		}

	case *i3.BindingEvent:
		if evt.Change != "run" {
			return nil
		}

		flds := strings.Fields(evt.Binding.Command)
		if len(flds) == 0 {
			return nil
		}

		switch flds[0] {
		case "split", "layout", "focus":
			// pass
		default:
			return nil
		}

	default:
		return nil
	}

	return predictLayout(outputName, 0)
}

func predictLayout(oname string, window int64) error {
	log.V(1).Infof("predictLayout(%q, %d)", oname, window)

	s, err := windowStack(window)
	if err != nil {
		return err
	}

	if oname != "" && s.output != oname {
		log.V(1).Infof("Skipping active window %s:%d != %s:%d", s.output, s.window, oname, window)
		return nil
	}

	if s.isFloating() {
		log.Infof("Skipping floating container for window: %d", s.window)
		return nil
	}

	log.Infof("EVT: win=%d  [peers=%2d] %-75s  [next=%s]\n", s.window, s.peers(), s, s.nextLayoutString())
	emitr.emit(s)

	return nil
}

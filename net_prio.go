// +build linux

package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func NewNetPrio(root string) *netprioController {
	return &netprioController{
		root: filepath.Join(root, string(NetPrio)),
	}
}

type netprioController struct {
	root string
}

func (n *netprioController) Name() Name {
	return NetPrio
}

func (n *netprioController) Path(path string) string {
	return filepath.Join(n.root, path)
}

func (n *netprioController) Create(path string, resources *specs.LinuxResources) error {
	if err := os.MkdirAll(n.Path(path), defaultDirPerm); err != nil {
		return err
	}
	if resources.Network != nil {
		for _, prio := range resources.Network.Priorities {
			if err := ioutil.WriteFile(
				filepath.Join(n.Path(path), "net_prio_ifpriomap"),
				formatPrio(prio.Name, prio.Priority),
				defaultFilePerm,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func formatPrio(name string, prio uint32) []byte {
	return []byte(fmt.Sprintf("%s %d", name, prio))
}

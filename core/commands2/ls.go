package commands

import (
	"fmt"

	cmds "github.com/jbenet/go-ipfs/commands"
	"github.com/jbenet/go-ipfs/core/commands2/internal"
	merkledag "github.com/jbenet/go-ipfs/merkledag"
)

type Link struct {
	Name, Hash string
	Size       uint64
}

type Object struct {
	Hash  string
	Links []Link
}

type LsOutput struct {
	Objects []Object
}

var lsCmd = &cmds.Command{
	Arguments: []cmds.Argument{
		cmds.Argument{"object", cmds.ArgString, false, true},
	},
	// TODO UsageLine: "ls",
	// TODO Short:     "List links from an object.",
	// TODO docs read ipfs-path. argument says option. which?
	Help: `ipfs ls <ipfs-path> - List links from an object.

    Retrieves the object named by <ipfs-path> and displays the links
    it contains, with the following format:

    <link base58 hash> <link size in bytes> <link name>

`,
	Run: func(res cmds.Response, req cmds.Request) {
		node := req.Context().Node

		paths, err := internal.ToStrings(req.Arguments())
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		dagnodes := make([]*merkledag.Node, 0)
		for _, path := range paths {
			dagnode, err := node.Resolver.ResolvePath(path)
			if err != nil {
				res.SetError(err, cmds.ErrNormal)
				return
			}
			dagnodes = append(dagnodes, dagnode)
		}

		output := make([]Object, len(req.Arguments()))
		for i, dagnode := range dagnodes {
			output[i] = Object{
				Hash:  paths[i],
				Links: make([]Link, len(dagnode.Links)),
			}
			for j, link := range dagnode.Links {
				output[i].Links[j] = Link{
					Name: link.Name,
					Hash: link.Hash.B58String(),
					Size: link.Size,
				}
			}
		}

		res.SetOutput(&LsOutput{output})
	},
	Marshallers: map[cmds.EncodingType]cmds.Marshaller{
		cmds.Text: func(res cmds.Response) ([]byte, error) {
			s := ""
			output := res.Output().(*LsOutput).Objects

			for _, object := range output {
				if len(output) > 1 {
					s += fmt.Sprintf("%s:\n", object.Hash)
				}

				for _, link := range object.Links {
					s += fmt.Sprintf("-> %s %s (%v bytes)\n", link.Name, link.Hash, link.Size)
				}

				if len(output) > 1 {
					s += "\n"
				}
			}

			return []byte(s), nil
		},
	},
	Type: &LsOutput{},
}

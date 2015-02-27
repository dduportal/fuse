// Copyright 2015 Google Inc. All Rights Reserved.
// Author: jacobsa@google.com (Aaron Jacobs)

package samples

import (
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseutil"
	"github.com/jacobsa/gcsfuse/timeutil"
	"golang.org/x/net/context"
)

// A file system with a fixed structure that looks like this:
//
//     hello
//     dir/
//         world
//
// Each file contains the string "Hello, world!".
type HelloFS struct {
	fuseutil.NotImplementedFileSystem
	Clock timeutil.Clock
}

var _ fuse.FileSystem = &HelloFS{}

const (
	rootInode fuse.InodeID = fuse.RootInodeID + iota
	helloInode
	dirInode
	worldInode
)

func (fs *HelloFS) OpenDir(
	ctx context.Context,
	req *fuse.OpenDirRequest) (resp *fuse.OpenDirResponse, err error) {
	// We always allow opening the root directory.
	if req.Inode == rootInode {
		resp = &fuse.OpenDirResponse{}
		return
	}

	// TODO(jacobsa): Handle others.
	err = fuse.ENOSYS
	return
}

// We have a fixed directory structure.
var gDirectoryEntries = map[fuse.InodeID][]fuseutil.Dirent{
	// root
	rootInode: []fuseutil.Dirent{
		fuseutil.Dirent{
			Offset: 1,
			Inode:  helloInode,
			Name:   "hello",
			Type:   fuseutil.DT_File,
		},
		fuseutil.Dirent{
			Offset: 2,
			Inode:  dirInode,
			Name:   "dir",
			Type:   fuseutil.DT_Directory,
		},
	},
}

func (fs *HelloFS) ReadDir(
	ctx context.Context,
	req *fuse.ReadDirRequest) (resp *fuse.ReadDirResponse, err error) {
	resp = &fuse.ReadDirResponse{}

	// Find the entries for this inode.
	entries, ok := gDirectoryEntries[req.Inode]
	if !ok {
		err = fuse.ENOENT
		return
	}

	// Grab the range of interest.
	if req.Offset > fuse.DirOffset(len(entries)) {
		err = fuse.EIO
		return
	}

	entries = entries[req.Offset:]

	// Resume at the specified offset into the array.
	for _, e := range entries {
		resp.Data = fuseutil.AppendDirent(resp.Data, e)
		if len(resp.Data) > req.Size {
			resp.Data = resp.Data[:req.Size]
			break
		}
	}

	return
}

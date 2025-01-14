package fuse

import (
    "context"
    "flag"
    "log"
    "os"
    "syscall"

	"bazil.org/fuse"
    "bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)


func main() {
    var mountpoint string
    flag.StringVar(&mountpoint, "mountpoint", "", "mount point(dir)?")
    flag.Parse()

    if mountpoint == "" {
        log.Fatal("please input invalid mount point\n")
    }
    // 建立一个负责解析和封装 FUSE 请求监听通道对象；
    c, err := fuse.Mount(mountpoint, fuse.FSName("helloworld"), fuse.Subtype("hellofs"))
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // 把 FS 结构体注册到 server，以便可以回调处理请求
    err = fs.Serve(c, FS{})
    if err != nil {
        log.Fatal(err)
    }
}

// hellofs 文件系统的主体
type FS struct{}

func (FS) Root() (fs.Node, error) {
    return Dir{}, nil
}

// hellofs 文件系统中，Dir 是目录操作的主体
type Dir struct{}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
    a.Inode = 20210601
    a.Mode = os.ModeDir | 0555
    return nil
}

// 当 ls 目录的时候，触发的是 ReadDirAll 调用，这里返回指定内容，表明只有一个 hello 的文件；
func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
    // 只处理一个叫做 hello 的 entry 文件，其他的统统返回 not exist
    if name == "hello" {
        return File{}, nil
    }
    return nil, syscall.ENOENT
}

// 定义 Readdir 的行为，固定返回了一个 inode:2 name 叫做 hello 的文件。对应用户的行为一般是 ls 这个目录。
func (Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
    var dirDirs = []fuse.Dirent{{Inode: 2, Name: "hello", Type: fuse.DT_File}}
    return dirDirs, nil
}

// hellofs 文件系统中，File 结构体实现了文件系统中关于文件的调用实现
type File struct{}

const fileContent = "hello, world\n"

// 当 stat 这个文件的时候，返回 inode 为 2，mode 为 444
func (File) Attr(ctx context.Context, a *fuse.Attr) error {
    a.Inode = 20210606
    a.Mode = 0444
    a.Size = uint64(len(fileContent))
    return nil
}

// 当 cat 这个文件的时候，文件内容返回 hello，world
func (File) ReadAll(ctx context.Context) ([]byte, error) {
    return []byte(fileContent), nil
}





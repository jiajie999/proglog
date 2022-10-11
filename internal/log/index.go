package log

import (
	"github.com/tysonmote/gommap"
	"io"
	"os"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth // 单条record的大小
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
	Size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}

	/*
	   这里把文件裁成index的实际大小，newstore里是扩成定义的大小，末尾填零.
	   原因是为了在重启的时候读index的最后一条就是record的最后一条，而不是一个空值。
	*/
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}

func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 { // 空index
		return 0, 0, io.EOF
	}
	if in == -1 {
		out = uint32((i.size / entWidth) - 1) // 返回最后一个record, 如果idx是空的则out=0/width-1=-1，uint32转类型会报错
	} else {
		out = uint32(in)
	}
	//这里的out是要读消息的上一条编号，所以编号*消息长度=目前在index中要读数据的起始位置。
	pos = uint64(out) * entWidth
	println("i.size < pos+entWidth ::::", i.size, "=====", pos+entWidth)
	// 起始位置+消息长度如果大于index大小说明越界了
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}
	//这个out和上面的out其实不是一回事，用个新变量更易读
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])
	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {
	// mmap的size要能放的下现有size+新record
	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	i.size += uint64(entWidth)
	i.Size = i.size
	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}

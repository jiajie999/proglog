package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexSimple(t *testing.T) {

	f, err := ioutil.TempFile(os.TempDir(), "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	idx, err := newIndex(f, c)
	require.NoError(t, err)
	_, _, err = idx.Read(-1)
	require.Error(t, err)
	require.Equal(t, f.Name(), idx.Name())

	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 9},
		{Off: 1, Pos: 10},
		{Off: 2, Pos: 11},
		{Off: 3, Pos: 13},
		{Off: 4, Pos: 16},
		{Off: 5, Pos: 122},
	}

	for _, want := range entries {

		err = idx.Write(want.Off, want.Pos)
		require.NoError(t, err)
		println(idx.Size)
	}

	//_, pos, err := idx.Read(5)
	//
	//println("===")
	//x := []int{0, 1, 2, 3, 4, 5}
	//for i := range x[0:3] {
	//	println(i)
	//
	//}

	println("========")
	println(len([]byte("hello world")))

}

//
//func TestIndex(t *testing.T) {
//	t.Helper()
//	f, err := ioutil.TempFile(os.TempDir(), "index_test")
//	require.NoError(t, err)
//	defer os.Remove(f.Name())
//
//	c := Config{}
//	c.Segment.MaxIndexBytes = 1024
//	idx, err := newIndex(f, c)
//	require.NoError(t, err)
//	_, _, err = idx.Read(-1)
//	require.Error(t, err)
//	require.Equal(t, f.Name(), idx.Name())
//
//	entries := []struct {
//		Off uint32
//		Pos uint64
//	}{
//		{Off: 0, Pos: 0},
//		{Off: 1, Pos: 10},
//		{Off: 2, Pos: 11},
//	}
//
//	for _, want := range entries {
//
//		err = idx.Write(want.Off, want.Pos)
//		require.NoError(t, err)
//
//		_, pos, err := idx.Read(int64(want.Off))
//
//		require.NoError(t, err)
//		require.Equal(t, want.Pos, pos)
//	}
//
//	// index and scanner should error when reading past existing entries
//	println(int64(len(entries)))
//	_, _, err = idx.Read(int64(len(entries)))
//	require.Equal(t, io.EOF, err)
//	_ = idx.Close()
//
//	// index should build its state from the existing file
//	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
//	idx, err = newIndex(f, c)
//	require.NoError(t, err)
//	off, pos, err := idx.Read(-1)
//	require.NoError(t, err)
//	require.Equal(t, uint32(2), off)
//	require.Equal(t, entries[2].Pos, pos)
//}

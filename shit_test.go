package main

import (
	"io"
	"os"
	"testing"
)

// osFS implements fileSystem using the local disk.
type osFS struct{}

var fs fileSystem = osFS{}

type fileSystem interface {
    Open(name string) (file, error)
	Create(name string) (*os.File,error)
	Read(name string) ([]byte, error)
	// Writer(name string) (int, error)
    Stat(name string) (os.FileInfo, error)
}

type file interface {
    io.Closer
    io.Reader
    io.ReaderAt
    io.Seeker
	io.Writer
	io.Reader
    Stat() (os.FileInfo, error)
    // Create() (os.Create, error)

}



func (osFS) Open(name string) (file, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (osFS) Create(name string) (*os.File, error) { return os.Create(name)}
func (osFS) Read(name string) ([]byte, error) { return os.ReadFile(name)}
// func (osFS) Writer(name string) (int, error) { return os.WriteFile(name, data []byte )}




func TestDumb(t *testing.T){
	testFile, err := fs.Create("shit.txt")
	if err != nil{
		t.Log("Err ", err)
	}
	t.Log(testFile)


	// testFile, err := fs.Open("shit.txt")
	// if err != nil{
	// 	t.Log("Err ", err)
	// }

	// write a chunk
	if _, err := testFile.WriteString("buffered\n"); err != nil {
		t.Log("Err ", err)

	}



	buf := make([]byte, 1024)
	n, err := testFile.Read(buf)
	if err != nil && err != io.EOF {
		t.Log("Err ", err)

	}


	


	t.Log(n)
	t.Log(string(buf))


	testFile.Close()
}
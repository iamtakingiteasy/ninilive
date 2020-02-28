// Code generated for package migration by go-bindata DO NOT EDIT. (@generated)
// sources:
// 0001_init.down.sql
// 0001_init.up.sql
package migration

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __0001_initDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\x29\xca\x2f\x50\x28\x49\x4c\xca\x49\x55\xc8\x4d\x2d\x2e\x4e\x4c\x4f\x2d\xb6\xe6\x42\x12\x2c\x2d\x4e\x2d\x42\x15\x49\xce\x48\xcc\xcb\x4b\xcd\x29\xb6\x06\x04\x00\x00\xff\xff\xae\xd0\xe5\x10\x3b\x00\x00\x00")

func _0001_initDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__0001_initDownSql,
		"0001_init.down.sql",
	)
}

func _0001_initDownSql() (*asset, error) {
	bytes, err := _0001_initDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "0001_init.down.sql", size: 59, mode: os.FileMode(420), modTime: time.Unix(1582660007, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __0001_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x52\xc1\x8e\xab\x30\x0c\xbc\xe7\x2b\x7c\x2b\x48\x3d\xbc\x7b\x4f\xef\x4b\x56\xa1\xb8\xd4\xda\x24\x66\x6d\xb3\xbb\xfc\xfd\xaa\x34\x29\x85\x52\x90\x90\x92\xcc\x84\x19\x3c\x73\x16\xf4\x86\x60\xbe\x09\x08\x83\xa2\x28\x54\x0e\x00\xa6\xf5\x07\xb5\x70\x7f\x1a\xea\x14\x85\x7c\x80\x5e\x28\x7a\x19\xe1\x13\xc7\xe3\x4c\x4c\x3e\xe2\x44\x34\xfc\xb5\xa7\xf3\xc0\x1d\xa5\x72\x0e\x43\xa2\xaf\x01\x9f\xe0\xde\xab\xfe\xb0\xb4\xeb\x6b\x91\xb3\x70\xc3\x1c\x5c\x7d\x72\x6e\xe1\xf3\x7c\xf5\x29\x61\x28\x56\xf3\xf6\xee\x76\xc7\x69\xe1\x4d\x66\x5f\x0c\x15\x94\xa5\x45\x01\x4a\xf6\xaa\x1b\x51\xd5\x77\x58\x74\xf3\x76\x9e\xd2\xae\x7a\x61\x3f\xb9\x6d\xa8\xa3\x64\x20\x78\x41\xc1\x74\x46\x7d\xfc\x59\x35\xb3\xea\xe5\xf5\x86\xdb\x31\x8b\xcd\x43\x2b\xa0\x51\x8e\x01\xe0\xb6\x54\xf3\xb1\x5f\x32\xb0\x25\xdb\x67\x98\x50\xff\x56\xe0\x91\xf3\x16\xc8\x42\x39\xed\x0d\x50\x30\xb2\xad\x3b\x52\xc0\x0b\x85\xf2\xed\x37\x60\xef\xed\xba\x05\xde\xea\x32\xcf\x7e\x35\xcd\xa9\xcf\x55\x6e\x72\x3d\xe5\x49\x49\x51\xec\x16\x2f\x6f\xd5\x7d\x55\xe9\x75\x93\x37\xaa\xbb\x6c\xad\xab\xe1\xdb\x87\x01\xd5\x55\xff\x8e\x70\xf8\x9f\x38\x8d\x91\x07\x3d\x1c\xe1\x90\xdf\x8b\x0f\x8a\xf5\xe9\x2f\x00\x00\xff\xff\x91\xe1\xb4\xe2\x79\x03\x00\x00")

func _0001_initUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__0001_initUpSql,
		"0001_init.up.sql",
	)
}

func _0001_initUpSql() (*asset, error) {
	bytes, err := _0001_initUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "0001_init.up.sql", size: 889, mode: os.FileMode(420), modTime: time.Unix(1582842245, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"0001_init.down.sql": _0001_initDownSql,
	"0001_init.up.sql":   _0001_initUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"0001_init.down.sql": &bintree{_0001_initDownSql, map[string]*bintree{}},
	"0001_init.up.sql":   &bintree{_0001_initUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

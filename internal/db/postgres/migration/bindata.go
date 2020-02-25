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

var __0001_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\xc1\x6e\x83\x30\x0c\x86\xef\x79\x0a\xdf\x0a\x52\x0f\xbb\xf7\xb4\x27\x99\x42\x71\xa9\xb5\x24\x66\xb6\xd9\xc6\xdb\x4f\xa5\xc9\x28\x94\x82\x84\x94\xe4\xff\x23\xff\xb1\xbf\xb3\xa0\x37\x04\xf3\x4d\x40\x18\x14\x45\xa1\x72\x00\x30\xad\x3f\xa8\x85\xfb\xd7\x50\xa7\x28\xe4\x03\xf4\x42\xd1\xcb\x08\x9f\x38\x1e\x67\x63\xf2\x11\x27\xa3\xe1\xaf\x3d\x9c\x07\xee\x28\x95\x73\x18\x12\x7d\x0d\xf8\x20\xf7\x5e\xf5\x87\xa5\x5d\x5f\x8b\x9c\x0b\x37\xcc\xc1\xd5\x27\xe7\x16\x39\xcf\x57\x9f\x12\x86\x12\x35\x6f\xef\x69\x77\x92\x16\xdf\x14\xf6\x29\x50\x51\x59\x5a\x14\xa0\x64\xcf\x75\x23\xaa\xfa\x0e\x4b\xdd\xbc\x9d\xbb\xb4\x5b\xbd\xb8\x1f\xd2\x36\xd4\x51\x32\x10\xbc\xa0\x60\x3a\xa3\xfe\xbf\xac\x9a\x5d\xf5\xf2\x7a\xc3\xed\x98\x8b\xcd\x4d\x2b\xa2\x51\x1e\x03\xc0\x6d\xa9\xe6\x63\xbf\x74\x60\x4b\xb6\xef\x30\xa1\xfe\x65\x01\x16\xca\x03\xdd\x10\x05\x23\xdb\x1a\x83\x22\x5e\x28\x60\xc6\xe4\x85\xd8\x7b\xbb\x6e\x89\x37\x22\xe6\xf6\xae\x1a\x36\x21\x5b\x65\x58\xeb\x69\x64\x94\x14\xc5\x6e\x13\xe4\x2d\xa2\x57\xd4\xae\x61\xdd\xa0\x73\x09\xa6\xab\xe1\xdb\x87\x01\xd5\x55\x6f\x47\x38\xbc\x27\x4e\x63\xe4\x41\x0f\x47\x38\xe4\xff\xe2\x83\x62\x7d\xfa\x0b\x00\x00\xff\xff\xc3\x12\x55\x6c\x5c\x03\x00\x00")

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

	info := bindataFileInfo{name: "0001_init.up.sql", size: 860, mode: os.FileMode(420), modTime: time.Unix(1582659909, 0)}
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

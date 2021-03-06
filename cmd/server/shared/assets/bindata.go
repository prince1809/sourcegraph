// Code generated by go-bindata. DO NOT EDIT.
// sources:
// nginx.conf (43B)
// redis.conf.tmpl (300B)
// nginx/sourcegraph_backend.conf (0)

package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _nginxConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\x08\xc9\xc8\x2c\x56\x48\xce\xcf\x4b\xcb\x4c\x57\x28\x4f\x2c\x56\x48\x4f\xcd\x4b\x2d\x4a\x2c\x49\x4d\x51\x48\xaa\x54\x08\xce\x2f\x2d\x4a\x4e\x4d\x2f\x4a\x2c\xc8\xe0\x02\x04\x00\x00\xff\xff\xca\x89\xcd\x21\x2b\x00\x00\x00")

func nginxConfBytes() ([]byte, error) {
	return bindataRead(
		_nginxConf,
		"nginx.conf",
	)
}

func nginxConf() (*asset, error) {
	bytes, err := nginxConfBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "nginx.conf", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x6b, 0x85, 0x68, 0x9a, 0xd3, 0x2e, 0xa0, 0x65, 0x6b, 0xba, 0x44, 0x65, 0x3, 0x26, 0xb0, 0xcd, 0x4, 0x37, 0xa0, 0x75, 0x44, 0x63, 0xcf, 0xf6, 0xae, 0xe1, 0x2f, 0xce, 0xf6, 0x34, 0xad, 0xb1}}
	return a, nil
}

var _redisConfTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\xcd\x41\x6e\x02\x31\x0c\x05\xd0\x7d\x4e\xf1\x25\xb6\x50\xc1\xa6\x27\xe8\x45\x42\xe2\x06\xab\x89\x3d\xb2\xcd\x4c\x47\x88\xbb\x57\x29\x95\xba\xf4\xd7\xf7\x7f\x07\xe4\xde\x75\x43\x2e\x85\xdc\xf1\x69\x3a\x66\x02\x16\x8f\x2c\x85\x3c\x2d\xa6\x41\x25\xa8\x9e\x86\x56\x82\x68\x4a\x07\x74\x1e\x1c\x18\x34\xd4\x76\xdc\x3d\x37\x3a\xc2\x28\xee\x26\x20\x33\x35\x6c\x37\x12\xdc\x38\x82\xa5\xbd\xda\x69\xe4\xef\xbf\x87\x4b\xbb\xfe\x5f\xa7\x45\x3b\x97\x1d\xa2\xb4\x72\x09\x56\x79\x01\x2b\xa1\xe8\x98\x4c\xd7\x86\x50\x54\xf6\xaf\x23\x72\xad\x3c\x4b\xb9\xf7\x1d\x2e\x79\xf1\x9b\x86\x83\x56\x9a\xc3\x67\x0c\x96\x7b\x90\xa7\xca\x86\xc7\x03\x6f\x1f\x6c\x78\x3e\x53\x5e\x16\x92\xaa\xd2\x77\xec\xe4\xc9\xf3\x4a\x78\x3f\x9f\x71\xf9\xc5\x28\x7b\x60\x25\xbb\xaa\xd3\xe4\x1a\x4b\x4b\x5d\x5b\xa7\x95\x3a\xb6\x6c\x32\x83\x9f\x00\x00\x00\xff\xff\x3d\xe3\xfc\x7d\x2c\x01\x00\x00")

func redisConfTmplBytes() ([]byte, error) {
	return bindataRead(
		_redisConfTmpl,
		"redis.conf.tmpl",
	)
}

func redisConfTmpl() (*asset, error) {
	bytes, err := redisConfTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "redis.conf.tmpl", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x39, 0xdc, 0xa1, 0x4e, 0x7, 0x3e, 0xee, 0x45, 0xc6, 0x63, 0x5d, 0xf, 0x9c, 0xfc, 0xc1, 0x41, 0xd0, 0xc1, 0xc5, 0x78, 0xf2, 0x89, 0xe0, 0xe1, 0x8f, 0x38, 0xbc, 0x9b, 0x8d, 0xff, 0xcf, 0x0}}
	return a, nil
}

var _nginxSourcegraph_backendConf = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func nginxSourcegraph_backendConfBytes() ([]byte, error) {
	return bindataRead(
		_nginxSourcegraph_backendConf,
		"nginx/sourcegraph_backend.conf",
	)
}

func nginxSourcegraph_backendConf() (*asset, error) {
	bytes, err := nginxSourcegraph_backendConfBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "nginx/sourcegraph_backend.conf", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"nginx.conf": nginxConf,

	"redis.conf.tmpl": redisConfTmpl,

	"nginx/sourcegraph_backend.conf": nginxSourcegraph_backendConf,
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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"nginx": {nil, map[string]*bintree{
		"sourcegraph_backend.conf": {nginxSourcegraph_backendConf, map[string]*bintree{}},
	}},
	"nginx.conf":      {nginxConf, map[string]*bintree{}},
	"redis.conf.tmpl": {redisConfTmpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}

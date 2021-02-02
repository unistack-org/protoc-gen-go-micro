// Code generated by vfsgen; DO NOT EDIT.

// +build !dev

package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Assets statically implements the virtual filesystem provided to vfsgen.
var Assets = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2021, 2, 2, 5, 31, 18, 783804630, time.UTC),
		},
		"/e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyb19ncnBjLnBiLmdvLnRtcGw=": &vfsgen۰CompressedFileInfo{
			name:             "e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyb19ncnBjLnBiLmdvLnRtcGw=",
			modTime:          time.Date(2021, 2, 2, 6, 6, 27, 628460932, time.UTC),
			uncompressedSize: 5889,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xcc\x57\x4d\x6f\xa3\x48\x13\x3e\xd3\xbf\xa2\x06\xbd\x1a\xc1\x2b\xbb\x7d\xd8\x5b\x46\x39\x8c\xa2\x1d\xed\x4a\x9b\xcc\x6a\xb2\xf7\x88\x69\x97\x09\x0a\x34\xa4\x69\x12\x47\xa4\xff\xfb\xaa\x3f\xc0\x80\x31\xb6\x43\xb2\xbb\x39\xc4\x72\x7f\x55\x3d\xf5\xf1\xd4\xe3\xd5\x0a\xae\xf2\x35\x42\x8c\x1c\x45\x24\x71\x0d\x3f\x5f\xa0\x10\xb9\xcc\xd9\x32\x46\xbe\xcc\x12\x26\x72\xb2\x5a\x41\x99\x57\x82\xe1\x05\xd4\x35\xfd\x96\xa4\x48\x6f\xa2\x0c\x95\x22\x45\xc4\x1e\xa2\x18\xa1\xae\xe3\xfc\xcf\x87\xf8\x8f\xa8\x94\xbf\xa6\x98\x21\x97\x60\xce\xc1\x2b\x94\x45\x9a\xc8\xaf\x42\x44\x2f\xe0\x7f\xf1\xe1\x15\xd2\xa8\x94\xf0\x0a\x02\x8b\x34\x62\x08\x3e\xf5\xc1\xbf\xf3\x95\x22\x24\xc9\x8a\x5c\x48\x08\x88\xe7\xb3\x9c\x4b\xdc\x4a\x9f\x10\xcf\xf8\x70\xc7\xd2\x44\xbf\xea\xc7\x89\xbc\xaf\x7e\x52\x96\x67\xab\x8a\x27\xa5\x8c\xd8\xc3\x32\x17\xf1\xca\x9c\x5a\x3d\xfd\xb2\xb2\x07\x7d\x68\x2e\x96\x28\x9e\x50\x9c\x70\xd1\x1e\xf4\x49\x48\xc8\x53\x24\x20\x20\x00\x00\x77\xd0\x7d\x85\x7e\x2f\x64\x92\xf3\xde\x8e\xb5\xd7\xec\x84\x84\xd4\xf5\x12\xfe\x67\xc0\x5f\x5c\xba\x28\x28\x65\x57\x6f\x51\x3c\x25\x0c\x75\xf0\xcc\xa6\xfb\x6e\xa2\x09\xaf\x20\x45\x92\xdd\x56\x9b\x4d\xb2\x05\xdf\x6d\xf9\xfa\x2e\x91\x2f\x85\x8e\x71\xef\xfe\x2b\xa4\xf9\x33\x8a\x6f\x89\x28\xa5\x52\x6e\x07\x4a\x29\x2a\x26\xa1\x36\x2e\xb2\xbe\x8b\x57\xe6\xc3\xec\x70\xfd\x42\x29\x45\xc2\x63\xa2\x88\x4e\xf0\xb5\x3e\x09\x2e\xca\xa5\xac\x36\x1b\xb3\x7c\x83\xcf\x7d\xbb\x3b\x53\x4c\x60\x24\x11\x38\x3e\x43\xd9\x2c\x59\x0b\x9b\x8a\xb3\x89\x9b\x41\xc7\xfa\x62\xdc\xc9\x10\x0e\x59\xb5\xc8\x04\xca\x4a\x70\xf8\x7c\x4a\x4c\x6a\x76\x01\x6c\x61\x20\x5f\x98\xff\x4a\x43\xae\x6b\x11\xf1\x18\x77\x39\xb8\x46\x79\x9f\xaf\x9b\x4c\x09\x7c\xb4\x0b\x26\x4f\xbf\xf3\xa2\x92\x7f\xe9\x24\xf4\xeb\x99\x36\xf5\xdc\x5e\x2b\x8b\xce\xb5\xef\x95\x3c\xed\x5e\xb2\x81\x88\xaf\x21\xe0\xb9\x04\x17\x81\x5b\x29\x30\xca\x12\x1e\x87\x10\x18\x1f\x51\x74\x96\x96\x4a\xd9\x28\x07\x0c\xfe\x7f\x4a\x10\x74\x40\x5d\xd3\x06\x4c\x6e\xc1\xf5\x17\xbd\xb2\x9f\x0b\x10\xf8\x68\x5e\x6a\x81\x2b\xb5\x80\xbc\x90\x25\x50\x4a\xfb\x19\x8a\xd2\xd4\x56\x7b\x08\xc1\x30\x4d\x77\xad\x19\xb7\xbc\x00\x14\x22\x17\xa1\x49\x9c\xc6\x8a\x69\x89\x1a\xf0\x10\xe7\x07\x60\xfa\x30\xf7\xff\x33\xe1\x37\x57\x9a\x9a\xd3\x57\x86\xde\xf2\xb5\x71\xd6\x7e\x4b\x36\x90\x8b\xb1\x72\x0a\xf6\x8a\xce\xdc\x01\xdd\xa4\x18\x65\xe6\x59\x5d\xd1\x8c\x32\x6a\xcf\x68\x14\x0b\xf3\xfd\x06\x9f\x7f\xe0\x63\x85\xa5\x0c\x18\xd5\xed\xb5\x00\x7f\x18\x56\xda\xa2\xf7\x17\xa6\x69\x3b\x30\x6b\x15\x5a\xa4\x94\xd2\xd0\x18\x4d\x36\xc6\xde\xa7\x4b\xe0\x49\xea\x1a\xbe\x6d\x79\x9e\xa4\xc6\x1d\xb3\xda\x01\x36\xd6\x3a\xa0\x14\xfc\xb3\x7f\x5d\xff\x2f\x2e\x5d\xf8\xe8\x2d\xf2\x75\x20\xf0\x31\xfc\x72\x2e\x30\xe4\x6b\x97\x89\xb3\x18\xaf\x89\x76\x6d\x1d\x50\x0b\xfd\x7a\xa7\x80\x9b\x37\xcb\x42\x7b\xf9\xb9\x5f\x44\xb5\xdd\xec\xa4\x5c\x97\xdc\x8c\x84\x6b\xe4\x0b\x6d\xec\xa4\x3c\x1f\x0e\x48\xbb\x63\x9e\xea\x20\x32\x41\x32\x9c\x7e\x5e\x85\x9f\x3c\x56\x5b\x34\xfd\xf9\x6a\xa3\xdb\x9f\x5f\xd6\x00\x51\xd6\x19\xc3\xea\x23\x8c\x6e\xaa\x75\xe8\x64\x4b\x2a\xdb\xd3\x48\xa5\xf5\x2a\x84\x1f\xc8\x9e\xbe\xf2\xf5\x55\x9a\x97\x18\x4c\xf2\x82\x97\x8d\xa7\xdc\x73\xf9\xde\x52\xfd\xd6\x75\x19\x07\x59\x48\x3c\x97\xa2\xcb\x26\x45\x9e\x39\xa6\x4f\x39\x53\xc4\x53\xa4\x3d\xf6\x69\x77\x6c\x98\x43\x4f\x11\x8f\x34\xab\x99\xcd\x9e\x9d\x7d\x36\x7b\x6f\x46\xde\x40\x36\x00\xfb\xea\x60\x4b\x5d\xff\x35\xbe\xce\x31\x63\x89\x3a\x08\x87\xd4\x7d\xc8\x64\x73\x7e\x8e\x51\xcd\x1b\x26\x13\x90\x70\x89\x62\x13\x31\xac\xd5\x34\x54\x43\x35\xd9\x2c\xab\x6d\xfe\x4f\xb7\xaa\xaf\x38\xab\x75\x7d\x64\xb8\xbf\x2d\x0c\x41\x36\x1c\x8f\xa7\xc6\xa1\x53\x63\xce\xb7\x41\xdf\xcd\xf2\xcd\x20\x3f\x32\x88\xb3\x43\x24\xbb\x1b\x14\xc3\x48\x8e\x0c\x8a\x7d\x56\xd4\x4f\x1c\xe8\x28\xb0\x70\x1b\xe4\x5d\x29\xb0\x13\xfb\xee\x97\x91\x13\xfb\xc7\xc9\xf0\xb7\x88\xaf\x53\x73\xa1\xe5\xc0\x21\xef\xbb\x23\xc6\x91\x51\x7d\x0d\xff\x9a\xc0\x3e\x63\x2e\xd8\x62\xb8\x3f\x52\x0c\x0e\xeb\x51\x69\xd7\x9b\x12\xee\x77\xa4\x35\xb7\x2b\xe1\x09\x11\x33\x43\xc3\x10\xaf\xc3\xf6\x5d\xc5\x35\x4f\xe3\x78\xde\x9e\xc2\x19\x2f\x5b\xcf\x6b\x07\x81\x99\x01\x5e\x6f\x8a\xdf\xd3\x03\xd5\x43\x7b\x01\x5d\xe8\xe2\x9e\x96\x3d\x3b\xd1\x6e\xbc\x69\x54\x4f\x48\x86\x7a\xe7\x4c\xc3\x73\xac\x3a\x4d\xd2\xf1\xe0\x7d\xab\x6a\xf4\x07\x83\x16\x74\x03\x26\xea\xd2\xe4\xb9\xf8\x05\x3e\x9a\x37\xfb\x34\xfa\xfe\x2a\x6b\x10\xc9\x2e\xc1\x1c\x6e\x9e\x8f\x95\x58\x03\x97\xec\x10\x6a\x15\x56\x36\x15\xe5\x96\xd2\xdb\xf1\x1d\xee\xb8\xfe\xb2\x4b\xe9\x8d\x96\x1a\xc8\x94\x1e\xb3\xeb\xc6\x79\x83\x50\xda\xf3\x7f\x5f\x27\xcd\x56\x49\xfb\x36\x26\x45\xd2\x3b\x48\xa4\xd1\xac\x4c\x6a\x95\xd9\xfa\x68\xcf\xe4\x51\x79\x34\x29\x8e\x96\x63\x0a\x64\x5e\x55\x4e\x57\xe3\x94\x30\x6a\x65\xd1\x88\x62\x7b\xbb\x4f\x3d\x49\xd4\x65\xa7\x51\x49\xd4\x1d\x4b\x1f\x25\x89\x5a\x9c\x03\x3d\xa4\x14\xf9\x3b\x00\x00\xff\xff\x55\x9e\x06\x89\x01\x17\x00\x00"),
		},
		"/e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyb19odHRwLnBiLmdvLnRtcGw=": &vfsgen۰CompressedFileInfo{
			name:             "e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyb19odHRwLnBiLmdvLnRtcGw=",
			modTime:          time.Date(2021, 2, 2, 6, 28, 27, 560473834, time.UTC),
			uncompressedSize: 7268,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xcc\x59\xdb\x6f\x9b\xc8\x1a\x7f\x86\xbf\xe2\x2b\x27\xa7\x82\x23\x3c\x7e\xe8\x5b\xaa\x3c\xf4\x44\xed\x39\xbb\xda\xb4\x55\x52\xed\x4b\x55\x45\x53\xfc\x99\xb0\x81\x81\xcc\x8c\x9d\x44\x64\xfe\xf7\xd5\x5c\xc0\x80\xb1\x63\x97\x74\x77\xfb\xd0\x84\xe1\xbb\x5f\x7f\x4c\xe6\x73\x38\x2f\x17\x08\x29\x32\xe4\x54\xe2\x02\xbe\x3f\x42\xc5\x4b\x59\x26\xb3\x14\xd9\xac\xc8\x12\x5e\xfa\xf3\x39\x88\x72\xc5\x13\x3c\x85\xba\x26\x1f\xb2\x1c\xc9\x47\x5a\xa0\x52\x7e\x45\x93\x5b\x9a\x22\xd4\x75\x5a\x7e\xbe\x4d\x7f\xa3\x42\xbe\xcf\xb1\x40\x26\xc1\xd0\xc1\x13\x88\x2a\xcf\xe4\x3b\xce\xe9\x23\x04\x6f\x03\x78\x82\x9c\x0a\x09\x4f\xc0\xb1\xca\x69\x82\x10\x90\x00\x82\xeb\x40\x29\xdf\xcf\x8a\xaa\xe4\x12\x42\xdf\x0b\x92\x92\x49\x7c\x90\x81\x0f\x10\x2c\x0b\x19\xf8\xbe\x67\x6c\xb9\x4e\xf2\x4c\x4b\x0f\xd2\x4c\xde\xac\xbe\x93\xa4\x2c\xe6\x2b\x96\x09\x49\x93\xdb\x59\xc9\xd3\xb9\xa1\x9a\xaf\xdf\xcc\x2d\x61\x00\x0d\xa3\x40\xbe\x46\x7e\x00\xa3\x25\x0c\xfa\x0a\xaf\x6f\xa4\xac\x9e\x61\x9e\x59\xd2\x99\x26\x9d\xaf\xdf\x04\x7e\xe4\xfb\x6b\xca\x21\xf4\x01\x00\xae\xa1\x6b\x07\xf9\x54\xc9\xac\x64\xbd\x37\x96\xbd\x79\x13\xf9\x7e\x5d\xcf\xe0\xc4\x84\xf1\xf4\xcc\xc5\x53\x29\x7b\x7a\x85\x7c\x9d\x25\xa8\xd3\x60\x5e\xba\x67\x93\x17\x78\x02\xc9\xb3\xe2\x6a\xb5\x5c\x66\x0f\x10\xb8\x57\x81\xe6\xf5\xe5\x63\xa5\xb3\xd5\xe3\x7f\x82\xbc\xbc\x47\xfe\x21\xe3\x42\x2a\xe5\xde\x80\x90\x7c\x95\x48\xa8\x8d\x89\x49\xdf\xc4\x73\xf3\xc3\xbc\x61\x5a\x82\x90\x3c\x63\xa9\xaf\x7c\x5d\x2a\x17\x9a\x12\x5c\x9e\x84\x5c\x2d\x97\xe6\xf8\x23\xde\xf7\xf5\x6e\x54\x25\x1c\xa9\x44\x60\x78\x0f\xa2\x39\xb2\x1a\x96\x2b\x96\xec\xe1\x0c\x3b\xda\xe3\x71\x23\x23\xd8\xa5\xd5\x7a\xc6\x51\xae\x38\x83\xd7\x87\xc4\xa4\x4e\x4e\x21\x89\x8d\xcb\xa7\xe6\x7f\xa5\x5d\xae\x6b\x4e\x59\x8a\x9b\x1c\x5c\xa0\xbc\x29\x17\x4d\xa6\x38\xde\xd9\x03\x93\xa7\x5f\x58\xb5\x92\x5f\x74\x12\xfa\x9d\x41\x9a\xce\x68\xd9\x44\xd5\x61\xfb\xb4\x92\x87\xf0\x41\xb6\x04\xca\x16\x10\xb2\x52\x82\x0b\xc0\x95\xe4\x48\x8b\x8c\xa5\x11\x84\xc6\x44\xe4\x9d\x23\xa5\x6c\x8c\xc3\x04\xfe\x73\x48\x08\x74\x38\x5d\xf3\x87\x89\x7c\x00\xd7\xa7\xe4\xdc\xfe\x8c\x81\xe3\x9d\x91\xd4\xba\xad\x54\x0c\x65\x25\x05\x10\x42\xfa\xf9\xa1\x79\x6e\x6b\x3d\x82\x70\x98\xa4\xeb\x56\x8d\x3b\x8e\x01\x39\x2f\x79\x64\xd2\x56\xd7\x80\xb9\x40\xed\xee\xd0\xcb\x97\xf7\xe8\x67\x19\xff\x4f\x09\xbd\x61\x69\xaa\x4d\xb3\xf4\x6c\x9d\x01\xb2\x85\xb6\x55\x37\x8b\x7e\xcc\x96\xa0\xab\x2b\xd4\xea\x69\xc6\x04\x84\x7f\x88\x92\x41\x58\x56\xc8\x68\x95\x59\xb1\x40\xa2\x08\x02\xb6\xca\xf3\x20\x1a\x30\x6f\x11\x92\x4b\x14\x55\xc9\x04\x8a\x86\x12\x39\x2f\x68\xa5\x0b\xbf\xa0\xb7\x18\x16\xb4\xfa\x6a\x7b\xfc\x5b\xc6\x24\xf2\x25\x4d\xb0\x56\xb1\x0e\x64\x8e\x6c\x9f\x40\xa5\xa2\x56\xb5\x6d\xd2\x93\xdb\x18\x4e\xd6\x5a\xf4\xc1\x76\x7c\x0d\xea\xfa\xe4\x56\xa9\xe0\x1b\x9c\xe9\x39\x31\x83\x30\x45\x79\x81\x42\xd0\x14\x4d\x4b\xda\x19\x7d\xb2\x26\x57\xc9\x0d\x16\x94\xfc\x2a\x4a\xe6\x7e\xbd\xc4\x65\x64\x87\xb2\x52\xf5\x26\x0e\x83\x98\xf6\x1e\xbb\x4f\xcc\xa4\xef\xf4\x0c\x68\x55\x21\x5b\x84\xfa\x31\x36\x6f\x00\xb6\xd6\x93\x1b\x3c\x61\x50\xd7\xfa\xf1\x77\xe4\xdf\x81\x28\x15\x44\x31\xec\x64\xf9\x4c\xe5\x4d\xc3\xa0\x7f\x77\x0c\x8e\x7e\x3c\xdf\x9a\xf6\xbf\xe5\xe2\x11\x48\x04\xc1\xff\xde\x7f\xd9\xa4\x78\x4c\x83\xa6\x6c\x34\x58\xae\xa1\x86\x8e\xbf\x13\x4b\xec\x98\x22\x1b\xb3\xf5\xbd\xae\xfc\x0b\x5a\x85\x36\xf3\xfb\xad\xec\x1c\x44\x7e\xab\xb9\xe4\x63\x63\x36\xdc\x1a\xc6\x8e\x53\x98\x03\xd3\x73\x3a\xcd\x09\x49\x88\xa5\xd1\x2d\x1e\x9b\xe7\x8f\x78\x7f\x89\x77\x2b\x14\x32\x4c\x88\xde\x3a\x31\x04\xc3\x89\x43\xda\xd1\x10\xc4\x66\x97\x75\x66\x40\xad\xa2\xd8\x16\x12\x21\xc4\xf6\x43\xb6\x34\x0a\x5f\x9d\x01\xcb\x72\xb7\x08\xdb\x55\xc8\xb2\xdc\xd8\x63\x4e\x95\xdf\x4b\xc9\xc8\xb4\x85\xbf\xf6\x5f\xd7\xfe\xd3\x33\x17\x3f\x72\xa5\x9b\x83\xe3\x5d\xf4\xf6\x58\xc7\x90\x2d\x5c\x2a\x8e\x42\x02\x4d\xb8\x6b\x6b\x80\x8a\xb5\xf4\x46\x64\x2e\xb0\x91\x29\xcc\x18\x7b\xdd\x1f\xb1\x75\x3b\x5e\x9a\x9c\xeb\x81\x3c\x21\xe3\xda\xf3\x58\x2b\x3b\x2c\xd1\xbb\x23\xd2\xbe\xb1\xb2\x36\x2e\x99\x28\x19\xb0\x73\x5c\x8d\x1f\x8c\x37\x5b\x77\xfa\xc0\xd3\x86\xb7\x0f\xec\xac\x02\x5f\x59\x63\x0c\xde\x19\xc1\x3a\xa6\x5c\x87\x46\xb6\x3b\xf7\xe1\xb0\x9d\xdb\x5a\x15\xc1\x25\x26\xeb\x77\x6c\x71\x9e\x97\x02\xc3\xbd\x6b\xd3\x2b\xc6\x73\xee\xb9\x84\x3f\x10\x2d\xeb\x42\xa4\x61\x11\xf9\x9e\x4b\xd1\x59\x93\x22\xcf\x90\x69\x2a\xa7\xca\xf7\x94\xdf\x92\xbd\xda\x90\x0d\x73\xe8\x29\xdf\xf3\x9b\xd3\xc2\x66\xcf\x82\x49\x9b\xbd\x1f\xf6\xbc\x71\xd9\x38\xd8\x87\xcd\x0f\xc4\x35\x60\x63\xeb\x14\x35\x16\xc7\x84\xd1\x10\xd9\xec\x52\xd9\xd0\x4f\x51\xaa\x07\x87\xc9\x04\x74\xa0\xc5\x7e\x57\xcd\xac\x29\x26\x69\x6d\xf3\x7f\xb8\x56\xcd\xe2\xb4\x5a\x9c\xbf\x35\x8a\x67\x13\x8a\xdb\xfa\x34\x44\x8f\x87\xc6\xa1\x53\x63\xce\xb6\x41\xdf\x4d\xb2\xcd\x78\xfe\x0c\x4e\x2d\x76\x4d\xd9\xcd\xa6\x18\x46\x72\x64\x53\x6c\x4f\x45\x2d\x62\x47\x47\x81\x75\xb7\xf1\xbc\x39\xd5\x9e\x6e\xbe\x82\xdd\xa5\x83\xfb\x0a\x7e\x7e\x18\xfe\x9f\xb2\x45\x6e\x18\xda\x19\x38\x1c\xfc\x8e\xc4\x18\x32\xfa\xe1\x09\x7f\xc7\x97\xe7\xb1\xd8\xc7\x16\xc3\xcd\x33\xc5\xe0\x7c\x7d\xf6\xcb\xa7\xb7\x25\xdc\x05\x8b\x55\xb7\x29\xe1\x3d\x28\x66\x02\x88\xf1\xbd\xce\xb4\xef\x62\xae\x69\x20\xc7\xf3\xb6\x20\xce\x78\xd9\x7a\x5e\xbb\x08\xcc\x0e\xf0\x7a\x5b\xfc\x86\xec\xa8\x1e\xd2\x0b\x68\xac\x8b\x7b\x3f\xee\xd9\x7c\xd1\x1a\x6b\x1a\xd8\x13\xf9\x43\xc0\x73\xa4\xe2\x29\x5a\x1d\x26\xe9\x58\xf0\xb2\x55\x35\xfa\x3d\xad\x11\xdd\x60\x12\x75\xc7\xe4\xb1\xfe\x73\xbc\x33\x32\xb7\xc6\xe8\x4b\xc3\xac\x41\x28\xbb\x13\x66\x77\xf7\xfc\x5c\x8c\x35\x30\xc9\x6e\xa1\x16\x62\x15\xfb\xc2\xdc\xce\xf4\x76\x7f\x47\x9b\x61\x7f\xd6\x9d\xe9\x0d\x98\x1a\xe0\x94\xde\x68\xd7\x9d\xf3\x03\x48\x69\xcb\xfe\x6d\xa0\x34\x19\x26\x6d\xeb\xd8\x8b\x92\x5e\x00\x23\x8d\x66\x65\x2f\x58\x99\x0c\x90\xb6\x54\x3e\x8b\x8f\x76\xa3\x23\x53\xaf\xc3\x9a\x9c\x56\x92\xfb\x4b\x71\x1f\x2c\xea\x76\xf3\xf6\xd2\x99\x12\x9d\x16\x10\x75\x67\xd3\x28\x20\xea\x2e\xa5\x9f\x05\x88\x5a\x3f\x1b\x2c\xd4\x1e\xcc\xe0\x04\x39\x2f\x44\xaa\x35\xe6\x99\x90\x16\x9d\xec\xbb\x31\xff\xa1\xdb\x9f\xc3\xee\x7d\xea\xfa\xd8\x6b\x40\xe3\x41\x21\x52\x33\x60\x35\xf9\x71\x57\x7f\xad\x08\x17\x84\xe6\x22\xaf\x3d\x68\x64\x37\x77\x84\x5d\x74\xf9\xec\xaf\x9b\xbf\x3d\x74\xfc\x69\x24\x3f\xc1\x8a\x65\x77\xc6\x82\xf9\x1c\xcc\xcd\x16\x14\x16\xd7\xc9\x12\x04\x95\x99\x58\x3e\xba\x3a\x6e\xbb\xcc\x95\x24\xea\x92\x9c\x69\x89\xba\xd6\x0d\x6f\x18\xb9\xbf\xb3\xf4\x4b\x63\x59\x48\x72\x55\xf1\x8c\xc9\x65\x18\xfc\xfb\x5f\xeb\x20\x06\x8c\xfa\x20\xf9\xcf\x00\x00\x00\xff\xff\x5e\x15\x5d\xbc\x64\x1c\x00\x00"),
		},
		"/e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyby5wYi5nby50bXBs": &vfsgen۰CompressedFileInfo{
			name:             "e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyby5wYi5nby50bXBs",
			modTime:          time.Date(2021, 2, 2, 6, 16, 50, 658000044, time.UTC),
			uncompressedSize: 5590,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xd4\x57\x51\x6f\xa3\x38\x10\x7e\x2e\xbf\x62\x84\xaa\x0a\x56\x29\xd9\xd3\xbd\xe5\x94\x87\x6e\xb5\xab\x3b\xe9\xda\x5d\x6d\x4f\x77\x0f\xab\x55\xe5\xc2\x84\xf8\x4a\x0c\x67\x9b\xb4\x15\xe5\xbf\x9f\x8c\x0d\x18\x97\x90\xa4\xaa\x56\xda\xbc\x04\xdb\xf3\x8d\x87\x99\xcf\xc3\xe7\xf9\x1c\x2e\xf3\x04\x21\x45\x86\x9c\x48\x4c\xe0\xee\x09\x0a\x9e\xcb\x3c\x3e\x4f\x91\x9d\x6f\x68\xcc\x73\x6f\x3e\x07\x91\x97\x3c\xc6\x05\x54\x55\xf4\x89\x66\x18\x5d\x93\x0d\xd6\xb5\x57\x90\xf8\x9e\xa4\x08\x55\x95\xe6\x5f\xee\xd3\x3f\x89\x90\x1f\x33\xdc\x20\x93\xd0\xd8\xc1\x33\x88\x22\xa3\xf2\x82\x73\xf2\x04\xfe\x6f\x3e\x3c\x43\x46\x84\x84\x67\xe0\x58\x64\x24\x46\xf0\x23\x1f\xfc\x5b\xbf\xae\x3d\x8f\x6e\x8a\x9c\x4b\x08\xbc\x13\x3f\xce\x99\xc4\x47\xe9\x7b\xde\x49\x13\xc3\x2d\x29\x28\xf8\x29\x95\xeb\xf2\x2e\x8a\xf3\xcd\xbc\x64\x54\x48\x12\xdf\x9f\xe7\x3c\x9d\x37\x26\xf3\xed\xaf\x73\x52\x50\x1f\x00\xa0\x45\xc5\x19\x55\xb1\xec\x07\x6a\x43\xbf\x03\x0a\xe4\x5b\xe4\x07\x00\xb5\xa1\xef\x85\x9e\x57\x55\xe7\x70\x7a\x83\x7c\x4b\x63\x54\xf9\x81\xc5\x12\x22\x33\x6e\x12\x06\xcf\x20\x39\xdd\xdc\x94\xab\x15\x7d\x04\xdf\x2c\xf9\x50\xd7\x1a\x8b\x05\x2b\x37\x0a\xf5\xbe\x9d\xe2\x84\xa5\xd8\x3b\xb9\x42\xb9\xce\x13\xb3\x46\x57\xc0\x72\x09\x81\x4a\x15\xa1\x4c\x40\xf0\xaf\xc8\x19\x04\x6b\x29\x8b\xcf\x85\xa4\x39\x83\x28\x0c\xc1\x67\x65\x96\xf9\x21\x58\x28\xd4\x46\x7f\x23\xbf\x83\x28\x04\xdf\x8d\x60\x09\x24\x49\xda\xc1\x2f\xc3\x58\x14\xf0\x0b\x91\x6b\x71\x91\x24\x54\x6d\x42\xb2\x0f\x94\x25\x94\xa5\x02\xa2\x83\xfc\x20\x4b\x0e\x7a\x54\xb4\xbb\xc6\x87\xaa\xb2\x73\x5a\xd7\x1f\x59\x52\xe4\x94\x49\xa1\x78\xba\xa5\x09\x0a\x50\xdc\xc0\x6e\x7a\x83\x32\x21\x92\xc0\x2a\xe7\xe0\x82\x41\xe8\x91\xb7\x2a\x59\x3c\xe9\x3d\x08\xe1\xdb\xf7\x77\x1d\xf7\xa2\x76\x01\x2a\x4f\x11\xac\xdf\x6e\xb1\x84\x0d\xb9\xc7\x60\xd4\x7a\x06\xef\x67\x50\x55\x6d\x0e\xea\x3a\x6c\xd0\x7d\x21\xcc\x42\x53\x70\xb5\xb2\x25\xbc\xf3\x0d\x23\x0e\x3b\xb8\x49\x52\x3b\xdc\x45\x14\x7b\xb7\xe3\xc8\xd2\x22\x4f\x2d\x93\xc5\xd2\x41\x80\xb3\x05\xda\xe6\x26\x0a\xc3\x2f\x3b\x6b\xb0\x84\xb3\x97\xef\xa6\x33\xab\x7e\xaa\x16\x0b\xf0\xdd\xe2\x44\x55\x65\x1a\x8f\x3f\xeb\x6c\x15\x19\x17\xf0\xed\xbb\x90\x9c\xb2\xb4\x52\x20\x2b\x04\xb5\x5a\xd7\x7e\xdd\xdb\xeb\xa0\x76\x23\xda\xd4\xd9\x98\x0f\x79\xf2\xa4\xe3\xb1\x0c\xd5\xe4\x20\x12\x93\x84\x9c\x43\x10\x5d\x36\xfd\xe4\x46\x72\x24\x1b\xca\xd2\x10\x82\xa6\x36\xc8\xfb\x29\x93\x13\xf5\xd3\x93\x0b\x90\xbc\xc4\xa1\x3f\x64\x89\x65\xf7\x3b\x61\x49\x86\x7c\x01\x3e\x2f\x62\xb3\x73\xed\xf0\x71\x09\xa4\x28\x90\x25\x41\x37\x35\xeb\x56\x43\x87\x2f\xa7\x94\x25\xf8\x38\x83\x53\x34\x4d\x7b\xb1\x1c\x14\xb0\x3f\xe3\x3f\xb8\x82\x26\x9e\x43\xcb\xd7\x9a\x4f\xd7\xae\xb5\xfa\xa9\x0b\xd7\xfb\x3d\x7c\xc4\x51\x96\x9c\xf5\x5b\x79\xba\xb9\xba\xc5\x31\x03\xa0\x4c\x22\x5f\x91\x18\x3d\xf9\x54\xe0\x7e\xb3\xa6\x25\x4e\x37\xa1\xa6\x91\x70\xfc\xcf\x74\x04\xf5\x61\xfc\x83\x15\xa5\xfc\x4b\x6d\x30\xd4\x08\x51\xa7\x11\x2c\xa4\x28\x2c\xe4\xe7\x52\x1e\x0c\x6d\xcb\xea\xd6\x70\xa4\xd2\x36\x86\xb0\x64\x14\xa4\x7a\xe8\x38\x50\x41\x0d\xb5\x03\x23\x5e\xa2\x4b\xfd\x3f\x83\x77\x55\xd5\xbf\x7d\x5d\xcf\x20\x8a\x22\x5b\xa0\x44\x97\x24\xcb\xf4\x99\x0b\x21\x70\x13\x7e\xdb\x79\x36\xd3\x33\x40\xce\x73\x1e\x9a\x78\x31\x13\xb8\x3f\x86\xb7\xdd\x92\xf5\x75\x3d\x6c\xfb\xe3\x52\xd0\x58\xb7\x45\x57\xd6\xbb\x77\x6f\x9e\x6b\xa5\xbe\xa6\xc4\xd2\xb1\xe4\xeb\x60\x47\x32\xaf\xc1\x1d\x47\xba\xd1\x43\xf6\xa2\x00\xce\x71\x03\x30\x99\x0d\x42\x70\x72\xdd\xac\xde\x20\x4b\xae\x44\x1a\x74\xa8\xaa\x0e\x75\x12\x9b\xe5\xaf\x18\x6f\x27\x96\x07\xa7\xe0\x65\x47\x6c\x4e\xc1\x8e\xae\xa8\x3c\x5f\xb0\xe4\x32\xcb\x05\x06\x93\x95\x74\xbb\x54\x0b\x79\x11\x86\x1b\x81\xb1\x57\xaf\x18\x38\xb4\x72\xd1\xc3\x8e\xa8\x7c\x39\x61\x5b\x51\x1f\x1e\x6d\x27\x53\x6b\xfb\x69\x3e\x87\x2b\x45\x68\x30\x77\x07\x21\xcb\xd5\x6a\xb4\xd1\x9a\xaf\x41\x6b\xb8\xd6\xc3\x71\x22\xb4\xb6\x76\xf9\x4f\xf6\x5d\x0d\x7e\x28\xdb\x8f\x6f\xb2\x03\x49\xba\xa3\xba\xc7\xb4\x92\x89\xc3\xd3\x78\x7d\xc1\x8a\x83\x3a\xd6\xd1\x5e\x7b\x3e\xbc\xae\x23\x3a\xe4\x6b\xfd\x8f\x71\x4d\x7f\xc0\xbf\x62\x4a\x85\x44\xbe\x8b\x33\xdc\xac\x0b\x97\x69\xcd\xed\x67\x0f\x3a\x10\x60\x5f\x85\x4d\x71\x67\x20\xd6\xbb\x38\x3a\x83\xbc\x90\xa2\x6f\xec\x06\x68\x56\xdb\xde\xde\xbc\x94\xe9\x61\x23\x8c\x57\xcc\xca\x1f\x90\x7f\xa2\x5c\xf1\x6b\xc8\xfb\x9f\x9f\xf8\x7b\x38\x31\xcc\xf8\x80\x66\x0e\xb1\xde\x88\x5c\x3a\xa5\x9a\x54\xf6\x73\x2f\x50\x47\xbb\x12\x08\xc9\xcb\xb8\xbd\x0c\xeb\x58\x26\xaa\x68\xb9\x5b\xab\x0c\x9f\x4d\x5a\x1b\xc2\x54\x62\xad\x21\xea\x22\x7f\xdb\x6b\x61\xe5\x40\x73\x60\xcf\xfd\xbd\x0f\xae\xe1\x65\xa7\xaf\xd5\xa8\x4d\xb5\xba\xbd\xfc\x43\xe5\xba\xc5\x75\xe2\x3b\x0c\xad\xa0\x8d\x80\x16\x86\xcb\x81\x88\xae\xf1\xa1\x3d\x27\x67\x6e\x0c\xd5\xba\xd6\x47\x21\x8a\xa2\x30\xd4\xda\xe4\xad\x69\xab\xdf\xeb\xb5\xe4\x6d\xc1\x6f\x2b\x54\x1a\x5b\x47\xa7\x4c\xa9\x94\x49\x8d\x32\xa9\x50\x5e\xad\x4f\xd4\x96\x9d\x3a\xd9\x71\x24\x86\xca\xc1\x55\x25\x53\x3a\xa2\x57\x24\xfb\xbc\x4e\x29\x9b\x81\x16\xb1\x8f\xf0\x88\x06\x1e\x57\x22\xff\x07\x00\x00\xff\xff\xf7\xa5\xf6\xe8\xd6\x15\x00\x00"),
		},
		"/micro_chi.pb.go.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "micro_chi.pb.go.tmpl",
			modTime:          time.Date(2021, 2, 2, 4, 57, 11, 52217487, time.UTC),
			uncompressedSize: 1576,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xa4\x54\x5d\x6f\xe3\x36\x10\x7c\x36\x7f\xc5\x96\x85\x01\x29\x27\x53\x3d\xb4\x4f\xee\xdd\x43\x7b\xc8\x01\x41\x9b\x34\x08\x8a\x04\x68\x10\x04\x0c\xb5\x92\x08\x4b\xa4\x4a\xae\x5c\x1b\xb6\xfe\x7b\x41\x49\xfe\x48\x9c\x04\x45\xcf\x80\x41\x93\xdc\xd9\x9d\xd9\x1d\x33\x4d\xe1\x8b\xcd\x10\x0a\x34\xe8\x24\x61\x06\x4f\x6b\x68\x9c\x25\xab\x66\x05\x9a\x59\xad\x95\xb3\xac\x91\x6a\x21\x0b\x84\xcd\xa6\xb0\xd7\x8b\xe2\x77\xe9\xe9\xbc\xc2\x1a\x0d\x81\xf8\xaa\x2b\x84\x2d\xf8\xa6\xd2\xf4\x8b\x73\x72\x0d\xfc\x67\x0e\x5b\xa8\xa4\x27\xd8\x82\xc3\xa6\x92\x0a\x81\x0b\x0e\xfc\x91\x77\x1d\x63\xba\x6e\xac\x23\x88\x18\x00\x57\xd6\x10\xae\x88\x87\xdf\x79\x3d\xac\x06\x29\x2d\x89\x9a\x7e\xe3\x30\xaf\x50\x0d\x17\x9e\x9c\x36\x85\xe7\x8c\x4d\x78\xa1\xa9\x6c\x9f\x84\xb2\x75\x5a\xd8\x99\x2a\x75\x1a\xbe\xcb\x9f\x38\x9b\xd4\x3a\xcb\x2a\xfc\x47\x3a\x84\x37\xc3\xd2\x43\x50\x8f\x50\xce\x3e\xca\x46\x3f\x03\xb4\x46\x7b\x92\x6a\x31\xb3\xae\x48\xfb\x90\x74\xf9\x63\x2a\x1b\xcd\x59\xcc\x18\xad\x1b\x04\x67\x5b\xc2\xdf\x70\x0d\x9e\x5c\xab\x68\xd3\x31\x96\xb7\x46\xc1\x4d\x38\xbf\x92\x35\x46\x8a\x56\x30\x8a\x14\x5f\x86\x35\x86\x68\x50\x92\xc0\x93\xb5\x55\x0c\x1b\x06\xe3\x67\x29\xab\x16\x13\xb0\x0b\x98\x7f\x06\x45\x2b\x71\x1b\x0e\xa2\x5d\x99\x4d\x17\x8b\x11\x1b\xef\x31\x0e\xa9\x75\xe6\x00\x65\x7b\x12\x58\x68\x4f\xe8\x22\x07\x67\xaa\xd4\xe2\xb2\x5d\x25\x50\x82\x36\x84\x2e\x97\x0a\x37\x5d\x02\xd8\x78\xb8\x7f\x38\xdb\x37\x40\x9c\x9b\xac\xb1\xda\x50\x0c\xe8\x9c\x75\xc7\xdc\x02\xa7\x71\x1c\x03\xaf\x3f\xf2\xa8\x8c\xd9\x3e\x40\xe7\xb0\x14\x57\x6d\x7d\x89\x54\xda\x2c\x8a\xe1\x13\x7c\x3c\xc2\xbf\xe0\x9b\xd7\x24\xce\x43\x89\x3c\xe2\xa5\x34\x59\x85\x0e\x4a\xe9\xc1\x58\xa8\xfb\x04\x7e\x0e\xd3\x3f\x79\x02\xe5\x41\x6a\x77\x28\x96\x5b\x07\x8f\x81\x7f\xcf\x4a\x9a\x02\x7b\x2d\xa7\xe5\x74\xb6\x0a\x21\xa3\x77\xc4\x85\xc9\x70\x15\x61\x23\xc2\x74\x92\x60\xcb\xf8\x14\x92\xf7\xa8\x40\x7f\xbb\x85\x0a\xcd\x2e\x3e\x86\x4f\x9f\xfb\xab\xd3\x32\xef\xa8\xd3\x66\x29\x2b\x9d\xc1\x71\x7b\xc1\xc8\x1a\xe7\x30\xf5\x3c\x68\x18\x92\x9f\xe4\xec\x4e\x4e\x02\x2a\xa8\x19\x21\xf7\x3a\x5b\x7d\xf8\x38\x7f\x38\x89\xab\x43\xd0\x52\x0c\x93\xf8\x75\xdd\x3b\xd1\xbc\x5a\x43\xe7\xf0\x5d\x2d\x2e\xfc\x6d\xa0\x18\xc5\x41\x70\xd8\xfe\x85\xce\x46\xf1\xff\xd3\x39\x4e\x33\x19\xe7\x08\xd3\x30\x54\x82\xdc\xb6\x26\xe3\x09\x98\x37\xb4\x9e\x1a\xa5\xdc\xfd\x13\x6a\x71\xb1\x33\x6d\x14\x8b\x28\xb8\x3b\x0a\x2f\x84\xb8\x41\xdf\x58\xe3\xf1\xce\x69\x0a\x15\xcf\xc6\xd3\xbf\x5b\xf4\x14\xbf\xae\xd6\x2e\xbe\x49\xd6\x1c\xa6\xdf\x7f\x58\xf2\xe4\x39\xa9\xff\x32\xbc\xd1\xb2\x63\x57\x8e\x6c\x3b\xce\xe9\x3d\x5a\xe2\x4e\x53\x19\x1d\x1e\xad\x7e\xff\xf2\x71\x38\x58\x29\x1e\x53\x7e\x0d\x9d\x1a\x0a\xf6\x97\xd7\x92\xca\xfb\x1f\x1e\x12\x70\xe5\x7b\x8c\x3b\x36\x99\x4c\xc6\x56\x18\x5d\xb1\x8e\xfd\x1b\x00\x00\xff\xff\xcb\x9f\xcd\xc3\x28\x06\x00\x00"),
		},
		"/micro_gorilla.pb.go.tmpl": &vfsgen۰CompressedFileInfo{
			name:             "micro_gorilla.pb.go.tmpl",
			modTime:          time.Date(2021, 2, 2, 4, 57, 24, 265846647, time.UTC),
			uncompressedSize: 1261,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xa4\x53\x5d\x6f\xd3\x30\x14\x7d\xae\x7f\xc5\xc5\xa8\x52\xb2\x65\x0e\x13\x6f\x65\x7b\x00\xb4\x89\x49\x30\xa6\x09\x0d\x89\x69\x9a\xbc\xe4\x26\xb1\x1a\x7f\x70\xed\x54\x9d\xda\xfc\x77\xe4\x34\x1b\x43\x19\x08\x81\x5f\x5a\xdb\xe7\xde\x73\x8e\xcf\x4d\x9e\xc3\x7b\x5b\x22\xd4\x68\x90\x64\xc0\x12\xee\xee\xc1\x91\x0d\xb6\x38\xa8\xd1\x1c\x68\x55\x90\x65\x4e\x16\x4b\x59\x23\x6c\x36\xb5\xbd\x58\xd6\x1f\xa5\x0f\x27\x2d\x6a\x34\x01\xc4\xa9\x6a\x11\xb6\xe0\x5d\xab\xc2\x5b\x22\x79\x0f\xfc\x0d\x87\x2d\xb4\xd2\x07\xd8\x02\xa1\x6b\x65\x81\xc0\x05\x07\x7e\xcb\xfb\x9e\x31\xa5\x9d\xa5\x00\x09\x03\xe0\x95\x0e\x3c\xfe\x1a\x0c\x79\x13\x82\x1b\x36\x84\x55\x8b\xc5\xee\xc2\x07\x52\xa6\xf6\x9c\xb1\x19\xaf\x55\x68\xba\x3b\x51\x58\x9d\xd7\x96\x54\xdb\xca\x5c\x77\x6b\xce\x66\x83\xca\x5b\xe9\x14\x3c\xc5\x74\x46\xf9\x20\x8b\xe5\x81\xa5\x3a\x1f\x20\xf9\xea\x75\x2e\x9d\xe2\x2c\x65\xac\xea\x4c\x01\x97\x58\x2b\x1f\x90\x12\x82\x3d\xdd\xad\xc5\xa5\xed\x02\x52\x06\x0d\x28\x13\x90\x2a\x59\xe0\xa6\xcf\x00\x9d\x87\xeb\x9b\xbd\x47\x1a\x71\x62\x4a\x67\x95\x09\x29\x20\x91\x25\xd8\x30\x18\xd7\x0a\x16\xc7\x30\x1a\x10\x57\xb2\xed\xf0\x73\x95\x34\x29\x7b\x04\xa8\x0a\x56\xe2\xbc\xd3\x9f\x30\x34\xb6\x4c\x52\x38\x82\xc3\x27\xf5\x0f\x8b\x30\x74\x64\xa0\xd2\x41\x9c\x44\x8a\x2a\xe1\x8d\x34\x65\x8b\x04\x8d\xf4\x60\x2c\xe8\xa1\x81\x5f\xc0\xfc\x0b\xcf\xa0\x49\x1f\x5b\xf4\x3f\xc9\x2a\x4b\x70\x1b\xf5\x0f\xaa\xa4\xa9\x71\xf0\x32\xa5\x53\xe5\x3a\x42\xc6\xd7\x16\x67\xa6\xc4\x75\x82\x4e\x9c\x4b\x8d\x59\x0c\x2f\x9d\x96\x54\x43\x55\x94\xbf\xdd\x42\x8b\xe6\x01\x9f\xc2\xd1\xf1\x70\x35\xa5\xf9\x83\x3b\x65\x56\xb2\x55\x25\x3c\x7d\x5e\x30\x52\xe3\x02\xe6\x9e\x47\x0f\xbb\xe6\x93\x9e\xfd\xe4\x24\x56\x45\x37\x63\xc9\xb5\x2a\xd7\xfb\x87\x8b\x9b\x09\x4e\x47\xd0\x4a\xec\x92\x78\x77\x1f\xb1\x89\x79\x96\x43\x55\xf0\x42\x8b\x33\x7f\x15\x25\x26\x69\x34\x1c\xb7\xdf\x90\x6c\x92\xfe\x9b\xcf\x31\xcd\x6c\xcc\x11\xe6\x31\xd4\x00\x95\xed\x4c\xc9\x33\x30\xbf\xf1\x3a\x1d\x94\x26\x03\xbb\x8c\x4e\xb4\x38\x7b\x18\xda\x24\x15\x49\x1c\xf0\x24\x7e\x53\xe2\x12\xbd\xb3\xc6\xe3\x57\x52\xc3\x74\xef\x8d\xa7\xdf\x3b\xf4\x21\x7d\xde\xad\x5d\xfe\x97\xad\x05\xcc\x5f\xee\xaf\x78\xf6\xab\xa8\xbf\x09\x8f\xc4\x87\xa1\xc5\x69\x94\x8f\x4e\x5c\xc8\xd0\x5c\xbf\xba\xc9\x80\x9a\x74\x8c\xca\xc7\xf3\xdd\x5f\x21\x44\x3a\x84\x9c\x4c\xe6\xa3\x67\xb3\xd9\x6c\x54\x6a\x54\xcb\x7a\xf6\x23\x00\x00\xff\xff\x0f\xc4\x43\xa4\xed\x04\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyb19ncnBjLnBiLmdvLnRtcGw="].(os.FileInfo),
		fs["/e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyb19odHRwLnBiLmdvLnRtcGw="].(os.FileInfo),
		fs["/e3t0cmltU3VmZml4ICIucHJvdG8iIC5GaWxlLk5hbWV9fV9taWNyby5wYi5nby50bXBs"].(os.FileInfo),
		fs["/micro_chi.pb.go.tmpl"].(os.FileInfo),
		fs["/micro_gorilla.pb.go.tmpl"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}

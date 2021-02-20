package auxiliary

import (
	"net/url"
	"os"
	"strings"
	"unicode"
)

// URI represents the full URI for a file.
type URI string

const fileScheme = "file"

func (uri URI) IsFile() bool {
	return strings.HasPrefix(string(uri), "file://")
}

//PathExists 判断文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//URIFromURI
func URIFromURI(s string) URI {
	if !strings.HasPrefix(s, "file://") {
		return URI(s)
	}

	if !strings.HasPrefix(s, "file:///") {
		// VS Code sends URLs with only two slashes, which are invalid. golang/go#39789.
		s = "file:///" + s[len("file://"):]
	}
	// Even though the input is a URI, it may not be in canonical form. VS Code
	// in particular over-escapes :, @, etc. Unescape and re-encode to canonicalize.
	path, err := url.PathUnescape(s[len("file://"):])
	if err != nil {
		panic(err)
	}

	// File URIs from Windows may have lowercase drive letters.
	// Since drive letters are guaranteed to be case insensitive,
	// we change them to uppercase to remain consistent.
	// For example, file:///c:/x/y/z becomes file:///C:/x/y/z.
	if isWindowsDriveURIPath(path) {
		path = path[:1] + strings.ToUpper(string(path[1])) + path[2:]
	}
	u := url.URL{Scheme: fileScheme, Path: path}
	return URI(u.String())
}

func isWindowsDriveURIPath(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}

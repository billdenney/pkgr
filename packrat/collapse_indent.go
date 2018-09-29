package packrat

import "bytes"

// CollapseIndentation collapses the pattern of newline then tab
// as it makes it easier for simple parsers to understand
// instead of having to look ahead/behind
func CollapseIndentation(b []byte) []byte {
	return bytes.Replace(b, []byte("\n\t"), []byte(""), -1)
}

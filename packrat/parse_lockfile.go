package packrat

import (
	"bytes"
	"fmt"
	"strings"
)

func trimmedString(b []byte) string {
	return strings.Trim(string(b), " ")
}

// ParsePackageReqs parses package reqs to a struct
func ParsePackageReqs(b []byte) PackageReqs {
	pr := PackageReqs{}
	lines := bytes.Split(b, []byte("\n"))
	for i, f := range lines {
		if len(f) == 0 {
			// empty line
			continue
		}
		fe := bytes.SplitN(f, []byte(":"), 2)
		if len(fe) != 2 {
			fmt.Printf("bad line %d nelements %d : %s", i, len(fe), f)
			continue
		}
		switch {
		case bytes.Equal(fe[0], []byte("Package")):
			pr.Package = trimmedString(fe[1])
		case bytes.Equal(fe[0], []byte("Version")):
			pr.Version = trimmedString(fe[1])
		case bytes.Equal(fe[0], []byte("Source")):
			pr.Source = trimmedString(fe[1])
		case bytes.Equal(fe[0], []byte("Hash")):
			pr.Hash = trimmedString(fe[1])
		case bytes.Equal(fe[0], []byte("Requires")):
			pr.Requires = strings.Fields(string(bytes.Replace(fe[1], []byte(","), []byte(" "), -1)))
		default:
			fmt.Printf("unrecognized fields %s-%s\n", fe[0], fe[1])
		}
	}
	return pr
}

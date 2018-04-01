package fdep

import (
	"fmt"
	"path"
	"strings"

	"github.com/RangelReale/fproto"
)

// The dependency file type.
type DepFileType int

const (
	// Your own project's proto files.
	DepType_Own DepFileType = iota

	// Imported proto file, that are not part of your project.
	DepType_Imported
)

func (dt DepFileType) String() string {
	switch dt {
	case DepType_Own:
		return "OWN"
	case DepType_Imported:
		return "IMPORTED"
	default:
		return "UNKNOWN"
	}
}

// DepFile represents one .proto file into the dependency.
type DepFile struct {
	// The INTERNAL path of the .proto file, for example "google/protobuf/empty.proto"
	// This is NOT the filesystem path
	FilePath string

	// The type of the file dependency, whether it is your own file, or an imported one.
	DepType DepFileType

	// The parent dependency list this file is contained.
	Dep *Dep

	// The parsed proto file. Can be nil it was from an ignored dependency.
	ProtoFile *fproto.ProtoFile
}

// Returns one named type from the dependency, in relation to the current file.
// If the type is from the current file, the "Alias" field is blank.
//
// If multiple types are found for the same name, an error is issued.
// If there is this possibility, use the GetTypes method instead.
//
// May return nil if type not found.
func (df *DepFile) FindType(name string) (*DepType, error) {
	t, err := df.GetTypes(name)
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, nil
	} else if len(t) > 1 {
		return nil, fmt.Errorf("More than one type found for '%s'", name)
	}

	return t[0], nil
}

// Like FindType, but returns an error if not found
func (df *DepFile) GetType(name string) (*DepType, error) {
	t, err := df.FindType(name)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, fmt.Errorf("Type %s not found", name)
	}
	return t, nil
}

// Returns all named types from the dependency, in relation to the current file.
// If the type is from the current file, the "Alias" field is blank.
//
// If not found using the name, the current file's package is searched recursivelly
// appending the name.
//
// Use this method if there is a possibility that one name resolves to more than one type.
func (df *DepFile) GetTypes(name string) ([]*DepType, error) {
	t, err := df.Dep.internalGetTypes(name, df)
	if err != nil {
		return nil, err
	}
	if len(t) > 0 {
		return t, nil
	}

	// if not found, search in each dotted scope of the current file's package
	if df.ProtoFile != nil && df.ProtoFile.PackageName != "" {
		scopes := strings.Split(df.ProtoFile.PackageName, ".")
		for si := 1; si <= len(scopes); si++ {
			t, err = df.Dep.internalGetTypes(fmt.Sprintf("%s.%s", strings.Join(scopes[:si], "."), name), df)
			if err != nil {
				return nil, err
			}

			if t != nil {
				return t, nil
			}
		}
	}

	// Not found in any method
	return nil, nil
}

func (df *DepFile) GetFileOfName(name string) (*DepFileOfName, error) {
	return df.Dep.GetFileOfName(name)
}

func (df *DepFile) GetFilesOfName(name string) ([]*DepFileOfName, error) {
	return df.Dep.GetFilesOfName(name)
}

// Checks if the passed DepFile refers to the same file as this one.
func (df *DepFile) IsSame(depfile *DepFile) bool {
	if depfile == nil {
		return false
	}

	if df == depfile {
		return true
	}

	if df.FilePath == depfile.FilePath && df.ProtoFile.PackageName == depfile.ProtoFile.PackageName {
		return true
	}

	return false
}

func (df *DepFile) OriginalAlias() string {
	if df.ProtoFile != nil {
		return df.ProtoFile.PackageName
	}
	return ""
}

// Checks if the passed DepFile refers to the same package as this one.
func (df *DepFile) IsSamePackage(depfile *DepFile) bool {
	if df == depfile {
		return true
	}

	if path.Dir(df.FilePath) == path.Dir(depfile.FilePath) && df.ProtoFile.PackageName == depfile.ProtoFile.PackageName {
		return true
	}

	return false
}

// Returns the go package of the file. If there is no "go_package" option, returns the "path" part of the package name.
func (df *DepFile) GoPackage() string {
	for _, o := range df.ProtoFile.Options {
		if o.Name == "go_package" {
			return o.Value.String()
		}
	}
	return path.Dir(df.ProtoFile.PackageName)
}

// Find all dependencies of file, include public ones from imports.
func (df *DepFile) FindDependencies() []string {
	var ret []string
	if df.ProtoFile != nil {
		ret = append(ret, df.ProtoFile.Dependencies...)
		for _, pd := range df.ProtoFile.PublicDependencies {
			if pdf, ispdf := df.Dep.Files[pd]; ispdf && pdf.ProtoFile != nil {
				ret = append(ret, pdf.ProtoFile.PublicDependencies...)
			}
		}
	}
	return ret
}

// Result of GetFilesOfName
type DepFileOfName struct {
	// File
	DepFile *DepFile
	// Package name
	Package string
	// Rest of name excluding the package name
	Name string
}

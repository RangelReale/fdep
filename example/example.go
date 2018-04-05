package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
)

var (
	lines = strings.Repeat("=", 20)
)

// This directory must be the working directory
func main() {
	// load test files
	pdep := fdep.NewDep()

	// add include path
	err := pdep.AddIncludeDir("proto_test/include")
	if err != nil {
		log.Fatal(err)
	}

	// add application files, with root as "app"
	// we use the "WithRoot" version, because if we added "proto_test" directly,
	// the "include" directory would be also included, but with a wrong path
	// (prepended by "include")
	err = pdep.AddPathWithRoot("app", "proto_test/app", fdep.DepType_Own)
	if err != nil {
		log.Fatal(err)
	}

	printFiles(pdep)

	printPackages(pdep)

	printExtensions(pdep)

	printTypes(pdep)

	printFields(pdep)

	printFieldTypes(pdep, "app.core.SendMailAttach")

	printFieldTypes(pdep, "fproto_wrap_headers.Headers")

	printMessageExtensions(pdep, fdep.FIELD_OPTION.MessageName())

	printOptionType(pdep)
}

// OUTPUT:
// ==================== PRINT FILES ====================
// File: app/core/sendmail.proto (OWN) [package: app.core]
// File: fproto-wrap/time.proto (IMPORTED) [package: fproto_wrap]
// File: google/protobuf/timestamp.proto (IMPORTED) [package: google.protobuf]
// File: google/protobuf/descriptor.proto (IMPORTED) [package: google.protobuf]
// File: google/protobuf/empty.proto (IMPORTED) [package: google.protobuf]
// File: app/base/pagination.proto (OWN) [package: app.base]
// File: fproto-wrap/uuid.proto (IMPORTED) [package: fproto_wrap]
// File: fproto-wrap/jsontag.proto (IMPORTED) [package: fproto_wrap]
// File: fproto-wrap-headers/headers.proto (IMPORTED) [package: fproto_wrap_headers]
// File: fproto-wrap-validate/validate.proto (IMPORTED) [package: validate]
// File: app/core/user.proto (OWN) [package: app.core]
func printFiles(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT FILES %s\n", lines, lines)
	for filepath, file := range pdep.Files {
		fmt.Printf("File: %s (%s) [package: %s]\n", filepath, file.DepType.String(), file.ProtoFile.PackageName)
	}
}

// OUTPUT:
// ==================== PRINT PACKAGES ====================
// Package: app.base [files: app/base/pagination.proto]
// Package: fproto_wrap [files: fproto-wrap/uuid.proto, fproto-wrap/time.proto, fproto-wrap/jsontag.proto]
// Package: google.protobuf [files: google/protobuf/timestamp.proto, google/protobuf/descriptor.proto, google/protobuf/empty.proto]
// Package: fproto_wrap_headers [files: fproto-wrap-headers/headers.proto]
// Package: validate [files: fproto-wrap-validate/validate.proto]
// Package: app.core [files: app/core/sendmail.proto, app/core/user.proto]
func printPackages(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT PACKAGES %s\n", lines, lines)
	for pkg, filelist := range pdep.Packages {
		fmt.Printf("Package: %s [files: %s]\n", pkg, strings.Join(filelist, ", "))
	}
}

// OUTPUT:
// ==================== PRINT EXTENSIONS ====================
// Extension: google.protobuf.FieldOptions [packages: fproto_wrap, validate]
func printExtensions(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT EXTENSIONS %s\n", lines, lines)
	for ext, pkglist := range pdep.Extensions {
		fmt.Printf("Extension: %s [packages: %s]\n", ext, strings.Join(pkglist, ", "))
	}
}

// OUTPUT:
// ==================== PRINT TYPES ====================
// Type 'app.core.User' is in file 'app/core/user.proto', package 'app.core' [name: User]
// Type 'app.core.SendMail.Body' is in file 'app/core/sendmail.proto', package 'app.core' [name: SendMail.Body, alias: app.core]
// Type 'Body' in the context of 'app.core.SendMail' is in file 'app/core/sendmail.proto', package 'app.core' [name: SendMail.Body, alias: ]
func printTypes(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT TYPES %s\n", lines, lines)

	//
	// app.core.User
	//
	tp_user, err := pdep.GetType("app.core.User")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Type 'app.core.User' is in file '%s', package '%s' [name: %s]\n",
		tp_user.DepFile.FilePath, tp_user.OriginalAlias, tp_user.Name)

	//
	// app.core.SendMail.Body
	//
	tp_sendmail_body, err := pdep.GetType("app.core.SendMail.Body")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Type 'app.core.SendMail.Body' is in file '%s', package '%s' [name: %s, alias: %s]\n",
		tp_sendmail_body.DepFile.FilePath, tp_sendmail_body.OriginalAlias, tp_sendmail_body.Name, tp_sendmail_body.Alias)

	//
	// app.core.SendMail.Body in the context of app.Core.Sendmail
	//
	tp_sendmail, err := pdep.GetType("app.core.SendMail")
	if err != nil {
		log.Fatal(err)
	}

	tp_sendmail_body2, err := tp_sendmail.GetType("Body")
	if err != nil {
		log.Fatal(err)
	}

	// When getting a type in the context of other type, the alias may be blank if the type is on the same file.
	fmt.Printf("Type 'Body' in the context of 'app.core.SendMail' is in file '%s', package '%s' [name: %s, alias: %s]\n",
		tp_sendmail_body2.DepFile.FilePath, tp_sendmail_body2.OriginalAlias, tp_sendmail_body2.Name, tp_sendmail_body2.Alias)
}

// OUTPUT:
// ==================== PRINT FIELDS ====================
// MESSAGE ELEMENT NAME: SendMail
// * FIELD: sendmail_id - TYPE: fproto_wrap.UUID
// * FIELD: sent_opt - TYPE: oneof
// * FIELD: tries - TYPE: int32
// * FIELD: error_message - TYPE: string
// * FIELD: last_try_at - TYPE: fproto_wrap.NullTime
// * FIELD: destination_to - TYPE: SendMailDestination
// * FIELD: destination_cc - TYPE: SendMailDestination
// * FIELD: destination_bcc - TYPE: SendMailDestination
// * FIELD: subject - TYPE: string
// * FIELD: body - TYPE: Body
// * FIELD: attach - TYPE: SendMailAttach
func printFields(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT FIELDS %s\n", lines, lines)

	//
	// app.core.Sendmail
	//
	tp_sendmail, err := pdep.GetType("app.core.SendMail")
	if err != nil {
		log.Fatal(err)
	}

	message_element, is_message_element := tp_sendmail.Item.(*fproto.MessageElement)
	if !is_message_element {
		log.Fatal("Should be a *fproto.MessageElement")
	}

	fmt.Printf("MESSAGE ELEMENT NAME: %s\n", message_element.Name)
	for _, fld := range message_element.Fields {
		switch xfld := fld.(type) {
		case *fproto.FieldElement:
			fmt.Printf("* FIELD: %s - TYPE: %s\n", xfld.Name, xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: map<%s, %s>\n", xfld.Name, xfld.KeyType, xfld.Type)
		case *fproto.OneOfFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: oneof\n", xfld.Name)
		}
	}
}

// OUTPUT 1:
// ==================== PRINT FIELD TYPES: app.core.SendMailAttach ====================
// Message type name: app.core.SendMailAttach
// * FIELD: attach_type - TYPE: attach_type
// TYPE: app.core.SendMailAttach.attach_type [ENUM] - from file: app/core/sendmail.proto
// * FIELD: content_type - TYPE: string
// TYPE: SCALAR string
// * FIELD: headers - TYPE: fproto_wrap_headers.Headers
// TYPE: fproto_wrap_headers.Headers [MESSAGE] - from file: fproto-wrap-headers/headers.proto
// * FIELD: filename - TYPE: string
// TYPE: SCALAR string
// * FIELD: content_opt - TYPE: oneof
// ONEOF FIELD: download_url
// ONEOF FIELD: content
//
// OUTPUT 2:
// ==================== PRINT FIELD TYPES: fproto_wrap_headers.Headers ====================
// Message type name: fproto_wrap_headers.Headers
// * FIELD: headers - TYPE: map<string, Values>
// KEY TYPE: SCALAR string
// TYPE: fproto_wrap_headers.Headers.Values [MESSAGE] - from file: fproto-wrap-headers/headers.proto
func printFieldTypes(pdep *fdep.Dep, typeName string) {
	fmt.Printf("%s PRINT FIELD TYPES: %s %s\n", lines, typeName, lines)

	tp_print, err := pdep.GetType(typeName)
	if err != nil {
		log.Fatal(err)
	}

	message_element, is_message_element := tp_print.Item.(*fproto.MessageElement)
	if !is_message_element {
		log.Fatal("Should be a *fproto.MessageElement")
	}

	fmt.Printf("Message type name: %s\n", tp_print.FullOriginalName())

	for _, fld := range message_element.Fields {
		switch xfld := fld.(type) {
		case *fproto.FieldElement:
			fmt.Printf("* FIELD: %s - TYPE: %s\n", xfld.Name, xfld.Type)
			Util_PrintType("", tp_print, xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: map<%s, %s>\n", xfld.Name, xfld.KeyType, xfld.Type)
			Util_PrintType("KEY ", tp_print, xfld.KeyType)
			Util_PrintType("", tp_print, xfld.Type)
		case *fproto.OneOfFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: oneof\n", xfld.Name)
			for _, oofld := range xfld.Fields {
				fmt.Printf("\tONEOF FIELD: %s\n", oofld.FieldName())
			}
		}
	}
}

// OUTPUT:
// ==================== PRINT MESSAGE EXTENSIONS: google.protobuf.FieldOptions ====================
// Message type name: google.protobuf.FieldOptions
// * EXTENSION: validate.google.protobuf.FieldOptions [package: validate] [proto path: fproto-wrap-validate/validate.proto]
// FIELD: field [type: FieldValidator]
// * EXTENSION: fproto_wrap.google.protobuf.FieldOptions [package: fproto_wrap] [proto path: fproto-wrap/jsontag.proto]
// FIELD: jsontag [type: JSONTag]
func printMessageExtensions(pdep *fdep.Dep, typeName string) {
	fmt.Printf("%s PRINT MESSAGE EXTENSIONS: %s %s\n", lines, typeName, lines)

	tp_print, err := pdep.GetType(typeName)
	if err != nil {
		log.Fatal(err)
	}

	extensions, err := tp_print.GetTypeExtensions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message type name: %s\n", tp_print.FullOriginalName())

	for extpkg, extType := range extensions {
		fmt.Printf("* EXTENSION: %s [package: %s] [proto path: %s]\n", extType.FullOriginalName(), extpkg, extType.DepFile.FilePath)
		if em, emok := extType.Item.(*fproto.MessageElement); emok {
			for _, emf := range em.Fields {
				switch xfld := emf.(type) {
				case *fproto.FieldElement:
					fmt.Printf("\tFIELD: %s [type: %s]\n", xfld.Name, xfld.Type)
				case *fproto.MapFieldElement:
					fmt.Printf("\tMAP FIELD: %s [type: map<%s, %s>]\n", xfld.Name, xfld.KeyType, xfld.Type)
				case *fproto.OneOfFieldElement:
					fmt.Printf("\tFIELD: %s [type: oneof]\n", xfld.Name)
				}
			}
		}
	}
}

// OUTPUT:
// ==================== PRINT OPTION TYPE ====================
// OPTION: validate.field
// Source option type: google.protobuf.FieldOptions
// Option type: validate.google.protobuf.FieldOptions
// Option name: field
// Field item fieldname: field [type: FieldValidator]
// TYPE: validate.FieldValidator [MESSAGE] - from file: fproto-wrap-validate/validate.proto
func printOptionType(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT OPTION TYPE %s\n", lines, lines)

	o, err := pdep.GetOption(fdep.FIELD_OPTION, "validate.field")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OPTION: %s\n", o.OptionName)

	fmt.Printf("Source option type: %s\n", o.SourceOption.FullOriginalName())
	if o.Option != nil {
		fmt.Printf("Option type: %s\n", o.Option.FullOriginalName())
	}
	if o.Name != "" {
		fmt.Printf("Option name: %s\n", o.Name)
	}
	if o.FieldItem != nil {
		parent_dt := pdep.DepTypeFromElement(o.FieldItem.ParentElement())

		switch xfld := o.FieldItem.(type) {
		case *fproto.FieldElement:
			fmt.Printf("Field item fieldname: %s [type: %s]\n", xfld.Name, xfld.Type)
			Util_PrintType("", parent_dt, xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("Field item fieldname: %s [type: map<%s, %s>]\n", xfld.Name, xfld.KeyType, xfld.Type)
			Util_PrintType("KEY ", parent_dt, xfld.KeyType)
			Util_PrintType("", parent_dt, xfld.Type)
		case *fproto.OneOfFieldElement:
			fmt.Printf("Field item fieldname: %s [type: oneof]\n", xfld.Name)
		}
	}
}

func Util_PrintType(desc string, parent *fdep.DepType, typeName string) {
	tp_fld, err := parent.GetType(typeName)
	if err != nil {
		log.Fatal(err)
	}
	if tp_fld.IsScalar() {
		fmt.Printf("\t%sTYPE: SCALAR %s\n", desc, tp_fld.ScalarType.ProtoType())
	} else {
		fmt.Printf("\t%sTYPE: %s [%s] - from file: %s\n", desc, tp_fld.FullOriginalName(), tp_fld.Item.ElementTypeName(), tp_fld.DepFile.FilePath)
	}
}

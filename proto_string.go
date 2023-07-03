package proto2

import (
	"fmt"
	"go/format"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GoStruct returns a Go struct literal for the given proto.Message.
func GoStruct(m proto.Message) string {
	if m == nil {
		return ""
	}
	var builder strings.Builder
	nestedFields := &nestingStack{}

	push := func(p protopath.Values) error {
		last := p.Index(-1)
		beforeLast := p.Index(-2)
		var fd protoreflect.FieldDescriptor

		switch last.Step.Kind() {
		case protopath.FieldAccessStep:
			fd = last.Step.FieldDescriptor()
			if oneof := oneofDescriptor(fd); oneof != nil {
				parentType, fieldName := parentAndField(fd)
				// oneof fields are expressed as parentType_fieldName in Go
				builder.WriteString(fmt.Sprintf("%s: &%s_%s{\n%s: ", oneof.Name(), parentType, fieldName, fieldName))
				// oneof needs an extra level of nesting
				nestedFields.push()
			} else {
				builder.WriteString(fmt.Sprintf("%s: ", fd.Name()))
			}

		case protopath.ListIndexStep:
			// lists always appear in the context of a repeated field
			fd = beforeLast.Step.FieldDescriptor()

		case protopath.MapIndexStep:
			// maps always appear in the context of a repeated field
			fd = beforeLast.Step.FieldDescriptor()
			builder.WriteString(fmt.Sprintf("%v: ", last.Step.MapIndex().Interface()))
		}

		switch v := last.Value.Interface().(type) {
		case protoreflect.EnumNumber:
			var fullName protoreflect.FullName
			var enumName protoreflect.Name
			if enumDesc := fd.Enum(); enumDesc != nil {
				fullName = enumDesc.FullName()
				enumValue := enumDesc.Values().ByNumber(v)
				if enumValue != nil {
					enumName = enumValue.Name()
				}
			}
			if enumName != "" {
				// enum fields are expressed as fullName_enumName in Go
				builder.WriteString(fmt.Sprintf("%s_%s", fullName, enumName))
			} else {
				builder.WriteString(fmt.Sprintf("%d", v))
			}
			builder.WriteString(nestedFields.commaUnlessLast())

		case protoreflect.Message:
			if last.Step.Kind() == protopath.ListIndexStep || last.Step.Kind() == protopath.MapIndexStep {
				builder.WriteString("{\n")
			} else {
				builder.WriteString(fmt.Sprintf("&%s{\n", v.Descriptor().FullName()))
			}
			nestedFields.push()

		case protoreflect.List:
			builder.WriteString(fmt.Sprintf("[]%s{\n", goTypeName(fd)))
			nestedFields.push()

		case protoreflect.Map:
			builder.WriteString(fmt.Sprintf("map[%s]%s{\n", fd.MapKey().Kind(), goTypeName(fd.MapValue())))
			nestedFields.push()

		case []byte:
			builder.WriteString(fmt.Sprintf("[]byte{%s}", bytesToStringList(v)))
			builder.WriteString(nestedFields.commaUnlessLast())

		case string:
			builder.WriteString(fmt.Sprintf("%q", v))
			builder.WriteString(nestedFields.commaUnlessLast())

		default:
			builder.WriteString(fmt.Sprintf("%v", v))
			builder.WriteString(nestedFields.commaUnlessLast())
		}
		return nil
	}

	pop := func(p protopath.Values) error {
		last := p.Index(-1)
		switch last.Value.Interface().(type) {
		case protoreflect.Message, protoreflect.List, protoreflect.Map:
			nestedFields.pop()
			builder.WriteString("}")
			builder.WriteString(nestedFields.commaUnlessLast())
		}
		if oneof := oneofDescriptor(last.Step.FieldDescriptor()); oneof != nil {
			nestedFields.pop()
			builder.WriteString("}")
			builder.WriteString(nestedFields.commaUnlessLast())
		}
		return nil
	}

	protorange.Options{Stable: true}.Range(m.ProtoReflect(), push, pop)

	// Format the generated Go code
	s, err := format.Source([]byte(builder.String()))
	if err != nil {
		panic(fmt.Errorf("error formatting generated code:\n%s\n %v", builder.String(), err))
	}
	return string(s)
}

// oneofDescriptor returns the oneof descriptor for the given field descriptor, if it has one.
func oneofDescriptor(fd protoreflect.FieldDescriptor) protoreflect.OneofDescriptor {
	if fd == nil {
		return nil
	}
	if oneof := fd.ContainingOneof(); oneof != nil && !oneof.IsSynthetic() {
		return oneof
	}
	return nil
}

func bytesToStringList(v []byte) string {
	byteAsDecString := make([]string, len(v))
	for i, b := range v {
		byteAsDecString[i] = fmt.Sprintf("%d", b)
	}
	return strings.Join(byteAsDecString, ", ")
}

// goTypeName returns the Go type name for the given field descriptor.
func goTypeName(fd protoreflect.FieldDescriptor) string {
	kind := fd.Kind()
	typeName := kind.String()
	if kind == protoreflect.MessageKind {
		typeName = fmt.Sprintf("*%s", fd.Message().FullName())
	}
	return typeName
}

// parentAndField returns the parent type and an exported field name for the given field descriptor.
func parentAndField(fd protoreflect.FieldDescriptor) (parentType, fieldName string) {
	return string(fd.FullName().Parent()), export(string(fd.Name()))
}

// export returns the given identifier so that it can be used as a Go identifier.
func export(s string) string { return strings.ToUpper(s[:1]) + s[1:] }

type nestingStack struct {
	stack []bool
}

func (s *nestingStack) push() {
	s.stack = append(s.stack, false)
}

func (s *nestingStack) pop() {
	if len(s.stack) > 0 {
		s.stack[len(s.stack)-1] = true
		s.stack = s.stack[:len(s.stack)-1]
	}
}

func (s *nestingStack) commaUnlessLast() string {
	if len(s.stack) > 0 && !s.stack[len(s.stack)-1] {
		return ",\n"
	}
	return "\n"
}

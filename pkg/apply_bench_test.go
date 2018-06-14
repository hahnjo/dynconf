// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"testing"
)

func BenchmarkApply_Empty(b *testing.B) {
	r := Recipe{}
	in := []byte{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(r, in)
	}
}

func BenchmarkApply_Delete(b *testing.B) {
	r := Recipe{
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
	}
	r.Compile()
	in := []byte("remove")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(r, in)
	}
}

func BenchmarkApply_Replace(b *testing.B) {
	r := Recipe{
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	r.Compile()
	in := []byte("replace")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(r, in)
	}
}

func BenchmarkApply_Append(b *testing.B) {
	r := Recipe{
		Append: "append",
	}
	in := []byte{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(r, in)
	}
}

func BenchmarkApply_Large(b *testing.B) {
	r := Recipe{
		Append: "append",
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	r.Compile()
	in := make([]byte, 0)
	for i := 0; i < 1000; i++ {
		in = append(in, "remove\n"...)
		in = append(in, "line\n"...)
		in = append(in, "text search text\n"...)
		in = append(in, "line\n"...)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(r, in)
	}
}

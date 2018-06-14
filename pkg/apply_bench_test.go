// SPDX-License-Identifier:	GPL-3.0-or-later

package dynconf

import (
	"testing"
)

func BenchmarkApply_Empty(b *testing.B) {
	c := Config{}
	in := []byte{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(c, in)
	}
}

func BenchmarkApply_Delete(b *testing.B) {
	c := Config{
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
	}
	c.Compile()
	in := []byte("remove")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(c, in)
	}
}

func BenchmarkApply_Replace(b *testing.B) {
	c := Config{
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	c.Compile()
	in := []byte("replace")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(c, in)
	}
}

func BenchmarkApply_Append(b *testing.B) {
	c := Config{
		Append: "append",
	}
	in := []byte{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(c, in)
	}
}

func BenchmarkApply_Large(b *testing.B) {
	c := Config{
		Append: "append",
		Delete: []DeleteEntry{
			{Search: "remove"},
		},
		Replace: []ReplaceEntry{
			{Search: "search", Replace: "replace"},
		},
	}
	c.Compile()
	in := make([]byte, 0)
	for i := 0; i < 1000; i++ {
		in = append(in, "remove\n"...)
		in = append(in, "line\n"...)
		in = append(in, "text search text\n"...)
		in = append(in, "line\n"...)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ApplyToInput(c, in)
	}
}

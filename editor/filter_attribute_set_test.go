package editor

import (
	"testing"
)

func TestAttributeSetFilter(t *testing.T) {
	cases := []struct {
		name    string
		src     string
		address string
		value   string
		ok      bool
		want    string
	}{
		{
			name: "simple top level attribute (reference)",
			src: `
a0 = v0
a1 = v1
`,
			address: "a0",
			value:   "v2",
			ok:      true,
			want: `
a0 = v2
a1 = v1
`,
		},
		{
			name: "simple top level attribute (string literal)",
			src: `
a0 = "v0"
a1 = "v1"
`,
			address: "a0",
			value:   `"v2"`,
			ok:      true,
			want: `
a0 = "v2"
a1 = "v1"
`,
		},
		{
			name: "simple top level attribute (number literal)",
			src: `
a0 = 0
a1 = 1
`,
			address: "a0",
			value:   "2",
			ok:      true,
			want: `
a0 = 2
a1 = 1
`,
		},
		{
			name: "simple top level attribute (bool literal)",
			src: `
a0 = true
a1 = true
`,
			address: "a0",
			value:   "false",
			ok:      true,
			want: `
a0 = false
a1 = true
`,
		},
		{
			name: "simple top level attribute (with comments)",
			src: `
// before attr
a0 = "v0" // inline
a1 = "v1"
`,
			address: "a0",
			value:   `"v2"`,
			ok:      true,
			want: `
// before attr
a0 = "v2" // inline
a1 = "v1"
`,
		},
		{
			name: "attribute in block",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a1",
			value:   "v2",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v2
}
`,
		},
		{
			name: "top level attribute not found",
			src: `
a0 = v0
`,
			address: "a1",
			value:   "v2",
			ok:      true,
			want: `
a0 = v0
`,
		},
		{
			name: "attribute not found in block",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b1.l1.a2",
			value:   "v2",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
		},
		{
			name: "block not found",
			src: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
			address: "b2.l1.a1",
			value:   "v2",
			ok:      true,
			want: `
a0 = v0
b1 "l1" {
  a1 = v1
}
`,
		},
		{
			name: "escaped address",
			src: `
a0 = v0
b1 "l.1" {
  a1 = v1
}
`,
			address: `b1.l\.1.a1`,
			value:   "v2",
			ok:      true,
			want: `
a0 = v0
b1 "l.1" {
  a1 = v2
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewEditOperator(NewAttributeSetFilter(tc.address, tc.value))
			output, err := o.Apply([]byte(tc.src), "test")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

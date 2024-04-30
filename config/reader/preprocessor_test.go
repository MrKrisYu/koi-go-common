package reader

import (
	"os"
	"strings"
	"testing"
)

func TestReplaceEnvVars(t *testing.T) {
	os.Setenv("myBar", "cat")

	testData := []struct {
		expected string
		data     []byte
	}{
		// Right use cases
		{
			`{"foo": "bar", "baz": {"bar": "cat"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${myBar}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "123"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${ABC:123}"}}`),
		},
		{
			`{"foo": "bar", "baz": {"bar": "22"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${CDB:22}"}}`),
		},
		// Wrong use cases
		{
			`{"foo": "bar", "baz": {"bar": "${myBar-}"}}`,
			[]byte(`{"foo": "bar", "baz": {"bar": "${myBar-}"}}`),
		},
	}

	for _, test := range testData {
		res, err := ReplaceEnvVars(test.data)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Compare(test.expected, string(res)) != 0 {
			t.Fatalf("Expected %s got %s", test.expected, res)
		}
	}
}

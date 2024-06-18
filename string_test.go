package copyfile

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/exp/maps"
	"gotest.tools/v3/assert"
)

func TestParse(t *testing.T) {
	for _, test := range []struct {
		name           string
		str            string
		expectedFields map[string]string
	}{
		{
			name: "simple",
			str:  "test {tenant} sample {nonna}",
			expectedFields: map[string]string{
				"tenant": "",
				"nonna":  "",
			},
		},
		{
			name: "not closed",
			str:  "test {tenant} sample {nonna",
			expectedFields: map[string]string{
				"tenant": "",
			},
		},
		{
			name: "not open",
			str:  "test tenant} sample {nonna}",
			expectedFields: map[string]string{
				"nonna": "",
			},
		},
		{
			name: "escaped open",
			str:  "test {{tenant} sample {nonna}",
			expectedFields: map[string]string{
				"nonna": "",
			},
		},
		{
			name: "escaped close",
			str:  "test {tenant}} sample {nonna}",
			expectedFields: map[string]string{
				"tenant} sample {nonna": "{tenant}} sample {nonna}",
			},
		},
		{
			name: "escaped repeated",
			str:  "test {{{tenant} sample {nonna}",
			expectedFields: map[string]string{
				"tenant": "",
				"nonna":  "",
			},
		},
		{
			name: "escaped repeated 4",
			str:  "test {{{{tenant} sample {nonna}",
			expectedFields: map[string]string{
				"nonna": "",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := Parse(test.str)
			assert.DeepEqual(t, maps.Keys(test.expectedFields), maps.Keys(p.fields), cmpopts.SortSlices(cmp.Less[string]))
			for fn, field := range p.fields {
				fv := p.str[field.start:field.end]
				fexpected := test.expectedFields[fn]
				if fexpected == "" {
					fexpected = fmt.Sprintf("{%s}", field.name)
				}
				assert.Equal(t, fexpected, fv)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	for _, test := range []struct {
		name        string
		str         string
		values      map[string]any
		expected    string
		expectedErr bool
	}{
		{
			name: "simple",
			str:  "test {tenant} sample {nonna}",
			values: map[string]any{
				"tenant": "666",
				"nonna":  "888",
			},
			expected: "test 666 sample 888",
		},
		{
			name: "missing",
			str:  "test {tenant} sample {nonna}",
			values: map[string]any{
				"tenant": "666",
			},
			expectedErr: true,
		},
		{
			name: "escape left",
			str:  "test {{tenant} sample {nonna}",
			values: map[string]any{
				"nonna": "888",
			},
			expected: "test {{tenant} sample 888",
		},
		{
			name: "escape right",
			str:  "test {tenant}} sample {nonna}",
			values: map[string]any{
				"tenant} sample {nonna": "666",
			},
			expected: "test 666",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := Parse(test.str)
			pr, err := p.Replace(test.values)
			if !test.expectedErr {
				assert.NilError(t, err)
				assert.Equal(t, test.expected, pr)
			} else {
				assert.Assert(t, err != nil)
			}
		})
	}
}

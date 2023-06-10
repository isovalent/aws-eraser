package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResourceString(t *testing.T) {
	cases := []struct {
		resourceStr string
		expected    []Resource
		isError     bool
	}{
		{
			resourceStr: "",
			expected:    nil,
			isError:     false,
		},
		{
			resourceStr: "vpc",
			expected:    nil,
			isError:     true,
		},
		{
			resourceStr: "vpc:vpc-1",
			expected: []Resource{
				{
					Type:    "vpc",
					Account: "default",
					Region:  "default",
					ID:      "vpc-1",
				},
			},
			isError: false,
		},
		{
			resourceStr: "bla:vpc-1",
			expected:    nil,
			isError:     true,
		},
		{
			resourceStr: "vpc:region-1:vpc-1",
			expected:    nil,
			isError:     true,
		},
		{
			resourceStr: "vpc:account-1:region-1:vpc-1",
			expected: []Resource{
				{
					Type:    "vpc",
					Account: "account-1",
					Region:  "region-1",
					ID:      "vpc-1",
				},
			},
			isError: false,
		},
		{
			resourceStr: "vpc:account-1:region-1:vpc-1,eks:account-2:region-2:eks-2",
			expected: []Resource{
				{
					Type:    "vpc",
					Account: "account-1",
					Region:  "region-1",
					ID:      "vpc-1",
				},
				{
					Type:    "eks",
					Account: "account-2",
					Region:  "region-2",
					ID:      "eks-2",
				},
			},
			isError: false,
		},
	}
	for _, tt := range cases {
		actual, err := parseResourceString(tt.resourceStr)
		if tt.isError {
			assert.Error(t, err, "error expected")
			continue
		}
		assert.NoError(t, err, "no error expected")
		assert.Equal(t, tt.expected, actual)
	}
}

package db

import "testing"

func TestValidateName(t *testing.T) {
	testCases := []struct {
		desc string
		name string
		err  error
	}{
		{
			desc: "Empty name",
			name: "",
			err:  ErrInvalidName,
		},
		{
			desc: "simple valid name",
			name: "database",
			err:  nil,
		},
		{
			desc: "complex valid name",
			name: "this-is-a-complex-name-0123",
			err:  nil,
		},
		{
			desc: "name with more than 63 chars",
			name: "this-has-too-many-characters-to-fit-1231231231231231231231231231",
			err:  ErrNameTooLong,
		},
		{
			desc: "name that starts with a number",
			name: "9123-asdfasdf",
			err:  ErrInvalidName,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := validateName(tC.name)
			if err != tC.err {
				t.Errorf("expected err to be '%v', got '%v'", tC.err, err)
			}
		})
	}
}

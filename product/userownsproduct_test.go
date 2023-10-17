package product

import "testing"

func TestUserOwnsProduct(t *testing.T) {
	type SampleParams struct {
		UserID        string
		ProductUserID string
		Expected      bool
	}

	validInputs := SampleParams{"1234-asdcd-8302023-ds3134-dfdf", "1234-asdcd-8302023-ds3134-dfdf", true}
	notMatchingInputs := SampleParams{"1234-asdcd-8302023-ds334-dfdf", "1234-asdcd-8302023-ds3134-dfdf", false}
	emptyInput := SampleParams{"", "", false}
	spaceInfront := SampleParams{" 1234-asdcd-8302023-ds334-dfdf", "1234-asdcd-8302023-ds3134-dfdf", false}
	trailingSpace := SampleParams{"1234-asdcd-8302023-ds334-dfdf", "1234-asdcd-8302023-ds3134-dfdf", false}

	cases := []SampleParams{
		validInputs,
		notMatchingInputs,
		emptyInput,
		spaceInfront,
		trailingSpace,
	}

	for _, item := range cases {
		result, _ := ValidateUserOwnsProduct(item.UserID, item.ProductUserID)
		if result != item.Expected {
			t.Errorf("test failed")
		}
	}
	t.Logf("all test passed")
}

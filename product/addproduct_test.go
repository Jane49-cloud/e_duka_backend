package product

import "testing"

func TestAddproduct(t *testing.T) {
	type testCase struct {
		name string
		data AddProductInput
		want bool
	}

	// cases for product name
	emptyName := testCase{
		name: "should return empty product name",
		data: AddProductInput{
			"",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	emptyName2 := testCase{
		name: "should return empty product name",
		data: AddProductInput{
			"   ",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	invalidSpecialCharactersInProductName := testCase{
		name: "should return invalida characters in name",
		data: AddProductInput{
			"productname@@%$@%$#",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}

	// test cases for product description
	emptyDescription := testCase{
		name: "should return empty product description",
		data: AddProductInput{
			"this is product name",
			"200",
			"",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	emptyDescription2 := testCase{
		name: "should return empty product description",
		data: AddProductInput{
			"product name",
			"200",
			" ",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	invalidSpecialCharactersInDescription := testCase{
		name: "should return invalida characters in name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text )*&(&^*%*^%&*of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}

	// cases for main image
	emptyImageString := testCase{
		name: "should return empty image ",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	emptyImageString2 := testCase{
		name: "should return empty image 2",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"   ",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}

	// test cases for product type
	emptyProductType := testCase{
		name: "should return empty image string",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	emptyProductType2 := testCase{
		name: "should return empty product type",
		data: AddProductInput{
			"name",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"   ",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}
	invalidSpecialCharactersInProductType := testCase{
		name: "should return invalida characters in name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"ne@**)*w",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}

	// test for brand
	emptyBrandName := testCase{
		name: "should return empty brand name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"",
			"category",
			"subcategory",
		},
		want: false,
	}
	emptyBrandName2 := testCase{
		name: "should return empty brand name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"    ",
			"category",
			"subcategory",
		},
		want: false,
	}
	invalidCharactersInBrand := testCase{
		name: "should return invalid characters in brand",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"br#$^%*&%(&*and",
			"category",
			"subcategory",
		},
		want: false,
	}

	// test for category
	emptyCategoryName := testCase{
		name: "should return empty category name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"",
			"subcategory",
		},
		want: false,
	}
	emptyCategoryName2 := testCase{
		name: "should return invalid empty category name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"    ",
			"subcategory",
		},
		want: false,
	}
	invaliCharactersInCategoryName := testCase{
		name: "should return invalid char  in category",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"cateE%()@_)*&!@#gory",
			"subcategory",
		},
		want: false,
	}

	// test for sub category
	emptySubCategoryName := testCase{
		name: "should return empty category name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"",
		},
		want: false,
	}
	emptySubCategoryName2 := testCase{
		name: "should return invalid empty category name",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"brand",
			"    ",
		},
		want: false,
	}
	invaliCharactersInSubCategoryName := testCase{
		name: "should return invalid char  in category",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"cateE%()@_)*&!@#gory",
			"subc@@$(*%(ategory",
		},
		want: false,
	}

	// valid data
	invalidMainImage := testCase{
		name: "should return true",
		data: AddProductInput{
			"productname",
			"200",
			"is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.",
			"image string",
			[]string{"image", "image2"},
			20,
			"new",
			"brand",
			"category",
			"subcategory",
		},
		want: false,
	}

	cases := []testCase{
		emptyName,
		invalidMainImage,
		emptyName2,
		invalidSpecialCharactersInProductName,
		emptyDescription,
		emptyDescription2,
		invalidSpecialCharactersInDescription,
		emptyImageString,
		emptyImageString2,
		invalidSpecialCharactersInProductType,
		emptyProductType2,
		emptyProductType,
		emptyBrandName,
		emptyBrandName2,
		invalidCharactersInBrand,
		emptyCategoryName,
		emptyCategoryName2,
		invaliCharactersInCategoryName,
		emptySubCategoryName,
		emptySubCategoryName2,
		invaliCharactersInSubCategoryName,
	}

	for _, item := range cases {
		result, _ := ValidateProductInput(&item.data)

		if result != item.want {
			t.Errorf("test failed!!")
		}
	}
	t.Logf("test passed")
}

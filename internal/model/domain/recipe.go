package domain

import "errors"

var (
	ErrTitleParamInvalid  = errors.New("invalid recipe title, must be less than 100 characters")
	ErrRatingParamInvalid = errors.New("invalid rating query parameter, must be an integer and between 0 and 5")
	ErrTagsParamInvalid   = errors.New("invalid tag values, values must be less than 32 characters")
)

const (
	MaxTitleLength = 100
	MaxTagLength   = 32
	MinRating      = 0
	MaxRating      = 5
)

type RecipeAllParams struct {
	Title  string   `query:"title"`
	Rating int      `query:"rating"`
	Tags   []string `query:"tags"`
}

type Recipe struct {
	ID    uint
	UUID  string
	Title string
}

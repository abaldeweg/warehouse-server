package models

// AnalyzeShopSearch struct for analyze shop search data.
type AnalyzeShopSearch struct {
	ID     string `bson:"_id,omitempty" json:"id,omitempty"`
	Branch int    `bson:"branch" json:"branch,omitempty"`
	Date   string `bson:"date" json:"date,omitempty"`
	Term   string `bson:"term" json:"term,omitempty"`
	Page   int    `bson:"page" json:"page,omitempty"`
	Genre  *string `bson:"genre" json:"genre,omitempty"`
}

// AnalyzeShopSearchFilter struct for analyze shop search filter.
type AnalyzeShopSearchFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// AnalyzeShopSearchOrderField struct for analyze shop search order field.
type AnalyzeShopSearchOrderField struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// AnalyzeShopSearchOrderBy struct for analyze shop search order by.
type AnalyzeShopSearchOrderBy struct {
	Article []AnalyzeShopSearchOrderField `json:"article,omitempty"`
}

// AnalyzeShopSearchOptions struct for analyze shop search options.
type AnalyzeShopSearchOptions struct {
	Term    string                    `json:"term,omitempty"`
	Filter  []AnalyzeShopSearchFilter `json:"filter,omitempty"`
	OrderBy AnalyzeShopSearchOrderBy  `json:"orderBy,omitempty"`
	Offset  int                       `json:"offset,omitempty"`
}

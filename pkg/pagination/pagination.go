package pagination

import (
	"math"
	"strings"
)

type PaginationParams struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	SortBy  string `json:"sort_by"`
	SortDir string `json:"sort_dir"`
}

type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
	NextPage   *int        `json:"next_page,omitempty"`
	PrevPage   *int        `json:"prev_page,omitempty"`
}

type PaginationConfig struct {
	DefaultPage    int
	DefaultSize    int
	MaxSize        int
	DefaultSortBy  string
	DefaultSortDir string
}

var DefaultConfig = PaginationConfig{
	DefaultPage:    1,
	DefaultSize:    10,
	MaxSize:        100,
	DefaultSortBy:  "created_at",
	DefaultSortDir: "DESC",
}

func NormalizeParams(params *PaginationParams, config PaginationConfig) {
	if params.Page <= 0 {
		params.Page = config.DefaultPage
	}
	if params.Size <= 0 {
		params.Size = config.DefaultSize
	}
	if params.Size > config.MaxSize {
		params.Size = config.MaxSize
	}
	if params.SortBy == "" {
		params.SortBy = config.DefaultSortBy
	}
	if params.SortDir == "" {
		params.SortDir = config.DefaultSortDir
	}
	params.SortDir = strings.ToUpper(params.SortDir)
	if params.SortDir != "ASC" && params.SortDir != "DESC" {
		params.SortDir = config.DefaultSortDir
	}
}

func CalculateOffset(page, size int) int {
	return (page - 1) * size
}

func CalculateTotalPages(total int64, size int) int {
	if size <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(size)))
}

func HasNextPage(page, totalPages int) bool {
	return page < totalPages
}

func HasPrevPage(page int) bool {
	return page > 1
}

func CreatePaginationResult(data interface{}, total int64, params PaginationParams) PaginationResult {
	totalPages := CalculateTotalPages(total, params.Size)

	result := PaginationResult{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Size:       params.Size,
		TotalPages: totalPages,
		HasNext:    HasNextPage(params.Page, totalPages),
		HasPrev:    HasPrevPage(params.Page),
	}

	if result.HasNext {
		nextPage := params.Page + 1
		result.NextPage = &nextPage
	}
	if result.HasPrev {
		prevPage := params.Page - 1
		result.PrevPage = &prevPage
	}

	return result
}

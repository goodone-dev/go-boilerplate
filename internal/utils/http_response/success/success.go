package success

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse[T any] struct {
	Message    string                `json:"message"`
	Data       T                     `json:"data"`
	Metadata   *PaginationMetadata   `json:"metadata,omitempty"`
	Navigation *PaginationNavigation `json:"navigation,omitempty"`
}

type PaginationMetadata struct {
	Total *int64 `json:"total,omitempty"`
	Pages *int   `json:"pages,omitempty"`
	Page  *int   `json:"page,omitempty"`
	Size  *int   `json:"size,omitempty"`
}

type PaginationNavigation struct {
	First *string `json:"first,omitempty"`
	Next  *string `json:"next,omitempty"`
	Prev  *string `json:"prev,omitempty"`
	Last  *string `json:"last,omitempty"`
}

func Send[T any](c *gin.Context, data T, pagination ...any) {
	var metadata *PaginationMetadata
	var navigation *PaginationNavigation

	for _, v := range pagination {
		switch v := v.(type) {
		case *PaginationMetadata:
			metadata = v
		case *PaginationNavigation:
			navigation = v
		}
	}

	c.JSON(http.StatusOK, SuccessResponse[T]{
		Message:    "success",
		Data:       data,
		Metadata:   metadata,
		Navigation: navigation,
	})
}

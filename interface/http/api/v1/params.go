package v1

const (
	UuidFieldName  = "uuid"
	LoginFieldName = "login"

	PageFieldName  = "page"
	LimitFieldName = "limit"

	CountHeaderName = "X-Pagination-Count"
	PageHeaderName  = "X-Pagination-Page"
	LimitHeaderName = "X-Pagination-Limit"

	LimitDefault = uint(20)
	PageDefault  = uint(1)
)

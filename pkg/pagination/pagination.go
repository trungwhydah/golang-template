package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/golang/be/pkg/logger"
)

const (
	PageSizeLimit  = 100
	DefaultOrderBy = "createdAt"
)

var (
	ErrorInvalidLenCursor = errors.New("invalid length of cursor")
	ErrorEncode           = errors.New("encode cursor error")
)

type OrderDirectionType string

const (
	AscOrderDirection  OrderDirectionType = "asc"
	DescOrderDirection OrderDirectionType = "desc"
)

type DefaultCursor struct {
	CreatedAt *time.Time `json:"createdAt"`
}

type Pagination struct {
	OrderBy        string   `json:"orderBy" form:"orderBy" default:"createdAt" enums:"createdAt,totalDistance,avgPace" url:"orderBy"`
	OrderDirection string   `json:"orderDirection" form:"orderDirection" default:"desc" url:"orderDirection"`
	Limit          int64    `json:"limit" form:"limit" default:"50" url:"limit"`
	Page           int64    `json:"page" form:"page" default:"1" url:"page"`
	Cursors        []string `json:"cursors" form:"cursors" url:"cursors"`
}

type Paging struct {
	Cursors []string `json:"cursors,omitempty"`
	Count   int64    `json:"count,omitempty"`
}

func (p *Pagination) IsAsc() bool {
	return p.OrderDirection == string(AscOrderDirection)
}

func (p *Pagination) BuildSortAndFindOperator() (compareOp string, sortOp int) {
	compareOp = "$lt"
	sortOp = -1

	if p.Limit <= 0 || p.Limit > PageSizeLimit {
		p.Limit = PageSizeLimit
	}

	if p.IsAsc() {
		compareOp = "$gt"
		sortOp = 1
	}

	if p.OrderBy == "" {
		p.OrderBy = DefaultOrderBy
	}

	return compareOp, sortOp
}

func (p *Pagination) DecodeCursor(dest interface{}) error {
	if p.Cursors == nil || len(p.Cursors) == 0 {
		return nil
	}

	if len(p.Cursors) != 1 {
		logger.Errorw(
			"invalid cursors, cursors must have one item",
			"cursors", p.Cursors,
		)

		return ErrorInvalidLenCursor
	}

	bytesData, err := base64.StdEncoding.DecodeString(p.Cursors[0])
	if err != nil {
		logger.Errorw(
			"decode base64 string of cursors err",
			"cursors", p.Cursors,
			"err", err,
		)

		return err
	}

	err = json.Unmarshal(bytesData, dest)
	if err != nil {
		logger.Errorw(
			"json unmarshal error",
			"cursors", p.Cursors,
			"entity", reflect.TypeOf(dest).Name(),
			"err", err,
		)

		return err
	}

	return nil
}

func EncodeCursor(source interface{}) (*string, error) {
	bytesData, err := json.Marshal(source)
	if err != nil {
		logger.Errorw(
			"encode struct err",
			"entity", reflect.TypeOf(source).Name(),
			"err", err,
		)

		return nil, err
	}

	cursorString := base64.StdEncoding.EncodeToString(bytesData)
	if cursorString == "" {
		return nil, ErrorEncode
	}

	return &cursorString, nil
}

func (p *Pagination) Fulfill() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > PageSizeLimit {
		p.Limit = PageSizeLimit
	}

	if p.OrderBy == "" {
		p.OrderBy = DefaultOrderBy
	}

	if p.OrderDirection == "" {
		p.OrderDirection = string(DescOrderDirection)
	}
}

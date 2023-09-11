package mongo

import (
	"errors"
	"reflect"
	"strings"

	"github.com/golang/be/pkg/logger"
	paginationpkg "github.com/golang/be/pkg/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrorInvalidModel = errors.New("required model")
	ErrorModelStruct  = errors.New("model must be struct")
)

type SortOperationType int8

const (
	AscSortMongo  SortOperationType = 1
	DescSortMongo SortOperationType = -1
)

type CompareOperationType string

const (
	LtOperationMongo CompareOperationType = "$lt"
	GtOperationMongo CompareOperationType = "$gt"
)

// Pagination using Cursor
func BuildPaginationPipeline(
	model any,
	pagination *paginationpkg.Pagination,
) (mongo.Pipeline, error) {
	pipeline := mongo.Pipeline{}

	nextCondAgg, err := GetNextCondAggregation(
		model,
		pagination,
	)

	if !errors.Is(err, nil) {
		return nil, err
	}

	if len(nextCondAgg) > 0 {
		pipeline = append(pipeline, nextCondAgg)
	}

	sortAgg, err := GetSortAggregation(
		model,
		pagination.OrderBy,
		pagination.OrderDirection,
	)

	if !errors.Is(err, nil) {
		return nil, err
	}

	if len(sortAgg) > 0 {
		pipeline = append(pipeline, sortAgg)
	}

	if pagination.Limit == 0 || pagination.Limit > paginationpkg.PageSizeLimit {
		pagination.Limit = paginationpkg.PageSizeLimit
	}

	limitAgg := bson.D{{"$limit", pagination.Limit}}
	pipeline = append(pipeline, limitAgg)

	return pipeline, nil
}

func GetSortAggregation(model any, orderBy, orderDirection string) (bson.D, error) {
	if model == nil {
		logger.Errorw("required model to get sort aggregation",
			"model", reflect.TypeOf(model),
		)

		return nil, ErrorInvalidModel
	}

	tags, err := getTagValues(model, true)
	if !errors.Is(err, nil) {
		logger.Errorw("get json tags of model err", "model", reflect.TypeOf(model))

		return nil, ErrorInvalidModel
	}

	sortCond := bson.D{}
	sortOp := DescSortMongo

	if orderDirection == string(paginationpkg.AscOrderDirection) {
		sortOp = AscSortMongo
	}

	for _, tag := range tags {
		if orderBy == "" || orderBy == tag.Tag {
			sortCond = append(sortCond, bson.E{Key: tag.Tag, Value: sortOp})
		}
	}

	if len(sortCond) > 0 {
		return bson.D{{"$sort", sortCond}}, nil
	}

	return bson.D{}, nil
}

func GetNextCondAggregation(model any, pagination *paginationpkg.Pagination) (bson.D, error) {
	if model == nil {
		logger.Errorw("required model to get next cond aggregation",
			"model", reflect.TypeOf(model),
		)

		return nil, ErrorInvalidModel
	}

	if len(pagination.Cursors) == 0 {
		return bson.D{}, nil
	}

	err := pagination.DecodeCursor(model)

	if !errors.Is(err, nil) {
		return nil, err
	}

	tagValues, err := getTagValues(model, false)

	if !errors.Is(err, nil) {
		return nil, err
	}

	compareOp := LtOperationMongo
	if pagination.IsAsc() {
		compareOp = GtOperationMongo
	}

	nextCond := bson.M{}
	listCond := bson.A{}
	lenTagValue := len(tagValues)

	for idx := range tagValues {
		tagValue := tagValues[lenTagValue-idx-1]

		if idx == 0 {
			listCond = append(listCond, bson.M{tagValue.Tag: bson.M{string(compareOp): tagValue.Value}})
		} else {
			listCond = append(listCond, bson.M{tagValue.Tag: tagValue.Value})
			nextCond = bson.M{"$or": bson.A{
				bson.M{tagValue.Tag: bson.M{string(compareOp): tagValue.Value}},
				bson.M{"$and": listCond},
			}}
			listCond = bson.A{nextCond}
		}
	}

	if len(tagValues) == 1 {
		nextCond = bson.M{"$and": listCond}
	}

	return bson.D{{"$match", nextCond}}, nil
}

// BuildPagePaginationPipeline Pagination using Page
func BuildPagePaginationPipeline(pagination *paginationpkg.Pagination) mongo.Pipeline {
	sortOperator := GetSortOperator(pagination)
	sortStage := bson.D{{Key: "$sort", Value: bson.M{pagination.OrderBy: sortOperator}}}

	skipStage := bson.D{{Key: "$skip", Value: (pagination.Page - 1) * pagination.Limit}}

	limitStage := bson.D{{Key: "$limit", Value: pagination.Limit}}

	pipeline := mongo.Pipeline{sortStage, skipStage, limitStage}

	return pipeline
}

func GetSortOperator(pagination *paginationpkg.Pagination) SortOperationType {
	sortOperator := DescSortMongo

	if pagination.IsAsc() {
		sortOperator = AscSortMongo
	}

	return sortOperator
}

type TagValue struct {
	Tag   string
	Value any
}

func getTagValues(model any, ignoreCheckValue bool) ([]TagValue, error) {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		logger.Errorw("model is not struct", "model", reflect.TypeOf(t))

		return nil, ErrorModelStruct
	}

	return getJSONTagValues(t, v, ignoreCheckValue), nil
}

func getJSONTagValues(t reflect.Type, v reflect.Value, ignoreCheckValue bool) (tagValues []TagValue) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.FieldByName(field.Name)

		if field.Type.Kind() == reflect.Struct {
			value := getJSONTagValues(field.Type, fieldValue, ignoreCheckValue)
			tagValues = append(tagValues, value...)

			continue
		}

		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag != "" && (!fieldValue.IsNil() || ignoreCheckValue) {
			tagValues = append(tagValues, TagValue{tag, fieldValue.Interface()})
		}
	}

	return tagValues
}

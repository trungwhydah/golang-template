package httpresp

import "errors"

type ErrorKey string

var (
	ErrorServiceUnavailable   ErrorKey = "EXTERNAL_SERVICE_UNAVAILABLE"
	ErrorInternalServer       ErrorKey = "INTERNAL_SERVER_ERROR"
	ErrorAuthFail             ErrorKey = "AUTH_FAIL"
	ErrorBadRequest           ErrorKey = "BAD_REQUEST"
	ErrorMissingRequiredField ErrorKey = "MISSING_REQUIRED_FIELD"
	ErrorInvalidFieldValue    ErrorKey = "INVALID_FIELD_VALUE"
	ErrorNoPermission         ErrorKey = "NO_PERMISSION"
	ErrorInvalidFieldType     ErrorKey = "INVALID_FIELD_TYPE"
	ErrorInvalidAuthScheme    ErrorKey = "INVALID_AUTH_SCHEME"

	ErrInternalSystem ErrorKey = "error.system.internal"
)

const (
	ErrorMessageServiceUnavailable = "External service is not available"
	ErrorMessageInternalServer     = "internal server error"
	ErrorMessageAuthFail           = "authentication fail"
	ErrorMessageNotFound           = "'{{ .Component }}' not found"
	ErrorMessageMissingField       = "'{{ .Field }}' param is required"
	ErrorMessageMissPermissionClub = "require role is '{{ .Role}}'"
	ErrorMessageAuthScheme         = "invalid authentication scheme. '{{ .Scheme }}' is required"
)

func NewError(key string) error {
	return errors.New(key)
}

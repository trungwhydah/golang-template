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

	ErrInternalSystem                          ErrorKey = "error.system.internal"
	ErrKeySystemInternalServer                          = errors.New("error.system.internal")
	ErrKeyAuthenticationNoPermission                    = errors.New("error.authentication.no_permission")
	ErrKeyAuthenticationInvalidAuthTokenFormat          = errors.New("error.authentication.invalid_auth_token_format")
	ErrKeyAuthenticationInvalidSignature                = errors.New("error.authentication.invalid_signature")
	ErrKeyHTTPValidatorsMissingRequiredField            = errors.New("error.http_validator.missing_required_field")
	ErrKeyHTTPValidatorsInvalidFieldType                = errors.New("error.http_validator.invalid_filed_type")
	ErrKeyHTTPValidatorsDecodeFail                      = errors.New("error.http_validator.decode_fail")
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

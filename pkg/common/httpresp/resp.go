package httpresp

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/golang/be/pkg/common/logger"
	"github.com/golang/be/pkg/common/msgtranslate"
)

type Response struct {
	Data       any      `json:"data,omitempty"`
	Count      int64    `json:"count,omitempty"`
	NextCursor string   `json:"nextCursor,omitempty"`
	ErrorKey   ErrorKey `json:"errorKey,omitempty" example:"error.system.internal"`
	Message    string   `json:"message,omitempty" example:"Internal System Error"`
}

func Error(
	c *gin.Context,
	status int,
	errorKey ErrorKey,
	msg string,
	msgArgs map[string]string,
) {
	msgRes := bytes.Buffer{}
	templ, err := template.New("").Parse(msg)

	if err != nil {
		logger.Errorw(
			"can not parse error template",
			"msg", msg,
			"args", msgArgs,
			"error", err,
		)
	} else {
		err := templ.Execute(&msgRes, msgArgs)
		if err != nil {
			logger.Errorw(
				"can not execute error template",
				"msg", msg,
				"args", msgArgs,
				"error", err,
			)
		}
	}

	c.AbortWithStatusJSON(
		status,
		Response{ErrorKey: errorKey, Message: msgRes.String()},
	)
}

// Success ...
func Success(c *gin.Context, result any, cursors ...string) {
	c.JSON(
		http.StatusOK, Response{
			Data:       result,
			NextCursor: cursors[0],
		},
	)
}

// SuccessNoContent returns success for rest api without content.
func SuccessNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// SuccessCount ...
func SuccessCount(c *gin.Context, result any, count int64, cursors ...string) {
	c.JSON(
		http.StatusOK, Response{
			Data:       result,
			NextCursor: cursors[0],
			Count:      count,
		},
	)
}

func ExternalServerError(g *gin.Context) {
	Error(
		g,
		http.StatusServiceUnavailable,
		ErrorServiceUnavailable,
		ErrorMessageServiceUnavailable,
		nil,
	)
}

func InternalServerError(g *gin.Context) {
	Error(
		g,
		http.StatusInternalServerError,
		ErrorInternalServer,
		ErrorMessageInternalServer,
		nil,
	)
}

func UnauthorizedError(g *gin.Context) {
	Error(
		g,
		http.StatusUnauthorized,
		ErrorAuthFail,
		ErrorMessageAuthFail,
		nil,
	)
}

func NotFoundError(g *gin.Context, comp string) {
	Error(g,
		http.StatusBadRequest,
		ErrorBadRequest,
		ErrorMessageNotFound,
		map[string]string{"Component": comp},
	)
}

func MissingRequiredFieldError(g *gin.Context, field string) {
	Error(
		g,
		http.StatusBadRequest,
		ErrorMissingRequiredField,
		ErrorMessageMissingField,
		map[string]string{"Field": field},
	)
}

func ForbiddenError(g *gin.Context, role string) {
	Error(
		g,
		http.StatusForbidden,
		ErrorNoPermission,
		ErrorMessageMissPermissionClub,
		map[string]string{"Role": role},
	)
}

func InvalidFieldTypeError(g *gin.Context, msgErr string) {
	Error(
		g,
		http.StatusBadRequest,
		ErrorInvalidFieldType,
		msgErr,
		nil,
	)
}

func InvalidFieldValueError(g *gin.Context, msg string) {
	Error(
		g,
		http.StatusBadRequest,
		ErrorInvalidFieldValue,
		msg,
		nil,
	)
}

func ErrorTranslated(
	g *gin.Context,
	msgTranslator *msgtranslate.Translator,
	key string,
	lang *string,
) {
	Error(
		g,
		http.StatusBadRequest,
		ErrorKey(key),
		msgTranslator.Translate(key, lang, nil),
		nil,
	)
}

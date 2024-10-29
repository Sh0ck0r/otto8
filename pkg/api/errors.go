package api

import (
	"errors"
	"net/http"

	"github.com/otto8-ai/otto8/apiclient/types"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func IsHTTPCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if errHttp := (*types.ErrHTTP)(nil); errors.As(err, &errHttp) {
		return errHttp.Code == code
	} else if errMeta := (*apierrors.StatusError)(nil); errors.As(err, &errMeta) {
		return errMeta.ErrStatus.Code == int32(code)
	}
	return false
}

func IsConflict(err error) bool {
	return IsHTTPCode(err, http.StatusConflict)
}

var ErrMustAuth = &types.ErrHTTP{
	Code:    http.StatusUnauthorized,
	Message: "unauthorized request, must authenticate",
}

package validator

import (
	"github.com/kataras/iris/v12"
)

// ReadAndValidate reads the request body into v and validates it.
// Uses iris's built-in JSON decoding.
func ReadAndValidate(ctx iris.Context, v interface{}) error {
	return ctx.ReadJSON(v)
}

package sample

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Controller handles HTTP requests for sample resources.
type Controller struct {
	SampleSvc SampleService
}

func (c *Controller) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "ListSamples")
	b.Handle("POST", "/", "CreateSample")
	b.Handle("GET", "/{id:int32}", "GetSample")
	b.Handle("PUT", "/{id:int32}", "UpdateSample")
	b.Handle("PATCH", "/{id:int32}", "PatchSample")
	b.Handle("DELETE", "/{id:int32}", "DeleteSample")
}

// ListSamples godoc
// @Summary      List all samples
// @Description  Returns a paginated list of samples
// @Tags         samples
// @Accept       json
// @Produce      json
// @Param        page   query  int  false  "Page number (default: 1)"
// @Param        limit  query  int  false  "Items per page (default: 10)"
// @Success      200  {object}  ListSamplesResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/samples [get]
func (c *Controller) ListSamples(ctx iris.Context) (mvc.Result, error) {
	var req ListSamplesRequest
	if err := ctx.ReadQuery(&req); err != nil {
		return nil, err
	}

	result, err := c.SampleSvc.ListSamples(ctx.Request().Context(), &req)
	if err != nil {
		return nil, err
	}

	return mvc.Response{Object: result}, nil
}

// CreateSample godoc
// @Summary      Create a sample
// @Description  Creates a new sample resource
// @Tags         samples
// @Accept       json
// @Produce      json
// @Param        request  body  CreateSampleRequest  true  "Sample payload"
// @Success      201  {object}  Sample
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/samples [post]
func (c *Controller) CreateSample(ctx iris.Context) (mvc.Result, error) {
	var req CreateSampleRequest
	if err := ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	// TODO: extract user details from auth middleware context
	// details := ctx.Values().Get("userDetails").(*domain.UserDetails)
	// req.CreatedBy = details.UserID

	result, err := c.SampleSvc.CreateSample(ctx.Request().Context(), &req)
	if err != nil {
		return nil, err
	}

	return mvc.Response{Code: iris.StatusCreated, Object: result}, nil
}

// GetSample godoc
// @Summary      Get a sample
// @Description  Returns a single sample by ID
// @Tags         samples
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Sample ID"
// @Success      200  {object}  Sample
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/samples/{id} [get]
func (c *Controller) GetSample(ctx iris.Context, id int32) (mvc.Result, error) {
	result, err := c.SampleSvc.GetSample(ctx.Request().Context(), id)
	if err != nil {
		return nil, err
	}

	return mvc.Response{Object: result}, nil
}

// UpdateSample godoc
// @Summary      Update a sample
// @Description  Replaces all fields of a sample by ID
// @Tags         samples
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "Sample ID"
// @Param        request  body  UpdateSampleRequest  true  "Sample payload"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/samples/{id} [put]
func (c *Controller) UpdateSample(ctx iris.Context, id int32) (mvc.Result, error) {
	var req UpdateSampleRequest
	if err := ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	if err := c.SampleSvc.UpdateSample(ctx.Request().Context(), id, &req); err != nil {
		return nil, err
	}

	return mvc.Response{Code: iris.StatusNoContent}, nil
}

// PatchSample godoc
// @Summary      Partially update a sample
// @Description  Updates only provided fields of a sample by ID
// @Tags         samples
// @Accept       json
// @Produce      json
// @Param        id       path  int                 true  "Sample ID"
// @Param        request  body  PatchSampleRequest  true  "Sample patch payload"
// @Success      200  {object}  Sample
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/samples/{id} [patch]
func (c *Controller) PatchSample(ctx iris.Context, id int32) (mvc.Result, error) {
	var req PatchSampleRequest
	if err := ctx.ReadJSON(&req); err != nil {
		return nil, err
	}

	result, err := c.SampleSvc.PatchSample(ctx.Request().Context(), id, &req)
	if err != nil {
		return nil, err
	}

	return mvc.Response{Object: result}, nil
}

// DeleteSample godoc
// @Summary      Delete a sample
// @Description  Deletes a sample by ID
// @Tags         samples
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Sample ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/samples/{id} [delete]
func (c *Controller) DeleteSample(ctx iris.Context, id int32) (mvc.Result, error) {
	if err := c.SampleSvc.DeleteSample(ctx.Request().Context(), id); err != nil {
		return nil, err
	}

	return mvc.Response{Code: iris.StatusNoContent}, nil
}

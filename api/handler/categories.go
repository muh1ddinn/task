package handler

import (
	"context"
	"fmt"
	"net/http"
	"task/api/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Createcategories godoc
// @Security ApiKeyAuth
// @Router      /categories [POST]
// @Summary     Create a categories
// @Description Create a new categories
// @Tags        categories
// @Accept      json
// @Produce     json
// @Param       categories body model.Categories true "categories"
// @Success     200 {object} model.Categories
// @Failure     400 {object} model.Response
// @Failure     404 {object} model.Response
// @Failure     500 {object} model.Response
func (h Handler) Createcat(c *gin.Context) {
	var categories model.Categories

	if err := c.ShouldBindJSON(&categories); err != nil {
		handleResponseLog(c, h.Log, "Invalid request", http.StatusBadRequest, err.Error())
		return
	}

	existingname, err := h.Services.Categoriess().Checkename(context.Background(), categories.Name)
	if err != nil {
		handleResponseLog(c, h.Log, "Error while checking name", http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(existingname, "gghj")
	if existingname != "" {
		handleResponseLog(c, h.Log, "This name is already used", http.StatusBadRequest, "name already exists")
		return
	}

	createdcategories, err := h.Services.Categoriess().Create(context.Background(), categories)
	if err != nil {
		handleResponseLog(c, h.Log, "Error while creating categories", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Created successfully", http.StatusOK, createdcategories)
}

// GetAllcategories godoc
// @Security ApiKeyAuth
// @Router 			/categories [GET]
// @Summary 		Get all categories
// @Description		Retrieves information about all categories.
// @Tags 			categories
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "categories"
// @Param 			page query uint64 false "page"
// @Param 			limit query uint64 false "limit"
// @Success 		200 {object} model.GetAllcategoriesResponse
// @Failure 		400 {object} model.Response
// @Failure 		500 {object} model.Response
func (h Handler) Getallcat(c *gin.Context) {

	var request = model.GetAllCategoriestRequest{}

	request.Search = c.Query("search")

	page, err := ParsePageQueryParam(c)

	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing page", http.StatusBadRequest, err.Error())

	}

	limit, err := ParseLimitQueryParam(c)
	if err != nil {
		handleResponseLog(c, h.Log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("page: ", page)
	fmt.Println("limit: ", limit)

	request.Page = page
	request.Limit = limit

	categories, err := h.Services.Categoriess().GetAll(context.Background(), request)
	if err != nil {

		handleResponseLog(c, h.Log, "error while getting user", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "", http.StatusOK, categories)
}

// GetById godoc
// @Security ApiKeyAuth
// @Router		/categories/{id} [GET]
// @Summary		get a categories by its id
// @Description This api gets a categories by its id and returns its info
// @Tags		categories
// @Accept		json
// @Produce		json
// @Param		id path string true "categories"
// @Success		200  {object}  model.Getcategoriest
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) GetByIDcat(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleResponseLog(c, h.Log, "missing car ID", http.StatusBadRequest, id)
		return
	}

	customer, err := h.Services.Categoriess().GetByID(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting user by ID", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "user was successfully gotten by Id", http.StatusOK, customer)
}

// Delete godoc
// @Security ApiKeyAuth
// @Router		/categories_s/{id} [DELETE]
// @Summary		deletesoft a user by its id
// @Description This api deletes a categories by its id and returns error or nil
// @Tags		categories
// @Accept		json
// @Produce		json
// @Param		id path string true "categories ID"
// @Success		200  {object}  string
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) Deletesoftcat(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id ", http.StatusBadRequest, err.Error())
		return
	}

	id, err = h.Services.Categoriess().Deletesoft(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while deleting user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "all good ", http.StatusOK, id)

}

// Delete godoc
// @Security ApiKeyAuth
// @Router		/categories/{id} [DELETE]
// @Summary		deletehard a user by its id
// @Description This api deletes a categories by its id and returns error or nil
// @Tags		categories
// @Accept		json
// @Produce		json
// @Param		id path string true "categories ID"
// @Success		200  {object}  string
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) Deletehardcat(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id ", http.StatusBadRequest, err.Error())
		return
	}

	id, err = h.Services.Categoriess().Delet(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while deleting user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "all good ", http.StatusOK, id)

}

// Update godoc
// @Security ApiKeyAuth
// @Router      /categories [PATCH]
// @Summary     Update a categories
// @Description Update a new categories
// @Tags        categories
// @Accept      json
// @Produce 	json
// @Param 		categories body model.Patchcategories true "categories"
// @Success 	200  {object}  model.Getcategoriest
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) Patchcat(c *gin.Context) {
	cus := model.Patchcategories{}

	if err := c.ShouldBindJSON(&cus); err != nil {
		return
	}
	fmt.Println(cus)

	existingname, err := h.Services.Categoriess().Checkename(context.Background(), cus.Name)

	if err != nil {
		handleResponseLog(c, h.Log, "Error while checking name", http.StatusInternalServerError, err.Error())
		return
	}

	if existingname != "" {
		handleResponseLog(c, h.Log, "This email is already used", http.StatusBadRequest, "name already exists")
		return
	}

	id, err := h.Services.Categoriess().Patchcat(context.Background(), cus)
	if err != nil {
		handleResponseLog(c, h.Log, "error while creating user", http.StatusBadRequest, err.Error())
		return
	}
	handleResponseLog(c, h.Log, "Created successfully", http.StatusOK, id)
}

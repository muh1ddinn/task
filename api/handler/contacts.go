package handler

import (
	"context"
	"fmt"
	"net/http"
	"task/api/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Createcontact godoc
// @Security ApiKeyAuth
// @Router      /contact [POST]
// @Summary     Create a contact
// @Description Create a new contact
// @Tags        contact
// @Accept      json
// @Produce     json
// @Param       contact body model.Contact true "contact"
// @Success     200 {object} model.Contact
// @Failure     400 {object} model.Response
// @Failure     404 {object} model.Response
// @Failure     500 {object} model.Response
func (h Handler) Create(c *gin.Context) {
	var contact model.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		handleResponseLog(c, h.Log, "Invalid request", http.StatusBadRequest, err.Error())
		return
	}

	existingEmail, err := h.Services.Contacts().Checkemail(context.Background(), contact.Email)
	if err != nil {
		handleResponseLog(c, h.Log, "Error while checking email", http.StatusInternalServerError, err.Error())
		return
	}

	if existingEmail != "" {
		handleResponseLog(c, h.Log, "This email is already used", http.StatusBadRequest, "email already exists ")
		return
	}

	createdContact, err := h.Services.Contacts().Create(context.Background(), contact)
	if err != nil {
		handleResponseLog(c, h.Log, "Error while creating contact", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "Created successfully", http.StatusOK, createdContact)
}

// GetAllcontact godoc
// @Security ApiKeyAuth
// @Router 			/contact [GET]
// @Summary 		Get all contact
// @Description		Retrieves information about all contact.
// @Tags 			contact
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "contact"
// @Param 			page query uint64 false "page"
// @Param 			limit query uint64 false "limit"
// @Success 		200 {object} model.GetAllContactResponse
// @Failure 		400 {object} model.Response
// @Failure 		500 {object} model.Response
func (h Handler) Getall(c *gin.Context) {

	var request = model.GetAllContactRequest{}

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

	contact, err := h.Services.Contacts().GetAll(context.Background(), request)
	if err != nil {

		handleResponseLog(c, h.Log, "error while getting user", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "", http.StatusOK, contact)
}

// GetById godoc
// @Security ApiKeyAuth
// @Router		/contact/{id} [GET]
// @Summary		get a contact by its id
// @Description This api gets a contact by its id and returns its info
// @Tags		contact
// @Accept		json
// @Produce		json
// @Param		id path string true "contact"
// @Success		200  {object}  model.GetAllContact
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		handleResponseLog(c, h.Log, "missing car ID", http.StatusBadRequest, id)
		return
	}

	customer, err := h.Services.Contacts().GetByID(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting user by ID", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "user was successfully gotten by Id", http.StatusOK, customer)
}

// Delete godoc
// @Security ApiKeyAuth
// @Router		/contact_s/{id} [DELETE]
// @Summary		deletesoft a user by its id
// @Description This api deletes a contact by its id and returns error or nil
// @Tags		contact
// @Accept		json
// @Produce		json
// @Param		id path string true "contact ID"
// @Success		200  {object}  string
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) Deletesoft(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id ", http.StatusBadRequest, err.Error())
		return
	}

	id, err = h.Services.Contacts().Deletesoft(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while deleting user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "all good ", http.StatusOK, id)

}

// Delete godoc
// @Security ApiKeyAuth
// @Router		/contact/{id} [DELETE]
// @Summary		deletehard a user by its id
// @Description This api deletes a contact by its id and returns error or nil
// @Tags		contact
// @Accept		json
// @Produce		json
// @Param		id path string true "contact ID"
// @Success		200  {object}  string
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) Deletehard(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while validating id ", http.StatusBadRequest, err.Error())
		return
	}

	id, err = h.Services.Contacts().Delet(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while deleting user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "all good ", http.StatusOK, id)

}

// Update godoc
// @Security ApiKeyAuth
// @Router      /contact [PATCH]
// @Summary     Update a contact
// @Description Update a new contact
// @Tags        contact
// @Accept      json
// @Produce 	json
// @Param 		contact body model.PatchContact false "contact"
// @Success 	200  {object}  model.GetAllContact
// @Failure		400  {object}  model.Response
// @Failure		404  {object}  model.Response
// @Failure		500  {object}  model.Response
func (h Handler) Patch(c *gin.Context) {
	cus := model.PatchContact{}

	if err := c.ShouldBindJSON(&cus); err != nil {
		return
	}

	existingEmail, err := h.Services.Contacts().Checkemail(context.Background(), cus.Email)
	if err != nil {
		handleResponseLog(c, h.Log, "Error while checking email", http.StatusInternalServerError, err.Error())
		return
	}

	if existingEmail != "" {
		handleResponseLog(c, h.Log, "This email is already used", http.StatusBadRequest, "email already exists ")
		return
	}
	id, err := h.Services.Contacts().Patch(context.Background(), cus)
	if err != nil {
		handleResponseLog(c, h.Log, "error while creating user", http.StatusBadRequest, err.Error())
		return
	}
	handleResponseLog(c, h.Log, "Created successfully", http.StatusOK, id)
}

// GetById godoc
// @Security ApiKeyAuth
// @Router      /contact/history/{id} [GET]
// @Summary     get a contact_history by its id
// @Description This api gets a contact by its id and returns its info
// @Tags        contact
// @Accept      json
// @Produce     json
// @Param       id path string true "contact"
// @Success     200  {object}    model.ContactHistory
// @Failure     400  {object}    model.Response
// @Failure     404  {object}    model.Response
// @Failure     500  {object}    model.Response
func (h Handler) History(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		handleResponseLog(c, h.Log, "missing ID", http.StatusBadRequest, nil)
		return
	}

	fmt.Println(id)
	history, err := h.Services.Contacts().History(context.Background(), id)
	if err != nil {
		handleResponseLog(c, h.Log, "error while getting history by ID", http.StatusBadRequest, err.Error())
		return
	}

	handleResponseLog(c, h.Log, "history was successfully retrieved", http.StatusOK, history)
}

package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-lock/models"
	"strings"
)

func (c Controller) CreateEntity(ctx *gin.Context) {

	var entity models.Entity

	// decoding json message to user model
	err := json.NewDecoder(strings.NewReader(ctx.PostForm("data"))).Decode(&entity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if entity.Password == "" {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  "password not provided for given entity",
		})
		return
	}

	if entity.SecretKey == "" {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  "secret key not provided",
		})
		return
	}

	if entity.Type < 6 && entity.Type > 0 {
		entity.IconPath = c.service.GetEntityIconPath(entity.Type)
	}
	entity.Password = c.service.EncryptPassword(entity.SecretKey, entity.Password)
	entity.UserUuid = c.service.Me(ctx)

	createdEntity, err := c.service.CreateEntity(ctx, entity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if entity.Type == 6 {
		path, err := c.service.UploadIconToBucket(ctx, createdEntity.Uuid)
		if err != nil {
			c.SendResponse(ctx, Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			})
			return
		}

		createdEntity.IconPath = path
		err = c.service.UpdateEntity(ctx, createdEntity)
		if err != nil {
			c.SendResponse(ctx, Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			})
			return
		}
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "entity successfully created",
	})
	return
}

func (c Controller) UpdateEntity(ctx *gin.Context) {
	var updatedEntity *models.Entity

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&updatedEntity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if updatedEntity.Password != "" {
		if updatedEntity.SecretKey == "" {
			c.SendResponse(ctx, Response{
				Status: http.StatusBadRequest,
				Error:  "password can not be updated without providing secret key",
			})
			return
		}
		updatedEntity.Password = c.service.EncryptPassword(updatedEntity.SecretKey, updatedEntity.Password)
	}

	if updatedEntity.Type != 0 {
		updatedEntity.IconPath = c.service.GetEntityIconPath(updatedEntity.Type)
	}

	err = c.service.UpdateEntity(ctx, updatedEntity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "entity successfully updated",
	})
	return

}

func (c Controller) DeleteEntity(ctx *gin.Context) {
	var request struct {
		LoginPassword string `json:"password"`
	}

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&request)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	// checking if user has permission to delete an entity
	err = c.service.Authorize(c.service.Me(ctx), request.LoginPassword)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		})
		return
	}

	entityUuid := ctx.Param("entity_uuid")

	err = c.service.DeleteEntity(ctx, entityUuid)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "entity successfully deleted",
	})
	return

}

func (c Controller) GetEntity(ctx *gin.Context) {
	var request struct {
		SecretKey string `json:"secret_key"`
	}

	// decoding json message to user model
	err := json.NewDecoder(ctx.Request.Body).Decode(&request)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	entity, err := c.service.GetEntityByUuid(ctx, ctx.Param("entity_uuid"), request.SecretKey)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(200, entity)

}

func (c Controller) ListEntities(ctx *gin.Context) {
	entities, err := c.service.ListEntities(ctx)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(200, entities)
}

func (c Controller) DownloadEntityIcon(ctx *gin.Context) {

	entityUuid := ctx.Param("entity_uuid")

	_, err := c.service.DownloadEntityIcon(ctx, entityUuid)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}
}

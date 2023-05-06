package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/UitsHabib/ecommerce-microservice/platform/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getFeaturesRequest struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit" binding:"required,min=1,max=10"`
}

func (server *Server) getFeatures(ctx *gin.Context) {
	var req getFeaturesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListFeaturesParams{
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}
	features, err := server.store.ListFeatures(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tmp_features := FormatFeatures(&features)
	ctx.JSON(http.StatusOK, tmp_features)
}

type getFeatureRequest struct {
	ID string `uri:"id" binding:"required"`
}

type getFeatureResponse struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	CreatedBy   uuid.NullUUID `json:"created_by"`
	UpdatedBy   uuid.NullUUID `json:"updated_by"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (server *Server) getFeature(ctx *gin.Context) {
	var req getFeatureRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	feature, err := server.store.GetFeature(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	f := formatFeature(&feature)
	ctx.JSON(http.StatusOK, f)
}

func FormatFeatures(features *[]db.Feature) []getFeatureResponse {
	tmp_features := []getFeatureResponse{}

	for _, row := range *features {
		f := formatFeature(&row)

		tmp_features = append(tmp_features, f)
	}

	return tmp_features
}

func formatFeature(feature *db.Feature) getFeatureResponse {
	f := &getFeatureResponse{
		ID:          feature.ID,
		Title:       feature.Title,
		Slug:        feature.Slug,
		Description: feature.Description.String,
		CreatedBy:   feature.CreatedBy,
		UpdatedBy:   feature.UpdatedBy,
		CreatedAt:   feature.CreatedAt,
		UpdatedAt:   feature.UpdatedAt,
	}
	return *f
}

package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heloayer/check-order-status/internal/resources/http/handler/dto"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (srv *Server) GetOrder(c *gin.Context) {
	log.Info().Msg("starting get order request")

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	order, err := srv.service.Order.GetOrderByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("error getting order by ID")
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Description: err.Error()})
		return
	}

	log.Info().Msg("success getting order")
	c.JSON(http.StatusOK, order)
}

func (srv *Server) CreateOrder(c *gin.Context) {
	log.Info().Msg("starting create order request")

	var request dto.OrderRequest
	if err := c.BindJSON(&request); err != nil {
		log.Error().Err(err).Msg("error binding json request")
		c.Set(ErrBindingJsonError, err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Description: err.Error()})
		return
	}
	/* provider query either glovo or wolt: glovo > Json Server, wolt > Mock API */
	provider, ok := c.GetQuery("provider")
	if !ok {
		log.Error().Msg("missing provider query parameter")
		c.Set(ErrMissingProvider, fmt.Errorf("missing provider"))
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "missing provider, choose between 'glovo' or 'wolt",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := srv.service.Order.CreateOrder(ctx, request, provider)
	if err != nil {
		log.Error().Err(err).Msg("error creating order")
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Description: err.Error()})
		return
	}

	log.Info().Msg("order created successfully")
	c.Status(http.StatusOK)
}

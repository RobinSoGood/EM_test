package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/RobinSoGood/EM_test/internal/logger"
	"github.com/RobinSoGood/EM_test/internal/models"
	"github.com/RobinSoGood/EM_test/internal/storage/storageerror"

	"github.com/gin-gonic/gin"
)

// @Summary Add a new subscription
// @Description Add a new subscription to the system
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.SubRequest true "Subscription details"
// @Security ApiKeyAuth
// @Success 201 {string} string "Sub {id} was saved"
// @Failure 400 {object} object "Bad request"
// @Failure 409 {object} object "Subscription already exists"
// @Failure 500 {object} object "Internal server error"
// @Router /subs [post]

func (s *Server) addSubHandler(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("ID")
	if !exist {
		log.Error().Msg("ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ID not found"})
		return
	}
	var subReq models.SubRequest
	err := ctx.ShouldBindBodyWithJSON(&subReq)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub := models.Sub{
		ID:          subReq.ID,
		UserID:      subReq.UserID,
		ServiceName: subReq.ServiceName,
		Price:       subReq.Price,
		StartDate:   subReq.StartDate,
		EndDate:     subReq.EndDate,
	}
	sid, err := s.sService.AddSub(sub)
	if err != nil {
		log.Error().Err(err).Msg("save sub failed")
		if errors.Is(err, storageerror.ErrSubAlredyExist) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusCreated, "Sub %s was saved", sid)
}

// @Summary Get all subscriptions
// @Description Retrieve all subscriptions from the system
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Sub "List of subscriptions"
// @Failure 204 {object} object "No content"
// @Failure 500 {object} object "Internal server error"
// @Router /subs [get]

func (s *Server) getSubsHandler(ctx *gin.Context) {
	log := logger.Get()
	subs, err := s.sService.GetSubs()
	if err != nil {
		log.Error().Err(err).Msg("get all subs form storage failed")
		if errors.Is(err, errors.New("empyt storage")) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, subs)
}
func (s *Server) getSubHandler(ctx *gin.Context) {
	log := logger.Get()
	sid := ctx.Param("id")
	log.Debug().Str("id", sid).Msg("check id from param")
	sub, err := s.sService.GetSub(sid)
	if err != nil {
		log.Error().Err(err).Msg("get all subs form storage failed")
		if errors.Is(err, storageerror.ErrSubNoFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reqSub := models.SubRequest{
		ID:          sub.ID,
		UserID:      sub.UserID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}
	ctx.JSON(http.StatusOK, reqSub)
}

// @Summary Delete subscription
// @Description Mark subscription as deleted (soft delete)
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {string} string "Sub {id} was deleted"
// @Failure 204 {object} object "Subscription not found"
// @Failure 500 {object} object "Internal server error"
// @Router /subs/{id} [delete]
func (s *Server) deleteSubHandler(ctx *gin.Context) {
	log := logger.Get()
	sid := ctx.Param("id")
	err := s.sService.SetDeleteStatus(sid)
	if err != nil {
		log.Error().Err(err).Msg("delete sub failed")
		if errors.Is(err, storageerror.ErrSubNoFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "Sub %s was deleted", sid)
}

// @Summary Calculate total price
// @Description Calculate total price for subscriptions in given period
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body models.SumRequest true "Period and filter criteria"
// @Success 200 {object} object "Total price calculation result"
// @Failure 400 {object} object "Invalid date format or period"
// @Failure 500 {object} object "Internal server error"
// @Router /subs/total [post]
func (s *Server) getTotalPriceHandler(ctx *gin.Context) {
	log := logger.Get()

	var req models.SumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	start, err := time.Parse("2006-01-02", req.Start)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format, use YYYY-MM-DD"})
		return
	}

	end, err := time.Parse("2006-01-02", req.End)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format, use YYYY-MM-DD"})
		return
	}

	if end.Before(start) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "end date must be after start date"})
		return
	}

	total, err := s.sService.GetTotalPriceByPeriod(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to calculate total price")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total_price": total,
		"period": gin.H{
			"start": start.Format("2006-01-02"),
			"end":   end.Format("2006-01-02"),
		},
		"user_id":      req.UserID,
		"service_name": req.ServiceName,
	})
}

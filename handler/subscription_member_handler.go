package handler

import (
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/jinzhu/copier"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/appcontext"
)

type SubscriptionMemberHandler struct {
	Db            *gorm.DB
	logger        *zap.Logger
	uuidGenerator models.UUIDGenerator
}

func NewSubscriptionMemberHandler(db *gorm.DB, logger *zap.Logger, uuidGenerator models.UUIDGenerator,
) *SubscriptionMemberHandler {
	return &SubscriptionMemberHandler{
		Db:            db,
		logger:        logger,
		uuidGenerator: uuidGenerator,
	}
}

type searchSubscriptionMemberParams struct {
	UserID                  string `form:"user_id" binding:"omitempty,uuid4"`
	MemberStatus            string `form:"member_status" binding:"omitempty,oneof=active inactive"`
	MemberStartDateFrom     string `form:"member_start_date_from" binding:"omitempty,datetime"`
	MemberEndDateTo         string `form:"member_end_date_to" binding:"omitempty,datetime"`
	MemberStopStartDateFrom string `form:"member_stop_start_date_from" binding:"omitempty,datetime"`
	MemberStopEndDateTo     string `form:"member_stop_end_date_to" binding:"omitempty,datetime"`
	Offset                  string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit                   string `form:"limit,default=10" binding:"omitempty,numeric"`
}

type subscriptionMemberResponseItem struct {
	ID                  string     `json:"id" example:"218c51c0-904e-4743-a2ae-94f0e34a0d6f"`
	UserID              string     `json:"user_id" example:"218c51c0-904e-4743-a2ae-94f0e34a0d6f"`
	MemberStatus        string     `json:"member_status" example:"basic"`
	MemberStartDate     time.Time  `json:"member_start_date" example:"2020-01-01T00:00:00+09:00"`
	MemberEndDate       *time.Time `json:"member_end_date" example:"2020-01-01T00:00:00+09:00"`
	MemberStopStartDate time.Time  `json:"member_stop_start_date" example:"2020-01-01T00:00:00+09:00"`
	MemberStopEndDate   *time.Time `json:"member_stop_end_date" example:"2020-01-01T00:00:00+09:00"`
}

type subscriptionMembersResponse struct {
	Total              int                              `json:"total"`
	SubscriptionMember []subscriptionMemberResponseItem `json:"subscription_members"`
}

// GetSubscriptionMember @title 一覧取得
// @id GetSubscriptionMember
// @tags subscription_members
// @version バージョン(1.0)
// @description 指定された条件に一致するproject一覧情報を取得する
// @Summary subscription_member一覧取得
// @Produce json
// @Success 200 {object} subscriptionMembersResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /subscriptionMembers [GET]
// @Param user_id query string false "ユーザID" minlength(36) maxlength(36) format(UUID v4)
// @Param member_status query string false "会員ステータス" maxlength(64)
// @Param member_start_date_from query string false "会員開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param member_end_date_to query string false "会員開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param member_stop_start_date_from query string false "会員停止開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param member_stop_end_date_to query string false "会員停止終了日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *SubscriptionMemberHandler) GetSubscriptionMember(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchSubscriptionMemberParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	var subscriptionMembers []models.SubscriptionMember
	query := createSubscriptionMemberQueryBuilder(params, h)
	query.Find(&subscriptionMembers)

	res := subscriptionMembersResponse{
		Total:              len(subscriptionMembers),
		SubscriptionMember: []subscriptionMemberResponseItem{},
	}
	res.Total = len(subscriptionMembers)
	if err := copier.Copy(&res.SubscriptionMember, &subscriptionMembers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetSubscriptionMemberDetail @title subscription_member詳細
// @id GetSubscriptionMemberDetail
// @tags subscription_members
// @version バージョン(1.0)
// @description subscription_member詳細を返す
// @Summary subscription_member詳細取得
// @Produce json
// @Success 200 {object} subscriptionMemberResponseItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /subscription_members/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *SubscriptionMemberHandler) GetSubscriptionMemberDetail(ctx *gin.Context) {
	var subscriptionMember models.SubscriptionMember
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).Find(&subscriptionMember).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("subscription_member not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}

	ctx.JSON(http.StatusOK, subscriptionMember)
}

func createSubscriptionMemberQueryBuilder(param searchSubscriptionMemberParams, h *SubscriptionMemberHandler) *gorm.DB {
	var subscriptionMember []models.SubscriptionMember
	query := h.Db.Find(&subscriptionMember)
	if param.UserID != "" {
		query = query.Where("user_id = ?", param.UserID)
	}
	if param.MemberStatus != "" {
		query = query.Where("member_status = ?", param.MemberStatus)
	}
	if param.MemberStartDateFrom != "" {
		query = query.Where("member_start_date >= ?", param.MemberStartDateFrom)
	}
	if param.MemberEndDateTo != "" {
		query = query.Where("member_start_date <= ?", param.MemberEndDateTo)
	}
	if param.MemberStopStartDateFrom != "" {
		query = query.Where("member_stop_start_date >= ?", param.MemberStopStartDateFrom)
	}
	if param.MemberStopEndDateTo != "" {
		query = query.Where("member_stop_start_date <= ?", param.MemberStopEndDateTo)
	}
	if param.Offset != "" {
		query = query.Offset(param.Offset)
	}
	if param.Limit != "" {
		query = query.Limit(param.Limit)
	}
	return query
}

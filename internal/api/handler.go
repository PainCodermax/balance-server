package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paincodermax/balance-service/internal/expense"
)

type Handler struct {
	service *expense.Service
}

func NewHandler(service *expense.Service) *Handler {
	return &Handler{service: service}
}

// CreateExpense handles the request to add a new expense.
func (h *Handler) CreateExpense(c *gin.Context) {
	var req expense.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid request data",
			StatusCode: http.StatusBadRequest})
		return
	}

	result, err := h.service.CreateExpense(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid request data",
			StatusCode: http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusCreated, response{
		Data:       result,
		StatusCode: http.StatusCreated,
	})
}

func (h *Handler) GetExpensesByMonthYear(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")
	page := c.Query("page")
	pageSize := c.Query("page_size")

	if year == "" || month == "" {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Year and month are required",
			StatusCode: http.StatusBadRequest})
		return
	}

	// Convert query parameters to integers
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid year format",
			StatusCode: http.StatusBadRequest})
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid month format",
			StatusCode: http.StatusBadRequest})
		return
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	expenses, err := h.service.GetExpensesByMonthYear(c.Request.Context(), yearInt, time.Month(monthInt), pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Message:    "Failed to retrieve expenses",
			StatusCode: http.StatusInternalServerError})
		return
	}
	total, err := h.service.GetTotalByMonthYear(c.Request.Context(), yearInt, time.Month(monthInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Message:    "Failed to retrieve total expenses",
			StatusCode: http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, response{
		Data: expenses,
		Pagination: map[string]int{
			"page":      pageInt,
			"page_size": pageSizeInt,
			"total":     total,
		},
		StatusCode: http.StatusOK,
	})
}

func (h *Handler) GetTotalByPayer(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")

	if year == "" || month == "" {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Year and month are required",
			StatusCode: http.StatusBadRequest})
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid year format",
			StatusCode: http.StatusBadRequest})
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid month format",
			StatusCode: http.StatusBadRequest})
		return
	}

	payers := []string{"Trung", "Thang"}
	totals := make(map[string]int)

	for _, payer := range payers {
		total, err := h.service.GetTotalByPayerMonthYear(c.Request.Context(), payer, yearInt, time.Month(monthInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Message:    "Failed to retrieve total for " + payer,
				StatusCode: http.StatusInternalServerError})
			return
		}
		totals[payer] = total
	}

	c.JSON(http.StatusOK, response{
		Data:       totals,
		StatusCode: http.StatusOK,
	})
}

func (h *Handler) GetOutstandingBalanceByMonthYear(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")

	if year == "" || month == "" {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Year and month are required",
			StatusCode: http.StatusBadRequest})
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid year format",
			StatusCode: http.StatusBadRequest})
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, errorResponse{
			Message:    "Invalid month format",
			StatusCode: http.StatusBadRequest})
		return
	}

	trungTotal, err := h.service.GetTotalByPayerMonthYear(c.Request.Context(), "Trung", yearInt, time.Month(monthInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Message:    "Failed to retrieve total for Trung",
			StatusCode: http.StatusInternalServerError})
		return
	}

	thangTotal, err := h.service.GetTotalByPayerMonthYear(c.Request.Context(), "Thang", yearInt, time.Month(monthInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Message:    "Failed to retrieve total for Thang",
			StatusCode: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, response{
		Data: map[string]int{
			"Trung": (trungTotal / 2) - (thangTotal / 2),
			"Thang": (thangTotal / 2) - (trungTotal / 2),
		},
		StatusCode: http.StatusOK,
	})
}

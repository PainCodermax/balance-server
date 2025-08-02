package expense

type CreateExpenseRequest struct {
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Payer       string  `json:"payer" binding:"required,oneof=Trung Thang"`
}

// BalanceResult định nghĩa cấu trúc cho kết quả trả về của API tính toán.
type BalanceResult struct {
	TotalFund       float64 `json:"total_fund"`
	TotalSpentByA   float64 `json:"total_spent_by_A"`
	TotalSpentByB   float64 `json:"total_spent_by_B"`
	EachPersonShare float64 `json:"each_person_share"`
	Debt            *Debt   `json:"debt"`
}

// Debt định nghĩa cấu trúc cho thông tin công nợ.
type Debt struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

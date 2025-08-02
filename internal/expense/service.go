package expense

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	collection *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		collection: db.Collection("expenses"),
	}
}

// CreateExpense tạo một bản ghi chi tiêu mới.
func (s *Service) CreateExpense(ctx context.Context, req CreateExpenseRequest) (*Expense, error) {
	expense := &Expense{
		Description: req.Description,
		Amount:      req.Amount,
		Payer:       req.Payer,
		CreatedAt:   time.Now(),
	}

	result, err := s.collection.InsertOne(ctx, expense)
	if err != nil {
		return nil, err
	}
	expense.ID = result.InsertedID.(primitive.ObjectID)
	return expense, nil
}

func (s *Service) GetTotalByMonthYear(ctx context.Context, year int, month time.Month) (int, error) {
	var totalCount int64
	filter := map[string]interface{}{
		"created_at": map[string]interface{}{
			"$gte": time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location()),
			"$lt":  time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 1, 0),
		},
	}
	totalCount, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(totalCount), nil
}

func (s *Service) GetExpensesByMonthYear(ctx context.Context, year int, month time.Month, page, pageSize int) ([]Expense, error) {
	location := time.Now().Location()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	filter := map[string]interface{}{
		"created_at": map[string]interface{}{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
	}

	findOptions := &options.FindOptions{}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(map[string]interface{}{"created_at": -1})

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (s *Service) GetTotalByPayerMonthYear(ctx context.Context, payer string, year int, month time.Month) (int, error) {
	matchStage := map[string]interface{}{
		"$match": map[string]interface{}{
			"payer": payer,
			"created_at": map[string]interface{}{
				"$gte": time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location()),
				"$lt":  time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 1, 0),
			},
		},
	}
	groupStage := map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":   nil,
			"total": map[string]interface{}{"$sum": "$amount"},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, []interface{}{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Total int `bson:"total"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].Total, nil
}

package expense

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Description string             `bson:"description" json:"description"`
	Amount      float64            `bson:"amount" json:"amount"`
	Payer       string             `bson:"payer" json:"payer"` // "A" hoặc "B"
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

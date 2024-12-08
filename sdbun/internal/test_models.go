package internal

import (
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"time"
)

type todoEntity1 struct {
	bun.BaseModel `json:"-" bun:"table:t_todo,alias:u"`
	Id_           int64     `json:"id_" bun:"id_,pk,autoincrement"`
	Id            string    `json:"id" bun:"id,unique" repoid:"true"`
	Title         string    `json:"title" bun:"title"`
	Description   string    `json:"description" bun:"description"`
	IsComplete    bool      `json:"is_complete" bun:"is_complete"`
	DueDate       time.Time `json:"due_date" bun:"due_date"`
	Priority      int       `json:"priority" bun:"priority"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at"`
}

func (e *todoEntity1) contentEqual(other *todoEntity1) bool {
	if e == nil && other == nil {
		return true
	}
	if e == nil || other == nil {
		return false
	}
	// 其余字段是数据自动生成的，不用于比较
	return e.Id == other.Id &&
		e.Title == other.Title &&
		e.IsComplete == other.IsComplete &&
		e.Priority == other.Priority
}

func todoEntity1ContainsById(entities []*todoEntity1, id string) bool {
	return lo.ContainsBy(entities, func(e *todoEntity1) bool {
		return e.Id == id
	})
}

func todoEntity1ContainsByTitle(entities []*todoEntity1, title string) bool {
	return lo.ContainsBy(entities, func(e *todoEntity1) bool {
		return e.Title == title
	})
}

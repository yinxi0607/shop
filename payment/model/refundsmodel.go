package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RefundsModel = (*customRefundsModel)(nil)

type (
	// RefundsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundsModel.
	RefundsModel interface {
		refundsModel
		withSession(session sqlx.Session) RefundsModel
	}

	customRefundsModel struct {
		*defaultRefundsModel
	}
)

// NewRefundsModel returns a model for the database table.
func NewRefundsModel(conn sqlx.SqlConn) RefundsModel {
	return &customRefundsModel{
		defaultRefundsModel: newRefundsModel(conn),
	}
}

func (m *customRefundsModel) withSession(session sqlx.Session) RefundsModel {
	return NewRefundsModel(sqlx.NewSqlConnFromSession(session))
}

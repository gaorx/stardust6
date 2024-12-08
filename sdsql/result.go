package sdsql

import (
	"database/sql"
	"database/sql/driver"
)

// Result 扩展 sql.Result 接口，加上一些扩展功能
type Result struct {
	sql.Result
}

// ResultNoRows 没有任何记录的返回结果
var ResultNoRows = ResultOf(driver.ResultNoRows)

// ResultOf 创建一个 Result 对象
func ResultOf(sr sql.Result) Result {
	return Result{sr}
}

// ResultOfRowsAffected 创建一个 Result 对象，只有影响的行数，它不支持LastInsertId
func ResultOfRowsAffected(rowsAffected int64) Result {
	return ResultOf(driver.RowsAffected(rowsAffected))
}

// ResultOfLastInsertIdAndRowsAffected 创建一个 Result 对象，包含 LastInsertId 和 RowsAffected
func ResultOfLastInsertIdAndRowsAffected(lastInsertId int64, rowsAffected int64) Result {
	return ResultOf(sqlResult{lastInsertId: lastInsertId, rowsAffected: rowsAffected})
}

// RowsAffectedAndLastInsertId 返回 RowsAffected 和 LastInsertId
func (sr Result) RowsAffectedAndLastInsertId() (int64, int64, error) {
	affected, err := sr.Result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	lastInsertId, err := sr.Result.LastInsertId()
	if err != nil {
		return affected, 0, err
	}
	return affected, lastInsertId, nil
}

// RowsAffectedOr 返回 RowsAffected，如果出错则返回默认值
func (sr Result) RowsAffectedOr(def int64) int64 {
	n, err := sr.Result.RowsAffected()
	if err != nil {
		return def
	}
	return n
}

// LastInsertIdOr 返回 LastInsertId，如果出错则返回默认值
func (sr Result) LastInsertIdOr(def int64) int64 {
	n, err := sr.Result.LastInsertId()
	if err != nil {
		return def
	}
	return n
}

// RowsAffectedAndLastInsertIdOr 返回 RowsAffected 和 LastInsertId，如果出错则返回默认值
func (sr Result) RowsAffectedAndLastInsertIdOr(affectedDef, lastInsertIdDef int64) (int64, int64) {
	affected, lastInsertId, err := sr.RowsAffectedAndLastInsertId()
	if err != nil {
		return affectedDef, lastInsertIdDef
	}
	return affected, lastInsertId
}

type sqlResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (sr sqlResult) LastInsertId() (int64, error) {
	return sr.lastInsertId, nil
}

func (sr sqlResult) RowsAffected() (int64, error) {
	return sr.rowsAffected, nil
}

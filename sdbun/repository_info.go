package sdbun

import (
	"github.com/gaorx/stardust6/sdparse"
	"github.com/gaorx/stardust6/sdreflect"
	"github.com/gaorx/stardust6/sdsql"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
	"reflect"
)

type repoInfo struct {
	table          *schema.Table
	idFieldSqlName string
	idFieldGoName  string
	idFieldGoIndex []int
}

type repoInfoLoaded struct {
	info *repoInfo
	ok   bool
}

var repoInfos = xsync.NewMapOf[reflect.Type, repoInfoLoaded]()

func getRepoInfo(db bun.IDB, modelType reflect.Type) *repoInfo {
	if sdreflect.IsStruct(modelType) {
		// OK
	} else if sdreflect.IsStructPtr(modelType) {
		modelType = modelType.Elem()
	} else {
		return nil
	}

	loadedInfo, _ := repoInfos.LoadOrCompute(modelType, func() repoInfoLoaded {
		table := db.Dialect().Tables().Get(modelType)
		if table.Type != modelType {
			panic("modelType is not a registered model")
		}
		var repoIdField *schema.Field
		numField := modelType.NumField()
		for i := 0; i < numField; i++ {
			f := modelType.Field(i)
			repoId := f.Tag.Get("repoid")
			if sdparse.Bool(repoId) {
				tf := getTableFieldByGoName(table, f.Name)
				if tf != nil {
					repoIdField = tf
					break
				}
			}
		}
		if repoIdField == nil {
			pks := table.PKs
			if len(pks) == 1 {
				repoIdField = pks[0]
			}
		}
		if repoIdField == nil {
			return repoInfoLoaded{ok: false}
		}

		info := &repoInfo{
			table:          table,
			idFieldSqlName: repoIdField.Name,
			idFieldGoName:  repoIdField.GoName,
			idFieldGoIndex: repoIdField.StructField.Index,
		}
		return repoInfoLoaded{info: info, ok: true}
	})
	if !loadedInfo.ok {
		return nil
	}
	return loadedInfo.info
}

func getEntityId[ID sdsql.EntityId](entity any, info *repoInfo) ID {
	entityVal := reflect.ValueOf(entity)
	if entityVal.Kind() == reflect.Ptr {
		entityVal = entityVal.Elem()
	}
	return entityVal.FieldByIndex(info.idFieldGoIndex).Interface().(ID)
}

func getTableFieldByGoName(t *schema.Table, goFieldName string) *schema.Field {
	for _, f := range t.Fields {
		if f.GoName == goFieldName {
			return f
		}
	}
	return nil
}

package sdcsv

import (
	"bytes"
	"encoding/csv"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdlo"
	"io"
	"os"
	"slices"
	"strings"
)

// Reader CSV读取器
type Reader struct {
	reader *csv.Reader
	fields []string
}

// Options 选项
type Options struct {
	// Header CSV是否有标题行
	Header bool
	// Fields 要读取的字段
	Fields []string
	// Comma 字段分隔符
	Comma rune
	// Comment 注释符
	Comment rune
	// FieldsPerRecord 每行期望的字段数
	FieldsPerRecord int
	// LazyQuotes 是否懒惰引号
	LazyQuotes bool
	// TrimLeadingSpace 是否去掉前导空格
	TrimLeadingSpace bool
	// ReuseRecord 是否重用记录
	ReuseRecord bool
}

// HandlerResult Handler的返回值
type HandlerResult int

const (
	// Stop 停止读取
	Stop HandlerResult = 0
	// Continue 继续读取
	Continue HandlerResult = 1
)

// NewReader 创建读取器
func NewReader(r io.Reader, opts *Options) (*Reader, error) {
	r1 := csv.NewReader(r)
	var fields []string
	if opts != nil {
		if opts.Comma != 0 {
			r1.Comma = opts.Comma
		}
		if opts.Comma != 0 {
			r1.Comment = opts.Comment
		}
		if opts.FieldsPerRecord > 0 {
			r1.FieldsPerRecord = opts.FieldsPerRecord
		}
		r1.LazyQuotes = opts.LazyQuotes
		r1.TrimLeadingSpace = opts.TrimLeadingSpace
		r1.ReuseRecord = opts.ReuseRecord
	}
	if opts != nil && len(opts.Fields) > 0 {
		fields = slices.Clone(opts.Fields)
	}
	if opts != nil && opts.Header {
		header, err := r1.Read()
		if err != nil {
			return nil, sderr.Wrapf(err, "read CSV error")
		}
		header = slices.Clone(header)
		if len(fields) == 0 {
			fields = header
		}
	}
	return &Reader{
		reader: r1,
		fields: fields,
	}, nil
}

// NewReaderBytes 从二进制内容创建读取器
func NewReaderBytes(b []byte, opts *Options) (*Reader, error) {
	return NewReader(bytes.NewReader(sdlo.EnsureSlice(b)), opts)
}

// NewReaderText 从文本内容创建读取器
func NewReaderText(s string, opts *Options) (*Reader, error) {
	return NewReader(strings.NewReader(s), opts)
}

// NewReaderFile 从文件创建读取器
func NewReaderFile(filename string, opts *Options) (*Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, sderr.Wrapf(err, "open CSV file error")
	}
	return NewReader(f, opts)
}

// Fields 获取字段列表
func (r *Reader) Fields() []string {
	return r.fields
}

// ReadRecord 读取一行记录
func (r *Reader) ReadRecord() ([]string, error) {
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.Wrapf(err, "read CSV record error")
	}
	return rec, nil
}

// ReadRecords 读取所有记录
func (r *Reader) ReadRecords() ([][]string, error) {
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.Wrapf(err, "read CSV all record error")
	}
	return recs, nil
}

// ReadMap 读取一行记录为map
func (r *Reader) ReadMap() (map[string]string, error) {
	if len(r.fields) <= 0 {
		return nil, sderr.Newf("no CSV field")
	}
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.Wrapf(err, "read CSV record error")
	}
	return makeMap(r.fields, rec), nil
}

// ReadMaps 读取所有记录为map列表
func (r *Reader) ReadMaps() ([]map[string]string, error) {
	if len(r.fields) <= 0 {
		return nil, sderr.Newf("no CSV field")
	}
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.Wrapf(err, "read CSV all record error")
	}
	maps := make([]map[string]string, 0, len(recs))
	for _, rec := range recs {
		maps = append(maps, makeMap(r.fields, rec))
	}
	return maps, nil
}

// ForeachRecord 遍历记录
func (r *Reader) ForeachRecord(h func(recNo int, rec []string) HandlerResult) error {
	recNo := 0
	for {
		rec, err := r.ReadRecord()
		if err != nil {
			if sderr.Is(err, io.EOF) {
				break
			} else {
				return sderr.Wrapf(err, "foreach CSV record error")
			}
		}
		hr := h(recNo, rec)
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

// ForeachMap 以map形式遍历记录
func (r *Reader) ForeachMap(h func(recNo int, rec map[string]string) HandlerResult) error {
	if len(r.fields) == 0 {
		return sderr.Newf("no CSV field")
	}
	recNo := 0
	for {
		rec, err := r.ReadRecord()
		if err != nil {
			if sderr.Is(err, io.EOF) {
				break
			} else {
				return sderr.Wrapf(err, "read CSV record error")
			}
		}
		hr := h(recNo, makeMap(r.fields, rec))
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

func makeMap(fields, record []string) map[string]string {
	fieldNum, valNum := len(fields), len(record)
	m := make(map[string]string, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := fields[i]
		v := ""
		if i < valNum {
			v = record[i]
		}
		m[field] = v
	}
	return m
}

package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	xerrors "github.com/pkg/errors"
	error2 "github.com/sp187/geekbang/final/internal/framework/error"
)

var _defaultError = error2.Server

// SetResponseDefaultError 设置默认的response错误类型
func SetResponseDefaultError(err *error2.Error) {
	_defaultError = err
}

func BuildErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var e *error2.Error
	if !errors.As(err, &e) {
		// no *RSError in err, use _defaultError
		w.WriteHeader(_defaultError.HttpCode)
		json.NewEncoder(w).Encode(_defaultError)
		return
	} else {
		w.WriteHeader(err.(*error2.Error).HttpCode)
		json.NewEncoder(w).Encode(err.(*error2.Error))
		return
	}
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func BuildOKResponse(w http.ResponseWriter, data interface{}) {
	var (
		err error
	)
	defer func() {
		if err != nil {
			// json序列化失败
			log.Printf("serialize response fail: %s", err.Error())
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if data == nil {
		ret := make(map[string]interface{})
		err = json.NewEncoder(w).Encode(SuccessResponse{Code: 0, Data: ret})
		return
	}
	rt := reflect.TypeOf(data)
	for rt.Kind() == reflect.Ptr {
		data = reflect.Indirect(reflect.ValueOf(data)).Interface()
		rt = reflect.TypeOf(data)
	}
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		if _, ok := data.([]interface{}); !ok {
			rd := reflect.ValueOf(data)
			res := make([]interface{}, rd.Len())
			for i := 0; i < rd.Len(); i++ {
				res[i] = rd.Index(i).Interface()
			}
			err = json.NewEncoder(w).Encode(SuccessResponse{Code: 0, Data: res})
		} else {
			err = json.NewEncoder(w).Encode(SuccessResponse{Code: 0, Data: data.([]interface{})})
		}
	default:
		err = json.NewEncoder(w).Encode(SuccessResponse{Code: 0, Data: data})
	}
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	vars    map[string]string
}

// NewContext 创建Context对象
func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{Writer: writer, Request: request}
}

// Context 返回context.Context
func (c *Context) Context() context.Context {
	return c.Request.Context()
}

// ReadJson 读取http请求中的json格式的body内容
func (c *Context) ReadJson(data interface{}) error {
	defer c.Request.Body.Close()
	return json.NewDecoder(c.Request.Body).Decode(&data)
}

func (c *Context) WriteJson(status int, data interface{}) error {
	err := json.NewEncoder(c.Writer).Encode(data)
	if err != nil {
		return err
	}
	c.Writer.WriteHeader(status)
	return nil
}

// BadJsonResponse 根据错误内容返回对应的错误状态码
func (c *Context) BadJsonResponse(err error) {
	BuildErrorResponse(c.Writer, xerrors.Cause(err))
}

// OKJsonResponse 返回http状态码为200，并且body为json格式的内容。
func (c *Context) OKJsonResponse(data interface{}) {
	BuildOKResponse(c.Writer, data)
}

// Query 获取url中的查询参数
func (c *Context) Query() url.Values {
	return c.Request.URL.Query()
}

// Var 获取路由中的参数
func (c *Context) Var(key string) string {
	return c.vars[key]
}

func (c *Context) WriteFileStream(f FileStream, name string) {
	defer f.Close()
	var err error
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+name)
	c.Writer.Header().Set("Accept-ranges", "bytes")
	var start, end int64
	if r := c.Request.Header.Get("Range"); r != "" {
		if strings.Contains(r, "bytes=") && strings.Contains(r, "-") {
			fmt.Sscanf(r, "bytes=%d-%d", &start, &end)
			if end == 0 {
				end = f.Size() - 1
			}
			if start > end || start < 0 || end < 0 || end >= f.Size() {
				fmt.Printf("文件范围错误, start:%d, end:%d, total size:%d\n", start, end, f.Size())
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			c.Writer.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
			c.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes %v-%v/%v", start, end, f.Size()))
			c.Writer.WriteHeader(http.StatusPartialContent)
		} else {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		start = 0
		end = f.Size() - 1
		c.Writer.Header().Set("Content-Length", strconv.FormatInt(f.Size(), 10))
		c.Writer.WriteHeader(http.StatusOK)
	}
	if c.Request.Method == http.MethodGet {
		_, err = f.Seek(start, 0)
		if err != nil {
			fmt.Println(err.Error())
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Sprintf("from %d to %d\n", start, end)
		size := 2 * 1024 * 1024
		buffer := make([]byte, size)
		for {
			if end-start+1 < int64(size) {
				size = int(end - start + 1)
			}
			if size <= 0 {
				return
			}
			_, err = f.Read(buffer[:size])
			if err != nil {
				if err == io.EOF {
					return
				} else {
					fmt.Println(err.Error())
					c.Writer.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			_, err = c.Writer.Write(buffer[:size])
			if err != nil {
				fmt.Println(err.Error())
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			start += int64(size)
		}
	}
	return
}

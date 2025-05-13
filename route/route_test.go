package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mail2fish/gorails/errors"
	"github.com/stretchr/testify/assert"
)

const (
	MODULE_TEST errors.ErrorModule = 999
)

// 测试用的参数结构体
type TestParams struct {
	Name string `json:"name" binding:"required"`
}

func (p *TestParams) Parse(c *gin.Context) errors.Error {
	if err := c.ShouldBindJSON(p); err != nil {
		return errors.NewError(http.StatusBadRequest, errors.THIRD_PARTY, MODULE_TEST, 1, "无效的请求参数", err)
	}
	return nil
}

// 测试用的响应结构体
type TestResponse struct {
	Message string `json:"message"`
}

func (r *TestResponse) Render(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

// 测试用的处理函数
func mockHandler(c *gin.Context, params Params) (Response, errors.Error) {
	p := params.(*TestParams)
	return &TestResponse{Message: "Hello, " + p.Name}, nil
}

func TestWrap(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "成功处理请求",
			requestBody: TestParams{
				Name: "World",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Hello, World",
			},
		},
		{
			name:           "无效的请求体",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    "2-999-1",
				"message": "无效的请求参数",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.POST("/test", Wrap[*TestParams, *TestResponse](mockHandler))

			// 创建测试请求
			reqBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// 发送请求
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// 验证状态码
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 验证响应体
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

// 测试 Wrap 并发时，参数传不同的值会不会有问题
func TestWrapConcurrent(t *testing.T) {
	router := gin.New()
	router.POST("/test", Wrap[*TestParams, *TestResponse](mockHandler))

	// 创建多个并发请求
	numRequests := 1000
	wg := sync.WaitGroup{}
	wg.Add(numRequests * 2) // 增加错误请求的测试

	// 测试正常请求
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer wg.Done()
			name := fmt.Sprintf("World-%d", index)
			reqBody := fmt.Sprintf(`{"name": "%s"}`, name)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte(reqBody)))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, http.StatusOK, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Hello, %s", name), response["message"])
		}(i)
	}

	// 测试错误请求
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{"invalid": "data"}`)))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// 验证错误响应
			assert.Equal(t, http.StatusBadRequest, w.Code)
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "2-999-1", response["code"])
			assert.Equal(t, "无效的请求参数", response["message"])
		}()
	}

	wg.Wait()
}

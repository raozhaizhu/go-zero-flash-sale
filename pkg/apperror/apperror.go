package apperror

import (
	"fmt"
	"path/filepath"
	"runtime"
)

/** ====================================================================================
 * 🏁 BizError
 * =====================================================================================
 *
 */

type BizError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	err  error
}

func New(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

func (e *BizError) Error() string {
	if e.err != nil {
		// 如果有底层错误，把内外两层错误拼起来，方便日志打印
		return fmt.Sprintf("ErrCode: %d, ErrMsg: %s, Cause: %v", e.Code, e.Msg, e.err)
	}
	return fmt.Sprintf("ErrCode: %d, ErrMsg: %s", e.Code, e.Msg)
}

func (e *BizError) Unwrap() error {
	return e.err
}

func (e *BizError) Is(targetErr error) bool {
	// 尝试转化为 BizError
	t, ok := targetErr.(*BizError)
	if !ok {
		return false
	}
	// 判断错误代码是否相同
	return e.Code == t.Code
}

func (e *BizError) WithErr(err error) *BizError {
	return &BizError{
		Code: e.Code,
		Msg:  e.Msg,
		err:  err,
	}
}

func NewSrvErr(err error) *BizError {
	return ErrServerErr.WithErr(err)
}

func WrapSrvErrf(err error, format string, args ...interface{}) *BizError {
	if err == nil {
		return nil
	}

	// 获取调用信息
	_, file, line, ok := runtime.Caller(1)
	callerInfo := "unknown"
	if ok {
		// 提取出简短的文件名和行号，如 "register.go:45"
		callerInfo = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// 格式化传入信息, 辅以调用信息
	inputMsg := fmt.Sprintf(format, args...)
	errMsg := fmt.Sprintf("[%s] %s", callerInfo, inputMsg)

	// 嵌套错误
	wrappedErr := fmt.Errorf("%s: %w", errMsg, err)

	// 打包进 ErrServerEr
	return ErrServerErr.WithErr(wrappedErr)
}

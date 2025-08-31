package middlewares

import (
	"net/http"
	"questionbasket/frame"
)

type CheckContentType struct {
}

func NewCheckContentTypeMiddleware() CheckContentType {
	return CheckContentType{}
}

func (cct CheckContentType) RunMiddleware(rs http.ResponseWriter, rq *http.Request) frame.APIMiddlewareResult {
	if rq.Method == http.MethodPost {
		if rq.Header.Get("Content-Type") != "application/json" {
			return frame.APIMiddlewareResult{
				IsSuccess: false,
				ApiError:  frame.NewAPIError("clienterror", "incorrect Content-Type", http.StatusUnsupportedMediaType),
			}
		} else {
			return frame.APIMiddlewareResult{
				IsSuccess: true,
			}
		}
	} else {
		return frame.APIMiddlewareResult{
			IsSuccess: true,
		}
	}
}

func (cct CheckContentType) GetMiddlewareInfo() frame.APIMiddlewareInfo {
	return frame.APIMiddlewareInfo{
		MiddlewareName: "ContentTypeCheck",
	}
}

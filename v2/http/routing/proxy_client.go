package gmrouting

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type responseWriter interface {
	WriteCustomJsonResponse(w http.ResponseWriter, statusCode int, res interface{}) (writtenRes []byte, err error)
}

type ProxyClient struct {
	routeTable     *RouteTable
	routeUrl       string
	httpCli        *http.Client
	responseWriter responseWriter
	ignoredPaths   map[string]bool

	onErr     func(context.Context, error)
	onReqRead func(context.Context, []byte)
	onResRead func(context.Context, []byte)
}

// NewProxyClient creates a new proxy client instance.
// Underlying HandleRequestAndRedirect method can be registered as a handler function.
// Handler function will redirect incoming request to the routeUrl.
//
// ignoredPaths will return 200 without any other http content.
// (ignoredPaths must be exact paths. Regex is not supported.)
//
// Since it will be registered to a blocking function (ListAndServe), internally occuring event data can be captured via following hooks:
//
// onErr: Can be registered to receive internal errors.
//
// onReqRead: Can be registered to get incoming request body.
//
// onResRead: Can be registered to get outgoing response body.
//
// Hooks will contain a second string value which represents session ID.
// Events which output the same session ID belong to same http session.
func NewProxyClient(routeTable *RouteTable, routeUrl string, httpCli *http.Client, responseWriter responseWriter, ignoredPaths []string, onErr func(context.Context, error), onReqRead func(context.Context, []byte), onResRead func(context.Context, []byte)) *ProxyClient {
	pc := &ProxyClient{
		routeTable:     routeTable,
		routeUrl:       routeUrl,
		httpCli:        httpCli,
		responseWriter: responseWriter,

		onErr:     onErr,
		onReqRead: onReqRead,
		onResRead: onResRead,
	}

	pc.ignoredPaths = make(map[string]bool, len(ignoredPaths))
	for _, path := range ignoredPaths {
		pc.ignoredPaths[path] = true
	}
	return pc
}

// HandleRequestAndRedirect can be registered to http.Handle() for redirecting requests to desired url.
func (pc *ProxyClient) HandleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	if pc.ignoredPaths[r.URL.RequestURI()] {
		w.Write(nil)
		return
	}

	uri := r.URL.RequestURI()

	isAllowed := false

	for _, e := range pc.routeTable.routeRules {
		if e.method != r.Method {
			continue
		}

		regexConv, err := RouteToRegExp(uri)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), err)
			}
			writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal server error",
			})
			if err != nil {
				if pc.onErr != nil {
					pc.onErr(r.Context(), err)
				}
				return
			}
			if pc.onResRead != nil {
				pc.onResRead(r.Context(), writtenRes)
			}
		}

		if e.regexp.MatchString(regexConv) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		if pc.onErr != nil {
			pc.onErr(r.Context(), fmt.Errorf("path is not allowed: %s", uri))
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"message": "unauthorized call",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(r.Context(), writtenRes)
		}
		return
	}

	redirectUrl := pc.routeUrl + r.URL.RequestURI()
	parsedRedirectUrl, err := url.Parse(redirectUrl)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(r.Context(), fmt.Errorf("unable to parse URL: '%s' error: %s", redirectUrl, err.Error()))
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(r.Context(), writtenRes)
		}
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(r.Context(), fmt.Errorf("error reading request body: %s", err.Error()))
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(r.Context(), writtenRes)
		}
		return
	}

	if pc.onReqRead != nil {
		pc.onReqRead(r.Context(), reqBytes)
	}

	buffer := bytes.NewBuffer(reqBytes)
	nopCloser := io.NopCloser(buffer)

	httpReq := &http.Request{
		Method: r.Method,
		URL:    parsedRedirectUrl,
		Header: r.Header,
		Body:   nopCloser,
	}

	httpRes, err := pc.httpCli.Do(httpReq)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(r.Context(), fmt.Errorf("error executing http request: %s", err.Error()))
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(r.Context(), writtenRes)
		}
		return
	}

	resBytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(r.Context(), fmt.Errorf("error reading response payload: %s", err.Error()))
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(r.Context(), writtenRes)
		}
		return
	}
	defer httpRes.Body.Close()

	for k, v := range httpRes.Header {
		for i := 0; i < len(v); i++ {
			w.Header().Add(k, v[i])
		}
	}

	if httpRes.StatusCode != http.StatusOK {
		w.WriteHeader(httpRes.StatusCode)
	}

	_, err = w.Write(resBytes)
	if err != nil {
		if pc.onErr != nil {
			pc.onErr(r.Context(), fmt.Errorf("error writing server response for client: %s", err.Error()))
		}
		writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "internal error",
		})
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
			}
			return
		}
		if pc.onResRead != nil {
			pc.onResRead(r.Context(), writtenRes)
		}
		return
	}

	if httpRes.Header.Get("Content-Encoding") == "gzip" {
		reader := bytes.NewReader(resBytes)
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("error creating gzip reader: %s", err.Error()))
			}
			writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal error",
			})
			if err != nil {
				if pc.onErr != nil {
					pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
				}
				return
			}
			if pc.onResRead != nil {
				pc.onResRead(r.Context(), writtenRes)
			}
			return
		}
		/* Modifying resBytes for logging decompressed content AFTER we've written the response body. */
		resBytes, err = io.ReadAll(gzipReader)
		if err != nil {
			if pc.onErr != nil {
				pc.onErr(r.Context(), fmt.Errorf("error reading from gzip reader: %s", err.Error()))
			}
			writtenRes, err := pc.responseWriter.WriteCustomJsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "internal error",
			})
			if err != nil {
				if pc.onErr != nil {
					pc.onErr(r.Context(), fmt.Errorf("write response error: %s", err.Error()))
				}
				return
			}
			if pc.onResRead != nil {
				pc.onResRead(r.Context(), writtenRes)
			}
			return
		}
	}
	if pc.onResRead != nil {
		pc.onResRead(r.Context(), resBytes)
	}
}

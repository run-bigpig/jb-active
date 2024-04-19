package handler

import (
	"encoding/json"
	"github.com/run-bigpig/jb-active/internal/license"
	"github.com/run-bigpig/jb-active/internal/model"
	"github.com/run-bigpig/jb-active/internal/response"
	"github.com/valyala/fasthttp"
)

func LicenseHandler(ctx *fasthttp.RequestCtx) {
	var licenseReq model.License
	err := json.Unmarshal(ctx.PostBody(), &licenseReq)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	lic, err := license.GenerateLicense(&licenseReq)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, model.ResponseData{
		License: lic,
	})
}

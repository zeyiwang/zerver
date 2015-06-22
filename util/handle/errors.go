package handle

import (
	"github.com/cosiner/gohper/errors/httperrs"
	"github.com/cosiner/ygo/log"
	"github.com/cosiner/zerver"
)

var (
	Logger   log.Logger
	KeyError = "error"
)

func Wrap(handle func(zerver.Request, zerver.Response) error) zerver.HandleFunc {
	return func(req zerver.Request, resp zerver.Response) {
		if err := handle(req, resp); err != nil {
			SendErr(resp, err)
		}
	}
}

func SendErr(resp zerver.Response, err error) {
	switch err := err.(type) {
	case httperrs.Error:
		resp.ReportStatus(err.Code())
		if err.Code() < int(httperrs.Server) {
			OnErrLog(resp.Send(KeyError, err.Error()))
			return
		}
	default:
		resp.ReportInternalServerError()
	}

	Logger.Errorln(err.Error())
}

func SendBadRequest(resp zerver.Response, err error) {
	SendErr(resp, httperrs.BadRequest.New(err))
}

func OnErrLog(err error) {
	if err != nil {
		Logger.Errorln(err)
	}
}
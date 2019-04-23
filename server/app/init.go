package app

import (
	//"fmt"
	"time"
	"net/http"
	//"errors"

	"github.com/Sirupsen/logrus"
	"github.com/syronz/ozzo-routing"
	"github.com/syronz/ozzo-routing/access"
	"github.com/syronz/ozzo-routing/fault"


)

func Init(logger *logrus.Logger) routing.Handler {
	return func(rc *routing.Context) error {
		now := time.Now()

		rc.Response = &access.LogResponseWriter{rc.Response, http.StatusOK, 0}

		ac := newRequestScope(now, logger, rc.Request)
		rc.Set("Context", ac)

		fault.Recovery(ac.Errorf, convertError)(rc)

		return nil
	}
}

func GetRequestScope(c *routing.Context) RequestScope {
	return c.Get("Context").(RequestScope)
}


func convertError(c *routing.Context, err error) error {
	return err
}


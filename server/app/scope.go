package app

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

type RequestScope interface {
	Logger
	UserID() string
	UserLanguage() string
	SetUserID(id string)
	SetUserLanguage(language string)
	RequestID() string
}

type requestScope struct {
	Logger
	now			time.Time
	requestID	string
	userID		string
	userLanguage string
}

func (rs *requestScope) UserID() string {
	return rs.userID
}

func (rs *requestScope) UserLanguage() string {
	return rs.userLanguage
}

func (rs *requestScope) SetUserID(id string) {
	//rs.Logger.SetField("UserID", id)
	rs.userID = id
}

func (rs *requestScope) SetUserLanguage(language string) {
	rs.userLanguage = language
}

func (rs *requestScope) RequestID() string {
	return rs.requestID
}

func newRequestScope(now time.Time, logger *logrus.Logger, request *http.Request) RequestScope {
	l := NewLogger(logger, logrus.Fields{})
	requestID := request.Header.Get("X-Request-Id")
	if requestID != "" {
		l.SetField("RequestID", requestID)
	}
	return &requestScope{
		Logger:		l,
		now:		now,
		requestID:	requestID,
	}
}

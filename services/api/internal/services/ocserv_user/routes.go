package ocserv_user

import (
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/ocserv-dashboard/api/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/ocserv/users", middlewares.AuthMiddleware())

	g.GET("", ctl.OcservUsers)
	g.GET("/:uid", ctl.OcservUser)
	g.POST("", ctl.CreateOcservUser)
	g.PATCH("/:uid", ctl.UpdateOcservUser)
	g.DELETE("/:uid", ctl.DeleteOcservUser)
	g.POST("/:uid/lock", ctl.LockOcservUser)
	g.POST("/:uid/unlock", ctl.UnLockOcservUser)
	g.POST("/:uid/activate", ctl.ActivateExpiredOcservUsers)
	g.POST("/:username/disconnect", ctl.DisconnectOcservUser)
	g.GET("/:uid/session_logs", ctl.OcservUserSessionLogs)
	g.GET("/:uid/statistics", ctl.OcservUserStatistics)

	g.GET("/ocpasswd", ctl.OcpasswdUsers, middlewares.AdminPermission())
	g.POST("/ocpasswd/sync", ctl.SyncToDB, middlewares.AdminPermission())
}

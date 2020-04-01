package main

import (
	"context"
	"net/http"

	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	fateConfig "github.com/godcong/fate/config"
	"github.com/labstack/echo/v4"
)

func fateName(ctx echo.Context) error {
	fateCfg := fateConfig.DefaultConfig()
	fateCfg.Database = cfg.Database.toFateDatabase()

	//生日：
	born := chronos.New(ctx.QueryParam("born")) // chronos.New("2020/01/23 11:31")
	//姓氏：
	lastName := ctx.QueryParam("last_name") //"张"

	//第一参数：姓氏
	//第二参数：生日
	f := fate.NewFate(lastName, born.Solar().Time(), fate.ConfigOption(fateCfg))

	if err := f.MakeName(context.Background()); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"name": "",
	})
}

func runServer() {
	e := echo.New()

	e.GET("/fatename", fateName)

	e.Logger.Fatal(e.Start(cfg.Server.Addr))
}

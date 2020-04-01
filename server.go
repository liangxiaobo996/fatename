package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/godcong/chronos"
	"github.com/godcong/fate"
	fateConfig "github.com/godcong/fate/config"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func fateName(ctx echo.Context) error {
	fateCfg := fateConfig.DefaultConfig()
	fateCfg.BaguaFilter = true
	fateCfg.SupplyFilter = true
	fateCfg.ZodiacFilter = true
	fateCfg.Database = cfg.Database.toFateDatabase()

	//生日：
	born := chronos.New(ctx.QueryParam("born")) // chronos.New("2020/01/23 11:31")
	//姓氏：
	lastName := ctx.QueryParam("last_name") //"张"

	// 临时文件
	fateCfg.FileOutput.OutputMode = fateConfig.OutputModelJSON
	fateCfg.FileOutput.Path = fmt.Sprintf("names-%s.txt", uuid.New().String())

	//第一参数：姓氏
	//第二参数：生日
	f := fate.NewFate(lastName, born.Solar().Time(), fate.ConfigOption(fateCfg))

	if err := f.MakeName(context.Background()); err != nil {
		zap.L().Error("make fail", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	content, err := ioutil.ReadFile(fateCfg.FileOutput.Path)

	if err != nil {
		zap.L().Error("read file fail", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// 删除临时文件
	go func(path string) {
		os.Remove(path)
	}(fateCfg.FileOutput.Path)

	return ctx.JSONBlob(http.StatusOK, content)
}

func runServer() {
	e := echo.New()
	e.GET("/fatename", fateName)
	e.Logger.Fatal(e.Start(cfg.Server.Addr))
}

package main

import (
	"os"

	fateConfig "github.com/godcong/fate/config"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

type Config struct {
	Debug    bool           `yaml:"debug" default:"false"`
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Pwd          string `yaml:"pwd"`
	Name         string `yaml:"name"`
	MaxIdleCon   int    `yaml:"max_idle_con"`
	MaxOpenCon   int    `yaml:"max_open_con"`
	Driver       string `yaml:"driver" default:"mysql"`
	File         string `yaml:"file"`
	Dsn          string `yaml:"dsn"`
	ShowSQL      bool   `yaml:"show_sql"`
	ShowExecTime bool   `yaml:"show_exec_time"`
}

func (d DatabaseConfig) Addr() string {
	return d.Host + ":" + d.Port
}
func (d DatabaseConfig) toFateDatabase() fateConfig.Database {
	return fateConfig.Database{
		Host:         d.Host,
		Port:         d.Port,
		User:         d.User,
		Pwd:          d.Pwd,
		Name:         d.Name,
		MaxIdleCon:   d.MaxIdleCon,
		MaxOpenCon:   d.MaxOpenCon,
		Driver:       d.Driver,
		File:         d.File,
		Dsn:          d.Dsn,
		ShowSQL:      d.ShowSQL,
		ShowExecTime: d.ShowExecTime,
	}
}

type LoggerConfig struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	// json or text
	Format string `yaml:"format,omitempty" default:"json"`
	// file
	Output string `yaml:"output,omitempty" default:""`
}

type ServerConfig struct {
	Addr string `yaml:"addr,omitempty" default:":1234"`
}

var cfg *Config

func initConfig(cfgFile string) {
	_ = os.Setenv("CONFIGOR_ENV_PREFIX", "-")

	cfg = new(Config)

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	if cfgFile != "" {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(cfg, cfgFile); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	} else {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(cfg); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	}

	zap.L().Debug("loaded config", zap.Any("config", cfg), zap.String("cfgFile", cfgFile))
}

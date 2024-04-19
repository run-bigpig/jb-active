package bootstrap

import (
	"github.com/fasthttp/router"
	"github.com/run-bigpig/jb-active/internal/cert"
	"github.com/run-bigpig/jb-active/internal/power"
	"github.com/run-bigpig/jb-active/internal/product"
	router2 "github.com/run-bigpig/jb-active/internal/router"
	"github.com/run-bigpig/jb-active/internal/utils"
	"log"
	"os"
	"path/filepath"
)

func Run() *router.Router {
	log.Println("正在创建相关目录")
	err := createDir()
	if err != nil {
		panic(err)
	}
	log.Println("正在生成证书")
	err = cert.CreateCert()
	if err != nil {
		panic(err)
	}
	log.Println("正在生成power.conf")
	err = power.GenerateEqualResult()
	if err != nil {
		panic(err)
	}
	log.Println("正在同步Jetbrains产品信息")
	err = product.SyncProduct()
	if err != nil {
		panic(err)
	}
	log.Println("正在启动服务")
	return router2.NewRouter()
}

// 创建目录
func createDir() error {
	certDir := filepath.Join(utils.GetCurrentPath(), "cert")
	jsonDir := filepath.Join(utils.GetStaticPath(), "json")
	confDir := filepath.Join(utils.GetStaticPath(), "conf")
	err := os.MkdirAll(certDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(jsonDir, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(confDir, 0755)
	if err != nil {
		return err
	}
	return nil
}

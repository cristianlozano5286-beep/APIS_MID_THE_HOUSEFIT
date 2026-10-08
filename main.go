package main

import (
	"log"

	_ "api_mid_the_housefit/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func main() {
	// Leer la cadena de conexión desde app.conf y conectar
	sqlconn, err := beego.AppConfig.String("sqlconn")
	if err != nil || sqlconn == "" {
		log.Fatal("Falta sqlconn en conf/app.conf")
	}
	if err := orm.RegisterDataBase("default", "postgres", sqlconn); err != nil {
		log.Fatal("Error registrando la base de datos: ", err)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
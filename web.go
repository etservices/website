package main

import (
	"fmt"
	"github.com/hoisie/web"
	"github.com/etservices/website/blog"
	"github.com/etservices/website/conf"
	"github.com/etservices/website/db"
	"github.com/etservices/website/template"
	"reflect"
	// apps
	"github.com/etservices/website/app"
	"github.com/etservices/website/gallery"
)

func main() {
	fmt.Printf("Using config %s\n", conf.Path)
	fmt.Printf("Using models:\n")
	for _, m := range db.Models {
		t := reflect.TypeOf(m)
		fmt.Printf("    %s\n", fmt.Sprintf("%s", t)[1:])
	}

	template.Init(conf.Config.TemplatePaths)
	fmt.Printf(conf.Config.String())
	web.Config.StaticDir = conf.Config.StaticPath

	db.Connect(conf.Config.DbHostString(), conf.Config.DbName)
	db.RegisterAllIndexes()

	if (!app.ValidateUser("admin", "passw0rd")) {
		fmt.Printf("Adding user admin\n")
		err := app.CreateUser("admin", "passw0rd")
		if(err != nil) {
			fmt.Println(err)
		}
	}

	blog.AttachAdmin("/admin/")
	blog.Attach("/")

	app.AttachAdmin("/admin/")
	app.Attach("/")

	gallery.AttachAdmin("/admin/")
	gallery.Attach("/")

	web.Get("/([^/]*)", blog.Index)
	web.Get("/(.*)", blog.Flatpage)

	app.AddPanel(&blog.PostPanel{})
	app.AddPanel(&blog.UnpublishedPanel{})
	app.AddPanel(&blog.PagesPanel{})
	app.AddPanel(&gallery.GalleryPanel{})

	web.Run(conf.Config.HostString())
}

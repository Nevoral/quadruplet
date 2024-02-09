package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"strings"
)

const pathBase = "./web"

func GetJs(c fiber.Ctx) (err error) {
	c.Set("Content-Type", fiber.MIMETextJavaScript)
	path := strings.Split(c.Path(), "/")
	if file := fmt.Sprintf("%s", c.Params("file")); strings.HasSuffix(file, ".js") {
		err = c.SendFile(fmt.Sprintf("%s/static/%s/%s", pathBase, path[1], c.Params("file")))
	} else {
		err = c.SendFile(fmt.Sprintf("%s/static/%s/%s.js", pathBase, path[1], c.Params("file")))
	}
	if err != nil {
		return err
	}
	return nil
}

func GetAsset(c fiber.Ctx) (err error) {
	path := strings.Split(c.Path(), "/")
	if path[1] == "assets" {
		err = c.SendFile(fmt.Sprintf("%s/home/static/%s/%s", pathBase, path[1], c.Params("file")))
	} else {
		err = c.SendFile(fmt.Sprintf("%s/%s/static/%s/%s", pathBase, path[1], path[2], c.Params("file")))
	}
	if err != nil {
		return err
	}
	return nil
}

func GetHTML(page templ.Component) func(c fiber.Ctx) (err error) {
	return func(c fiber.Ctx) (err error) {
		c.Set("Content-Type", fiber.MIMETextHTML)

		var output bytes.Buffer

		err = page.Render(context.Background(), &output)
		if err != nil {
			return err
		}

		err = c.Send(output.Bytes())
		if err != nil {
			return err
		}

		return nil
	}
}

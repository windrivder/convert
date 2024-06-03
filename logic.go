package main

import (
	"context"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type Logic struct {
	cfg    *Config
	logger echo.Logger
}

func NewLogic(cfg *Config, logger echo.Logger) *Logic {
	return &Logic{cfg: cfg, logger: logger}
}

func (l *Logic) OfficeConvert(ctx context.Context, req *OfficeConvertRequest) (*OfficeConvertReply, error) {
	path := filepath.Join(l.cfg.OutputDir, req.Input)

	// 调用 gotenberg 接口
	client := resty.New()
	client.SetOutputDirectory(filepath.Base(path))
	resp, err := client.R().
		SetFile("file", path).
		SetOutput("index.pdf").
		Post(l.cfg.GotenbergAddr + "/forms/libreoffice/convert")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, echo.ErrBadRequest
	}

	// 读取文件
	doc, err := fitz.New(filepath.Join(filepath.Base(path), "index.pdf"))
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	// Extract pages as images
	outputs := []string{}
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			l.logger.Error(err)
			continue
		}

		name := filepath.Join(filepath.Base(path), fmt.Sprintf("%04d.jpg", n))
		f, err := os.Create(name)
		if err != nil {
			l.logger.Error(err)
			continue
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
		if err != nil {
			l.logger.Error(err)
			continue
		}

		output, _ := filepath.Rel(l.cfg.OutputDir, name)
		outputs = append(outputs, output)
		f.Close()
	}
	return &OfficeConvertReply{Outputs: outputs}, nil
}

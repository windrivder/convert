package main

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
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
	pdfPath := filepath.Join(filepath.Dir(path), "index.pdf")
	if err := l.officeToPDF(path, pdfPath); err != nil {
		return nil, err
	}

	// 读取文件
	doc, err := fitz.New(pdfPath)
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

		name := filepath.Join(filepath.Dir(path), fmt.Sprintf("%04d.jpg", n))
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

func (l *Logic) officeToPDF(filePath, outputFilePath string) error {
	log.Printf("input: %s\n", filePath)
	log.Printf("output: %s\n", outputFilePath)

	// 打开文件用于读取
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %v\n", err)
		return err
	}
	defer file.Close()

	// 获取文件大小
	fileStat, err := file.Stat()
	if err != nil {
		log.Printf("Failed to get file stats: %v\n", err)
		return err
	}
	fileSize := fileStat.Size()

	// 读取文件内容到字节数组
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		log.Printf("Failed to read file content: %v\n", err)
		return err
	}
	// 创建HTTP POST请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", filePath)
	if err != nil {
		log.Printf("Failed to create form file: %v\n", err)
		return err
	}
	_, err = io.Copy(part, bytes.NewReader(buffer))
	if err != nil {
		log.Printf("Failed to copy file to body: %v\n", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Printf("Failed to close writer: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", l.cfg.GotenbergAddr+"/forms/libreoffice/convert", body)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("Request failed with status: %s, response: %s\n", resp.Status, string(bodyBytes))
		log.Println(err.Error())
		return err
	}

	// 将响应体写入到输出文件
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("Failed to create output file: %v\n", err)
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		log.Printf("Failed to write response to file: %v\n", err)
		return err
	}

	fmt.Println("Conversion successful. PDF saved as", outputFilePath)
	return nil
}

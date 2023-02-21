// Copyright 2023 The imgctl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/cosmtrek/imgctl"
)

func main() {
	_ = godotenv.Load()

	app := &cli.App{
		Name:    "imgcli",
		Usage:   "control images with ease",
		Version: "0.1.0",
		Commands: []*cli.Command{
			{
				Name:    "ocr",
				Aliases: []string{"o"},
				Usage:   "Extract text from image using OCR service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "image",
						Aliases: []string{"i"},
						Value:   "",
						Usage:   "Image file path, if not set, read from clipboard",
					},
					&cli.StringFlag{
						Name:    "service",
						Aliases: []string{"s"},
						Value:   imgctl.TencentOCRKey,
						Usage:   "OCR service",
					},
				},
				Action: func(cCtx *cli.Context) error {
					image := cCtx.String("image")
					srv := cCtx.String("service")
					var ocrServ imgctl.OCRServ
					switch srv {
					case imgctl.TencentOCRKey:
						secretID := os.Getenv("TENCENT_OCR_SECRET_ID")
						secretKey := os.Getenv("TENCENT_OCR_SECRET_KEY")
						if secretID == "" || secretKey == "" {
							panic("TENCENT_OCR_SECRET_ID or TENCENT_OCR_SECRET_KEY is empty")
						}
						ocrServ = imgctl.NewTencentOCR(secretID, secretKey)
					default:
						return errors.New("unsupported OCR engine")
					}
					return imgctl.NewEngine(ocrServ).Run(imgctl.EngineCtx{
						Ctx:     context.Background(),
						ImgPath: image,
					})
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

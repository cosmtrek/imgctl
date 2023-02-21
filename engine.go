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

package imgctl

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

type Engine struct {
	ocrServ OCRServ
}

type EngineCtx struct {
	Ctx     context.Context
	ImgPath string
}

func NewEngine(ocrServ OCRServ) *Engine {
	return &Engine{
		ocrServ: ocrServ,
	}
}

func (e *Engine) Run(eCtx EngineCtx) error {
	var err error
	var imgBuf []byte
	if eCtx.ImgPath != "" {
		if imgBuf, err = os.ReadFile(eCtx.ImgPath); err != nil {
			return err
		}
	} else {
		if err = clipboard.Init(); err != nil {
			return err
		}
		imgBuf = clipboard.Read(clipboard.FmtImage)
		if len(imgBuf) == 0 {
			fmt.Println("clipboard is empty")
			return nil
		}
	}

	imgBs64 := base64.StdEncoding.EncodeToString(imgBuf)
	ocrReq := &OCRReq{
		ImgBs64: imgBs64,
	}
	ocrResp, err := e.ocrServ.Accept(eCtx.Ctx, ocrReq)
	if err != nil {
		return err
	}
	fmt.Printf(`id: %s, lang: %s
content: %s`, ocrResp.ReqID, ocrResp.Lang, ocrResp.Content)
	return nil
}

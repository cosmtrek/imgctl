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
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	terrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

type OCRReq struct {
	ImgBs64 string
}

type OCRResp struct {
	ReqID   string
	Lang    string
	Content string
}

type OCRServ interface {
	Accept(ctx context.Context, req *OCRReq) (*OCRResp, error)
}

const (
	TencentOCRKey = "tencent"
)

type tencentOCR struct {
	client *ocr.Client
}

func NewTencentOCR(key string, secret string) OCRServ {
	credential := common.NewCredential(key, secret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	client, _ := ocr.NewClient(credential, "ap-beijing", cpf)
	return &tencentOCR{
		client: client,
	}
}

func (o *tencentOCR) Accept(ctx context.Context, req *OCRReq) (*OCRResp, error) {
	if o.client == nil {
		return nil, errors.New("tencent ocr client is nil")
	}
	oreq := ocr.NewGeneralFastOCRRequest()
	oreq.ImageBase64 = &req.ImgBs64
	resp, err := o.client.GeneralFastOCR(oreq)
	if _, ok := err.(*terrors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	r := resp.Response
	if r == nil {
		return nil, errors.New("tencent ocr response is nil")
	}
	oresp := &OCRResp{
		ReqID: *r.RequestId,
		Lang:  *r.Language,
	}
	strs := make([]string, 0, len(r.TextDetections))
	lo.ForEach(r.TextDetections, func(v *ocr.TextDetection, i int) {
		strs = append(strs, *v.DetectedText)
	})
	oresp.Content = strings.Join(strs, "\n")
	return oresp, nil
}

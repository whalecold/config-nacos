// Copyright 2023 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex-examples/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/kitex_gen/api/echo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	nacosclient "github.com/kitex-contrib/config-nacos/client"
	"github.com/kitex-contrib/config-nacos/nacos"
	"github.com/kitex-contrib/config-nacos/utils"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type configLog struct{}

func (cl *configLog) Apply(opt *utils.Options) {
	fn := func(cp *vo.ConfigParam) {
		klog.Infof("nacos config %v", cp)
	}
	opt.NacosCustomFunctions = append(opt.NacosCustomFunctions, fn)
}

func main() {
	klog.SetLevel(klog.LevelDebug)

	nacosClient, err := nacos.NewClient(nacos.Options{})
	if err != nil {
		panic(err)
	}

	cl := &configLog{}

	serviceName := "ServiceName"
	clientName := "ClientName"
	client, err := echo.NewClient(
		serviceName,
		client.WithHostPorts("0.0.0.0:8888"),
		client.WithSuite(nacosclient.NewSuite(serviceName, clientName, nacosClient, cl)),
	)
	if err != nil {
		log.Fatal(err)
	}
	for {
		req := &api.Request{Message: "my request"}
		resp, err := client.Echo(context.Background(), req)
		if err != nil {
			klog.Errorf("take request error: %v", err)
		} else {
			klog.Infof("receive response %v", resp)
		}
		time.Sleep(time.Second * 10)
	}
}

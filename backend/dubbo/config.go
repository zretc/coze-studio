/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dubbo

import (
	"context"
	"fmt"
	"log"
	"os"

	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// InitDubbo 初始化Dubbo服务
func InitDubbo(ctx context.Context) error {
	// 从环境变量获取Nacos配置
	nacosAddr := getEnv("NACOS_ADDR", "localhost:8848")
	nacosNamespace := getEnv("NACOS_NAMESPACE", "public")
	nacosUsername := getEnv("NACOS_USERNAME", "nacos")
	nacosPassword := getEnv("NACOS_PASSWORD", "nacos")

	log.Println("[Dubbo] 开始初始化Dubbo服务...")
	log.Printf("[Dubbo] Nacos配置: addr=%s, namespace=%s, username=%s", nacosAddr, nacosNamespace, nacosUsername)

	// 配置Nacos客户端
	log.Println("[Dubbo] 正在初始化Nacos客户端...")
	_, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig: &nacosConstant.ClientConfig{
				NamespaceId:         nacosNamespace,
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				LogDir:              "logs/nacos",
				CacheDir:            "cache/nacos",
				LogLevel:            "info",
				Username:            nacosUsername,
				Password:            nacosPassword,
			},
			ServerConfigs: []nacosConstant.ServerConfig{
				{
					IpAddr:      nacosAddr,
					ContextPath: "/nacos",
					Port:        8848,
				},
			},
		})
	if err != nil {
		log.Printf("[Dubbo] 初始化Nacos客户端失败: %v", err)
		return fmt.Errorf("init nacos client failed: %w", err)
	}
	log.Println("[Dubbo] Nacos客户端初始化成功")

	// 配置Dubbo
	log.Println("[Dubbo] 配置Dubbo服务...")
	config.SetProviderService(new(CozeService))
	log.Println("[Dubbo] 注册CozeService服务")

	// 初始化Dubbo
	log.Println("[Dubbo] 加载Dubbo配置...")
	if err := config.Load(); err != nil {
		log.Printf("[Dubbo] 加载Dubbo配置失败: %v", err)
		return fmt.Errorf("load dubbo config failed: %w", err)
	}
	log.Println("[Dubbo] Dubbo配置加载成功")

	log.Println("[Dubbo] Dubbo服务初始化完成")
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

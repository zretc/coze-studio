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
	"os"
	"time"

	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/registry"
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
	nacosAuthToken := getEnv("NACOS_AUTH_TOKEN", "Y296ZS1zdHVkaW8tbmFjb3MtdG9rZW4tMTc3NTY3ODkwNw==")
	nacosAuthIdentityKey := getEnv("NACOS_AUTH_IDENTITY_KEY", "coze-studio")
	nacosAuthIdentityValue := getEnv("NACOS_AUTH_IDENTITY_VALUE", "coze-studio-secret")

	// 配置Nacos客户端
	nacosClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig: &nacosConstant.ClientConfig{
				NamespaceId:         nacosNamespace,
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				LogDir:              "logs/nacos",
				CacheDir:            "cache/nacos",
				RotateTime:          "1h",
				MaxAge:              3,
				LogLevel:            "info",
				Username:            nacosUsername,
				Password:            nacosPassword,
				IdentityKey:         nacosAuthIdentityKey,
				IdentityValue:       nacosAuthIdentityValue,
			},
			ServerConfigs: []nacosConstant.ServerConfig{
				{
					IpAddr:      nacosAddr,
					ContextPath: "/nacos",
					Port:        8848,
					GrpcPort:    9848,
				},
			},
		})
	if err != nil {
		return fmt.Errorf("init nacos client failed: %w", err)
	}

	// 配置Dubbo
	config.SetProviderService(new(CozeService))

	// 设置应用配置
	applicationConfig := config.NewApplicationConfig(
		config.WithApplicationName("coze-studio-backend"),
	)

	// 设置注册中心配置
	registryConfig := config.NewRegistryConfig(
		config.WithProtocol("nacos"),
		config.WithAddress(nacosAddr),
		config.WithUsername(nacosUsername),
		config.WithPassword(nacosPassword),
		config.WithNamespace(nacosNamespace),
		config.WithParams(map[string]string{
			"nacos.auth.token":          nacosAuthToken,
			"nacos.auth.identity.key":   nacosAuthIdentityKey,
			"nacos.auth.identity.value": nacosAuthIdentityValue,
		}),
	)

	// 设置协议配置
	protocolConfig := config.NewProtocolConfig(
		config.WithName("dubbo"),
		config.WithPort(20000),
		config.WithParams(map[string]string{
			"threadpool": "fixed",
			"threads":    "200",
			"queues":     "1000",
		}),
	)

	// 初始化Dubbo
	if err := config.Load(
		config.WithApplicationConfig(applicationConfig),
		config.WithRegistryConfig(registryConfig),
		config.WithProtocolConfig(protocolConfig),
	); err != nil {
		return fmt.Errorf("load dubbo config failed: %w", err)
	}

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

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

	log.Println("[Dubbo] 开始初始化Dubbo服务...")
	log.Printf("[Dubbo] Nacos配置: addr=%s, namespace=%s, username=%s", nacosAddr, nacosNamespace, nacosUsername)

	// 配置Nacos客户端
	log.Println("[Dubbo] 正在初始化Nacos客户端...")
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
		log.Printf("[Dubbo] 初始化Nacos客户端失败: %v", err)
		return fmt.Errorf("init nacos client failed: %w", err)
	}
	log.Println("[Dubbo] Nacos客户端初始化成功")

	// 测试Nacos连接
	log.Println("[Dubbo] 测试Nacos连接...")
	services, err := nacosClient.GetServices(vo.GetServicesParam{
		PageNo:   1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("[Dubbo] Nacos连接测试失败: %v", err)
	} else {
		log.Printf("[Dubbo] Nacos连接测试成功，当前服务数量: %d", len(services.Doms))
	}

	// 配置Dubbo
	log.Println("[Dubbo] 配置Dubbo服务...")
	config.SetProviderService(new(CozeService))
	log.Println("[Dubbo] 注册CozeService服务")

	// 设置应用配置
	applicationConfig := config.NewApplicationConfig(
		config.WithApplicationName("coze-studio-backend"),
	)
	log.Println("[Dubbo] 设置应用配置: coze-studio-backend")

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
	log.Printf("[Dubbo] 设置注册中心配置: %s", nacosAddr)

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
	log.Println("[Dubbo] 设置协议配置: dubbo:20000")

	// 初始化Dubbo
	log.Println("[Dubbo] 加载Dubbo配置...")
	if err := config.Load(
		config.WithApplicationConfig(applicationConfig),
		config.WithRegistryConfig(registryConfig),
		config.WithProtocolConfig(protocolConfig),
	); err != nil {
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

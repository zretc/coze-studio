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

	"dubbo.apache.org/dubbo-go/v3/config"
)

// CozeService 定义Dubbo服务接口
type CozeService struct {
}

// Reference 服务引用
func (s *CozeService) Reference() string {
	return "coze.studio.CozeService"
}

// SayHello 示例方法
func (s *CozeService) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello, " + name, nil
}

// GetAppInfo 获取应用信息
func (s *CozeService) GetAppInfo(ctx context.Context, appID string) (map[string]interface{}, error) {
	// 这里可以实现获取应用信息的逻辑
	return map[string]interface{}{
		"appId":  appID,
		"status": "active",
		"name":   "Coze Studio App",
	}, nil
}

// RegisterConsumer 注册Dubbo消费者
func RegisterConsumer() error {
	// 注册Java服务消费者
	_, err := GetJavaService()
	if err != nil {
		return err
	}
	return nil
}

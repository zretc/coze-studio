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

// JavaService Java服务接口
type JavaService struct {
	SayHello func(ctx context.Context, name string) (string, error)
	GetData  func(ctx context.Context, id string) (map[string]interface{}, error)
}

// Reference 服务引用
func (s *JavaService) Reference() string {
	return "com.example.JavaService"
}

// GetJavaService 获取Java服务实例
func GetJavaService() (*JavaService, error) {
	var javaService JavaService
	if err := config.SetConsumerService(&javaService); err != nil {
		return nil, err
	}
	return &javaService, nil
}

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

package coze

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/dubbo"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// TestDubboGetUserById 测试调用Java Dubbo服务的getUserById方法
func TestDubboGetUserById(ctx context.Context, c *app.RequestContext) {
	// 获取userId参数
	userIdStr := c.Query("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "userId is required",
		})
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid userId",
		})
		return
	}

	// 获取Dubbo用户服务实例
	userService, err := dubbo.GetDubboUserService()
	if err != nil {
		logs.Errorf("GetDubboUserService failed: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get dubbo user service",
		})
		return
	}

	// 调用getUserById方法
	user, err := userService.GetUserById(ctx, userId)
	if err != nil {
		logs.Errorf("GetUserById failed: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

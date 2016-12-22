/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package register

import (
	"github.com/coreos/etcd/client"
	"time"
)

type Register interface {

	/**
	建立链接
	 */
	connect()

	/**
	设置元素
	 */
	Set(path string,value string) error

	/**
	设置元素 包含有效时间
	 */
	SetWithTime(path string,value string,times time.Duration) error

	/**
	获得元素
	 */
	Get(path string) (Node,error)

	/**
	获得元素
	 */
	GetChildren(path string) ([]Node,error)

	/**
	删除元素
	 */
	Delete(path string) error

	/**
	订阅服务变化
	 */
	Subscribe(path string  , cancel <- chan int, handler func(*client.Response))

	/**
	维持心跳  etcd需要自己维持心跳保持临时节点有效
	 */
	TimeTicker()

}

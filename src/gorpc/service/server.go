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
package service

/**
服务
 */
type Service struct {
	Servic		interface{} 	`注册服务`
	Protocol	string		`协议类型`
}

/**
服务节点机器信息
 */
type Provider struct {
	Address string 	`服务地址`
	Invoke 	int	`调用次数`
}

func (p *Provider) Invok()  {
	p.Invoke = p.Invoke + 1
}

/**
服务消费者
 */
type Consumer struct {
	Host	string	`消费机器地址`
	Name	string	`服务名称`
}





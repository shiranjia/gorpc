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
package monitor

import (
	"gorpc/register"
	"gorpc/utils"
	"log"
)

type Monitor struct {
	Register	register.Register		`注册中心`
	Service		map[string]MonitorService	`服务列表`
}

func (m *Monitor) GetDate()  {
	services,err := m.Register.GetChildren(utils.Separator)
	utils.CheckErr("monitor.GetDate",err)
	log.Println("service:",services)
	for _,s := range services{
		service := MonitorService{}
		service.Name = s.Key
		log.Println("servicePath:",s.Path)
		provs,err := m.Register.GetChildren(s.Key + utils.Separator + "provider")
		utils.CheckErr("monitor.GetProviders",err)
		providers := make([]string,10)
		for _ ,p := range provs {
			providers = append(providers,p.Key)
		}
		service.Provider = providers
		cons,err := m.Register.GetChildren(s.Key + utils.Separator + "consumer")
		utils.CheckErr("monitor.GetConsumers",err)
		consumers := make([]string,10)
		for _,c := range cons{
			consumers = append(consumers,c.Key)
		}
		m.Service[s.Key] = service
	}
}

/**
服务
 */
type MonitorService struct {
	Name		string 			`注册服务`
	Protocol	string			`协议类型`
	Provider	[]string		`服务提供者`
	Consumer	[]string		`服务消费者`
}


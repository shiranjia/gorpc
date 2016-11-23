package service

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





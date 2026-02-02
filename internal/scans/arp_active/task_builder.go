package arp_active

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol{

	return domain.ARP_ACTIVE
}


func (b Builder) Build(cidr string,cfg domain.Config)[]tasks.Task{

    fmt.Println("[tasks] requested protocol:", cfg.Type)
 

	var tasks []tasks.Task

	tasks = append(tasks, &Task{timeout: cfg.Timeout,cidr:cidr})

	return tasks
	
}


func init() {
	tasks.Register(Builder{})
}

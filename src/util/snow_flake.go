package util

import (
	"fmt"
	"github.com/holdno/snowFlakeByGo"
)

var worker *snowFlakeByGo.Worker

func init() {
	// 生成一个节点实例
	worker, _ = snowFlakeByGo.NewWorker(0) // 传入当前节点id 此id在机器集群中一定要唯一 且从0开始排最多1024个节点，可以根据节点的不同动态调整该算法每毫秒生成的id上限(如何调整会在后面讲到)
}

func GenerateId() string {
	return fmt.Sprintf("%v", worker.GetId())
}

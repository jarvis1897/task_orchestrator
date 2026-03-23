package master

import (
	"sync"
	"time"
	"github.com/jarvis1897/task_orchestrator/internal/common"
)

type Registry struct {
	nodes map[string]*common.Node
	mu sync.Mutex
}

func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[string]*common.Node),
	}
}

func (r *Registry) Register(node *common.Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[node.ID] = node
}

func (r *Registry) Heartbeat(nodeId string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if node, exists := r.nodes[nodeId]; exists {
		node.LastHeartbeat = time.Now()
		node.Status = "active"
	}
}

func (r *Registry) GetActiveNodes() []*common.Node {
	r.mu.Lock()
	defer r.mu.Unlock()
	active := make([]*common.Node, 0)
	for _, node := range r.nodes {
		if node.Status == "active" {
			active = append(active, node)
		}
	}
	return active
}

func (r *Registry) StartDeadNodeChecker() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			r.mu.Lock()
			for _, node := range r.nodes {
				if time.Since(node.LastHeartbeat) > 30*time.Second {
					node.Status = "dead"
				}
			}
			r.mu.Unlock()
		}
	}()
}


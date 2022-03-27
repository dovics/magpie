package scheduler

type Cache struct {
	nodes map[string]*nodeInfo
	tasks map[string]*taskInfo

	headNode *nodeInfo
}

type nodeInfo struct {
	name string

	tasks []*taskInfo

	next *nodeInfo
	prev *nodeInfo

	allocatable *Resource
	requested   *Resource
}

type TaskStatus int

const (
	Pending TaskStatus = iota
	Running
	Completed
	Failed
)

type taskInfo struct {
	name      string
	Requested *Resource

	status TaskStatus
	node   *nodeInfo
}

func cacheCleaner() {
	// todo: periodically clean up completed tasks
	for {

	}
}

func (c *Cache) AddNode(node *WorkerConfig) {
	if _, ok := c.nodes[node.name]; ok {
		c.DeleteNode(node.name)
	}

	nodeInfo := &nodeInfo{name: node.name, allocatable: node.resources, requested: &Resource{0, 0}, prev: c.headNode, next: c.headNode.next}

	c.headNode.next = nodeInfo
	nodeInfo.next.prev = nodeInfo
	c.nodes[node.name] = nodeInfo
}

func (c *Cache) DeleteNode(node string) {
	n, ok := c.nodes[node]
	if !ok {
		return
	}

	n.prev.next = n.next
	n.next.prev = n.prev
	delete(c.nodes, node)
}

func (c *Cache) AddTask(node string, task Task) {
	name := task.Name()
	n := c.nodes[node]

	t := &taskInfo{
		name:      name,
		Requested: task.Resource(),
		node:      n,
	}

	c.tasks[name] = t
	n.tasks = append(n.tasks, t)
}

func (c *Cache) UpdateTaskStatus(task string, status TaskStatus) {
	c.tasks[task].status = status
}

package scheduler

// Resource is a collection of compute resource.
type Resource struct {
	MilliCPU int64
	Memory   int64
}

func (r *Resource) Add(req *Resource) *Resource {
	return &Resource{
		r.MilliCPU + req.MilliCPU,
		r.Memory + req.Memory,
	}
}

func (r *Resource) Sub(req *Resource) *Resource {
	return &Resource{
		r.MilliCPU - req.MilliCPU,
		r.Memory - req.Memory,
	}
}

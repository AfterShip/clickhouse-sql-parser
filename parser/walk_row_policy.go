package parser

func walkRowPolicyOptions(options *RowPolicyOptions, rename *Ident, fn WalkFunc) bool {
	for _, target := range options.Targets {
		if !Walk(target, fn) {
			return false
		}
	}
	if !Walk(options.AccessStorageType, fn) || !Walk(options.OnCluster, fn) || !Walk(rename, fn) {
		return false
	}
	for _, filter := range options.Filters {
		if !Walk(filter, fn) {
			return false
		}
	}
	return Walk(options.To, fn)
}

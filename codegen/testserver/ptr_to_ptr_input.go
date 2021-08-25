package testserver

type PtrToPtrOuter struct {
	Name        string
	Inner       *PtrToPtrInner
	StupidInner *******PtrToPtrInner
}

type PtrToPtrInner struct {
	Key   string
	Value string
}

type UpdatePtrToPtrOuter struct {
	Name        *string
	Inner       **UpdatePtrToPtrInner
	StupidInner ********PtrToPtrInner
}

type UpdatePtrToPtrInner struct {
	Key   *string
	Value *string
}

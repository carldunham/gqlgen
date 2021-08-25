package testserver

type PtrToPtrOuter struct {
	Name  string
	Inner *PtrToPtrInner
}

type PtrToPtrInner struct {
	Key   string
	Value string
}

type UpdatePtrToPtrOuter struct {
	Name  *string
	Inner **UpdatePtrToPtrInner
}

type UpdatePtrToPtrInner struct {
	Key   *string
	Value *string
}

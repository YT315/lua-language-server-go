package analysis

import "sync"

//类型集合
type TypeSet struct {
	Count  int
	Types  []TypeInfo
	posMap map[TypeInfo]int
	mu     sync.Mutex
}

func NewTypeSet() *TypeSet {
	res := &TypeSet{
		Count:  0,
		Types:  nil,
		posMap: map[TypeInfo]int{},
	}
	return res
}
func NewTypeSetWithCap(cap int) *TypeSet {
	res := &TypeSet{
		Count:  0,
		Types:  make([]TypeInfo, cap),
		posMap: make(map[TypeInfo]int, cap),
	}
	return res
}

func NewTypeSetWithContent(tp ...TypeInfo) *TypeSet {
	cap := cap(tp)
	res := &TypeSet{
		Count:  0,
		Types:  make([]TypeInfo, cap),
		posMap: make(map[TypeInfo]int, cap),
	}
	res.AddRange(tp...)
	return res
}

//添加内容
func (ts *TypeSet) Add(tp TypeInfo) (res bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if _, ok := ts.posMap[tp]; !ok {
		ts.Types = append(ts.Types, tp)
		ts.posMap[tp] = ts.Count
		ts.Count++
		res = true
	} else {
		res = false
	}

	return
}

//批量添加内容
func (ts *TypeSet) AddRange(tp ...TypeInfo) (res []bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	res = make([]bool, len(tp))
	for idx, val := range tp {
		if _, ok := ts.posMap[val]; !ok {
			ts.Types = append(ts.Types, val)
			ts.posMap[val] = ts.Count
			res[idx] = true
			ts.Count++
		} else {
			res[idx] = false
		}
	}
	return
}

//删除内容
func (ts *TypeSet) Del(tp TypeInfo) (res bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if val, ok := ts.posMap[tp]; ok {
		delete(ts.posMap, tp)
		ts.Types = append(ts.Types[:val], ts.Types[val+1:]...)
		ts.Count--
		res = true
	} else {
		res = false
	}
	return
}

//是否包含
func (ts *TypeSet) Contain(tp TypeInfo) (res bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	_, res = ts.posMap[tp]
	return
}

//索引
func (ts *TypeSet) IndexOf(tp TypeInfo) (res int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if val, ok := ts.posMap[tp]; !ok {
		res = -1
	} else {
		res = val
	}
	return
}

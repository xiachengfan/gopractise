package main

import (
	"math/rand"
)

const (
	maxLevel    = 32
	probability = 0.25
)

type (
	SortedSet struct {
		record map[string]*SortedSetNode
	}
	SortedSetNode struct {
		dict map[string]*sklNode
		skl  *skipList
	}
	//节点的level数组，保存每层的后向指针和跨度
	sklLevel struct {
		forward *sklNode
		span    uint64 //这是用来记录结点在某一层上的*forward指针和该指针指向的结点之间，跨越了 level0 上的几个结点。
	}
	sklNode struct {
		member   string
		score    float64
		backward *sklNode
		level    []*sklLevel
	}
	skipList struct {
		head *sklNode
		//尾节点
		tail   *sklNode
		length int64
		level  int16
	}
)

func New() *SortedSet {
	return &SortedSet{
		make(map[string]*SortedSetNode),
	}
}

func (z *SortedSet) exist(key string) bool {
	_, exist := z.record[key]
	return exist
}

func (z *SortedSet) ZAdd(key string, score float64, member string) {
	//判断有序集合的key是否存在，不存在，则新建
	if !z.exist(key) {
		node := &SortedSetNode{
			dict: make(map[string]*sklNode),
			skl:  newSkipList(),
		}
		z.record[key] = node
	}
	item := z.record[key]
	//根据加入的member，先进行hash更改
	v, ok := item.dict[member]
	//进行调表更改
	var node *sklNode
	if ok {
		if score != v.score {
			item.skl.sklDelete(v.score, member)
			node = item.skl.sklInsert(score, member)
		}
	} else {
		node = item.skl.sklInsert(score, member)
	}
	if node != nil {
		item.dict[member] = node

	}
}

func (z *SortedSet) ZScore(key string, member string) (ok bool, score float64) {
	if z.exist(key) {
		return
	}

	node, exist := z.record[key].dict[member]
	if !exist {
		return
	}
	return true, node.score
}
func (z *SortedSet) ZCard(key string) int {
	if !z.exist(key) {
		return 0
	}
	return len(z.record[key].dict)
}
func (z *SortedSet) ZRank(key, member string) int64 {
	if !z.exist(key) {
		return -1
	}
	node, exist := z.record[key].dict[member]
	if !exist {
		return -1
	}
	rank := z.record[key].skl.sklGetRank(node.score, member)
	//减去头结点的1
	rank--
	return rank
}

func (z *SortedSet) ZRevRank(key, member string) int64 {
	if !z.exist(key) {
		return -1
	}
	node, exist := z.record[key].dict[member]
	if !exist {
		return -1
	}
	rank := z.record[key].skl.sklGetRank(node.score, member)
	return z.record[key].skl.length - rank
}
func sklNewNode(level int16, score float64, member string) *sklNode {
	node := &sklNode{
		score:  score,
		member: member,
		level:  make([]*sklLevel, level),
	}
	for i := range node.level {
		node.level[i] = new(sklLevel)
	}
	return node
}

func newSkipList() *skipList {
	return &skipList{
		level: 1,
		head:  sklNewNode(maxLevel, 0, ""),
	}
}

func randomLevel() int16 {
	var level int16 = 1
	for float32(rand.Int31()&0xFFFF) < (probability * 0xFFFF) {
		level++
	}
	if level < maxLevel {
		return level
	}

	return maxLevel
}

func (skl *skipList) sklInsert(score float64, member string) *sklNode {
	//①查找要插入的位置；②调整跳跃表高度；③插入节点；④调整backward。
	//插入节点时，需要更新被插入节点每层的前一个节点。由于每层更新的节点不一样，所以将每层需要更新的节点记录在update[i]中。
	updates := make([]*sklNode, maxLevel)
	//记录当前层从header节点到update[i]节点所经历的步长，在更新update[i]的span和设置新插入节点的span时用到。
	rank := make([]uint64, maxLevel)
	//调表的增删查都是依靠前节点，每一层的前节点不一致
	cur := skl.head
	for i := skl.level - 1; i >= 0; i-- {
		//记录最高层，而且是头结点。
		if i == skl.level-1 {
			rank[i] = 0
		} else {
			//不是最高层的时候都需要先加上一层所走的步长
			rank[i] = rank[i+1]
		}
		if cur.level[i] != nil {
			for cur.level[i].forward != nil &&
				(cur.level[i].forward.score < score ||
					(cur.level[i].forward.score == score && cur.level[i].forward.member < member)) {
				//加上当前节点的步长
				rank[i] += cur.level[i].span
				//前节点后移
				cur = cur.level[i].forward
			}
		}
		updates[i] = cur
	}
	//随机一个level，保持均衡
	level := randomLevel()
	//如果层高大于了原来的，需要修改头结点中增加的层数。
	if level > skl.level {
		//大于level的值，结点一定是头结点。头节点span是整个调表的长度。
		for i := skl.level; i < level; i++ {
			rank[i] = 0
			updates[i] = skl.head
			updates[i].level[i].span = uint64(skl.length)
		}
		skl.level = level
	}
	p := sklNewNode(level, score, member)
	//更改sklLevel
	for i := int16(0); i < level; i++ {
		//修改每一层的sklLevel中的forward值与新加节点的sklLevel中的forward赋值
		p.level[i].forward = updates[i].level[i].forward
		updates[i].level[i].forward = p
		//修改每一层的sklLevel中的span值与新加节点的sklLevel中的span赋值
		//设置加入节点的span
		p.level[i].span = updates[i].level[i].span - (rank[0] - rank[i])
		//设置上一层节点的span
		updates[i].level[i].span = (rank[0] - rank[i]) + 1
	}
	//level值小于skl.level时，所有层高之上的层只需要步长加一
	for i := level; i < skl.level; i++ {
		updates[i].level[i].span++
	}
	//头节点在有序集合中不存储任何member和score值，member值为""，
	//score值为0；所以指向头节点的backward的指针为NULL。其实在go语言中字符串不能为NULL，可以为""。
	if updates[0] == skl.head {
		p.backward = nil
	} else {
		//指向updates[0]
		p.backward = updates[0]
	}
	//调整后一节点的后向指正。如果新插入节点是最后一个节点，则需要更新跳跃表的尾节点为新插入节点。插入节点后，更新跳跃表的长度加1。
	if p.level[0].forward != nil {
		p.level[0].forward.backward = p
	} else {
		skl.tail = p
	}
	skl.length++
	return p
}

func (s *skipList) sklDelete(score float64, member string) {
	//删除节点的步骤：1）查找需要更新的节点；2）设置span和forward。
	if "" == member {
		return
	}
	cur := s.head
	update := make([]*sklNode, maxLevel)
	for i := s.level - 1; i >= 0; i-- {
		//寻找前置节点，重点。因为查找，删除，插入都会用到这个。
		for cur.level[i].forward != nil && (cur.level[i].forward.score < score ||
			(cur.level[i].forward.score == score && cur.level[i].forward.member < member)) { //当查找到的结点保存的元素score，等于要查找的score时，跳表会再检查该结点保存的 member 类型数据，是否比要查找的 member 数据小。
			cur = cur.level[i].forward

		}
		//保存每一层的前置节点。因为每一层的前置节点都不一样。
		update[i] = cur
	}
	//找到删除节点
	cur = cur.level[0].forward
	if cur != nil && score == cur.score && cur.member == member {
		s.sklDeleteNode(cur, update)
		return
	}

}

func (s *skipList) sklDeleteNode(p *sklNode, updates []*sklNode) {
	//计算每一层的当前节点的sklLevel的前节点的的关系，并调整前项指针。
	for i := int16(0); i < s.level; i++ {
		if updates[i].level[i].forward == p {
			updates[i].level[i].span += p.level[i].span - 1
			updates[i].level[i].forward = p.level[i].forward
		} else {
			updates[i].level[i].span--
		}
	}
	//删除当前数据节点
	if p.level[0].forward != nil {
		p.level[0].forward.backward = p.backward
	} else {
		s.tail = p.backward
	}
	//删除的层高是最高的，并且没有其他节点与p节点的高度相同时。高度降低后，头节点的内存并没有释放，增加高度会更新，这里申请的内存就不用释放。
	for s.level > 1 && s.head.level[s.level-1].forward == nil {
		s.level--
	}

	s.length--

}

func (s *skipList) sklGetRank(score float64, member string) int64 {
	var rank uint64 = 0
	cur := s.head

	//寻找前置节点，重点。因为查找，删除，插入都会用到这个。
	for i := s.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && (cur.level[i].forward.score < score ||
			(cur.level[i].forward.score == score && cur.level[i].forward.member < member)) { //当查找到的结点保存的元素score，等于要查找的score时，跳表会再检查该结点保存的 member 类型数据，是否比要查找的 member 数据小。
			rank += cur.level[i].span
			cur = cur.level[i].forward

		}
		if cur.member == member {
			return int64(rank)
		}
	}

	return 0
}

func (s *skipList) sklGetElementByRank(rank uint64) *sklNode {
	var traversed uint64 = 0
	cur := s.head
	for i := s.level - 1; i >= 1; i-- {
		for cur.level[i].forward != nil && (traversed+cur.level[i].span) <= rank {
			traversed += cur.level[i].span
			cur = cur.level[i].forward
		}
		if traversed == rank {
			return cur
		}
	}
	return nil
}

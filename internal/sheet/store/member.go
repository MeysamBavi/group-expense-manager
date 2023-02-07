package store

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"strings"
)

type MemberStore struct {
	indexByName   map[string]int
	memberByIndex []*model.Member
}

func NewMemberStore() *MemberStore {
	return &MemberStore{
		indexByName:   make(map[string]int),
		memberByIndex: make([]*model.Member, 0),
	}
}

func (ms *MemberStore) Count() int {
	return len(ms.memberByIndex)
}

func (ms *MemberStore) AddMember(member *model.Member) error {
	name := standardizeName(member.Name)
	if name == "" {
		return fmt.Errorf("empty or whitespace name")
	}

	_, found := ms.indexByName[name]
	if found {
		return fmt.Errorf("duplicate member name %q", member.Name)
	}

	ms.indexByName[name] = len(ms.memberByIndex)
	ms.memberByIndex = append(ms.memberByIndex, member)
	return nil
}

func (ms *MemberStore) GetMemberByName(name string) (*model.Member, bool) {
	name = standardizeName(name)
	i, ok := ms.indexByName[name]
	if !ok {
		return nil, false
	}

	return ms.GetMemberByIndex(i)
}

func (ms *MemberStore) GetMemberByIndex(index int) (*model.Member, bool) {
	if index < 0 || index >= len(ms.memberByIndex) {
		return nil, false
	}

	return ms.memberByIndex[index], true
}

func (ms *MemberStore) GetIndexByName(name string) int {
	name = standardizeName(name)
	i, ok := ms.indexByName[name]
	if !ok {
		return -1
	}
	return i
}

func (ms *MemberStore) RequireMemberByIndex(index int) *model.Member {
	return ms.memberByIndex[index]
}

func (ms *MemberStore) IsValid(name string, index int) bool {
	name = standardizeName(name)
	trueIndex, ok := ms.indexByName[name]
	return ok && index == trueIndex
}

func (ms *MemberStore) IsPresent(name string) bool {
	_, ok := ms.GetMemberByName(name)
	return ok
}

func (ms *MemberStore) Range(consumer func(index int, member *model.Member)) {
	for i, member := range ms.memberByIndex {
		consumer(i, member)
	}
}

func standardizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

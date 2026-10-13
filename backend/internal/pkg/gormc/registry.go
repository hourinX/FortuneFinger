package gormc

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

const MainSource = "main"

// Registry 仅在启动装配阶段写入，发布后只读；Close 必须在请求全部停止后调用。
type Registry struct {
	items map[string]*Source
}

func NewRegistry() *Registry {
	return &Registry{
		items: make(map[string]*Source),
	}
}

func (r *Registry) Main() (*gorm.DB, error) {
	return r.Resolve(MainSource)
}

func (r *Registry) Add(name string, source *Source) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("gormc: Database.Sources的键名不能为空")
	}
	if source == nil {
		return fmt.Errorf("gormc: 数据源「%s」不能为空", name)
	}
	if _, exists := r.items[name]; exists {
		return fmt.Errorf("gormc: 数据源「%s」已存在", name)
	}
	r.items[name] = source
	return nil
}

func (r *Registry) Resolve(name string) (*gorm.DB, error) {
	source, exists := r.items[name]
	if !exists {
		return nil, fmt.Errorf("gormc: 数据源「%s」不存在", name)
	}
	return source.DB(), nil
}

func (r *Registry) Close() error {
	names := make([]string, 0, len(r.items))
	for name := range r.items {
		names = append(names, name)
	}
	sort.Strings(names)

	var errs []error
	for _, name := range names {
		if err := r.items[name].Close(); err != nil {
			errs = append(errs, fmt.Errorf("gormc: 关闭数据源「%s」失败: %w", name, err))
			continue
		}
		delete(r.items, name)
	}
	return errors.Join(errs...)
}

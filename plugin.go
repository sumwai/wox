package wox

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"plugin"
	"sync"
)

type (
	// 插件接口规范
	Plug interface {
		Run()
	}

	// 插件管理器
	PlugManager struct {
		Plugins []PlugInfo // 插件列表
	}
	// 插件信息
	PlugInfo struct {
		Runner      Plug        // 入口函数
		Name        string      // 名称
		Description string      // 简介
		Version     PlugVersion // 版本
	}
	// 插件版本
	PlugVersion string
)

func NewPM() *PlugManager {
	return new(PlugManager)
}

// 解析文件到plugin
func (p PlugManager) ParseFile(file string) (info PlugInfo, err error) {
	var (
		f   *plugin.Plugin
		ok  bool
		sym plugin.Symbol

		name plugin.Symbol
		desc plugin.Symbol
		ver  plugin.Symbol
	)
	if f, err = plugin.Open(file); err != nil {
		return
	}
	if sym, err = f.Lookup("Plugin"); err != nil {
		return
	}
	if name, err = f.Lookup("Name"); err != nil {
		name = "Untitled"
	} else {
		name = *(name.(*string))
	}
	if desc, err = f.Lookup("Description"); err != nil {
		desc = ""
	} else {
		desc = *(desc.(*string))
	}
	if ver, err = f.Lookup("Version"); err != nil {
		ver = "0.0.0"
	} else {
		ver = *(ver.(*string))
	}
	info.Name = name.(string)
	info.Description = desc.(string)
	info.Version = PlugVersion(ver.(string))
	if info.Runner, ok = sym.(Plug); !ok {
		return info, errors.New("not a plugin")
	}
	return
}

// 加载插件目录下的所有插件
func (p *PlugManager) Load(s string) *PlugManager {
	err := filepath.WalkDir(s, func(file string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		var plugInfo PlugInfo
		if plugInfo, err = p.ParseFile(file); err != nil {
			log.Fatalf("load plugin %s failed: %v", d.Name(), err)
			return err
		}
		return p.Add(plugInfo)
	})
	if err != nil {
		log.Fatalf("load plugins failed: %v", err)
	}
	return p
}

// 添加插件信息到插件列表
func (p *PlugManager) Add(plug PlugInfo) (err error) {
	p.Plugins = append(p.Plugins, plug)
	log.Printf("loaded [%s] (v %s) - %s", plug.Name, plug.Version, plug.Description)
	return
}

// 运行所有插件
func (p *PlugManager) Run() {
	wg := sync.WaitGroup{}
	wg.Add(p.Count())
	for _, plug := range p.Plugins {
		go func(plug PlugInfo) {
			log.Printf("run [%s] (v %s)", plug.Name, plug.Version)
			plug.Runner.Run()
			wg.Done()
		}(plug)
	}
	wg.Wait()
}

// 获取插件数量
func (p *PlugManager) Count() int {
	return len(p.Plugins)
}

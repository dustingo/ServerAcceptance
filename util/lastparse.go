package util

import (
	"github.com/BurntSushi/toml"
)

// LastCheckInfo https://xuri.me/toml-to-go/

type LastCheckInfo struct {
	Service   Service             `toml:"service"`
	Package   map[string]PackInfo `toml:"package"`
	Directory map[string]DirInfo  `toml:"directory"`
	Ulimit    Ulimit              `toml:"ulimit"`
}
type Off struct {
	State int      `toml:"state"`
	Label string   `toml:"label"`
	Name  []string `toml:"name"`
}
type On struct {
	State int      `toml:"state"`
	Label string   `toml:"label"`
	Name  []string `toml:"name"`
}
type Service struct {
	Off Off `toml:"off"`
	On  On  `toml:"on"`
}

/*
type Yum struct {
	State  int      `toml:"state"`
	Name   []string `toml:"name"`
}

type Pip struct {
	State    int      `toml:"state"`
	Name     []string `toml:"name"`
}
type Perl struct {
	State    int      `toml:"state"`
	Name     []string `toml:"name"`
}
type Package struct {
	Yum Yum `toml:"yum"`
	Pip Pip `toml:"pip"`
	Perl Perl `toml:"perl"`
}
*/
type PackInfo struct {
	State int      `toml:"state"`
	Label string   `toml:"label"`
	Name  []string `toml:"name"`
}
type DirInfo struct {
	State  int    `toml:"state"`
	Action string `toml:"action"`
	Path   string `toml:"path"`
	Mode   int    `toml:"mode"`
	Owner  string `toml:"owner"`
}
type Soft struct {
	Domain string `toml:"domain"`
	Type   string `toml:"type"`
	Item   string `toml:"item"`
	Value  string `toml:"value"`
}
type Hard struct {
	Domain string `toml:"domain"`
	Type   string `toml:"type"`
	Item   string `toml:"item"`
	Value  string `toml:"value"`
}
type Core struct {
	State  int    `toml:"state"`
	Action string `toml:"action"`
	Soft   Soft   `toml:"soft"`
	Hard   Hard   `toml:"hard"`
}
type Nofile struct {
	State  int    `toml:"state"`
	Action string `toml:"action"`
	Soft   Soft   `toml:"soft"`
	Hard   Hard   `toml:"hard"`
}
type Nproc struct {
	State  int    `toml:"state"`
	Action string `toml:"action"`
	Soft   Soft   `toml:"soft"`
	Hard   Hard   `toml:"hard"`
}
type Ulimit struct {
	Core   Core   `toml:"core"`
	Nofile Nofile `toml:"nofile"`
	Nproc  Nproc  `toml:"nproc"`
}

// GetService 获取service区域方法
func (l *LastCheckInfo) GetService() (int, string, []string, int, string, []string) {
	offState := l.Service.Off.State
	offLabel := l.Service.Off.Label
	offName := l.Service.Off.Name
	onState := l.Service.On.State
	onLabel := l.Service.On.Label
	onName := l.Service.On.Name
	return offState, offLabel, offName, onState, onLabel, onName
}

//GetPackage 获取Package区域方法
func (l *LastCheckInfo) GetPackage() (int, string, []string, int, string, []string, int, string, []string) {
	yumState := l.Package["yum"].State
	yumLabel := l.Package["yum"].Label
	yumName := l.Package["yum"].Name

	pipState := l.Package["pip"].State
	pipLabel := l.Package["pip"].Label
	pipName := l.Package["pip"].Name

	perlState := l.Package["perl"].State
	perlLabel := l.Package["perl"].Label
	perlName := l.Package["perl"].Name
	return yumState, yumLabel, yumName, pipState, pipLabel, pipName, perlState, perlLabel, perlName
}

//GetDirectory 获取directory区域方法
func (l *LastCheckInfo) GetDirectory() map[string]DirInfo {
	return l.Directory
}

func LastParse(fp string) (*LastCheckInfo, error) {
	var lastconfig LastCheckInfo
	_, err := toml.DecodeFile(fp, &lastconfig)
	if err != nil {
		return nil, err
	}
	return &lastconfig, nil
}

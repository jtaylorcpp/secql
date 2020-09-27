package osquery

/*
[
  {"arch":"x86_64","build":"","codename":"focal","major":"20","minor":"4","name":"Ubuntu","patch":"0","platform":"ubuntu","platform_like":"debian","version":"20.04.1 LTS (Focal Fossa)"}
]
*/
type OSInfo struct {
	Arch         string
	Build        string
	Codename     string
	Major        string
	Minor        string
	Name         string
	Patch        string
	Platform     string
	PlatformLike string `json:"platform_like"`
	Version      string
}

/*
	example cases:
		osqueryi --json "select * from Deb_packages limit 1"
		[
  			{"arch":"amd64","maintainer":"Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>","name":"accountsservice","priority":"optional","revision":"0ubuntu12~20.04.1","section":"admin","size":"452","source":"","status":"install ok installed","version":"0.6.55-0ubuntu12~20.04.1"}
		]
*/
type OSPackage struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Source     string `json:"source"`
	Size       string `json:"size"`
	Arch       string `json:"arch"`
	Revision   string `json:"revision"`
	Status     string `json:"status"`
	Maintainer string `json:"maintainer"`
	Section    string `json:"section"`
	Priority   string `json:"priority"`
}

/*
osqueryi --json "select distinct process.name, listening.port, listening.address, process.pid from processes as process join listening_ports as listening on process.pid = listening.pid limit 1"
[
  {"address":"127.0.0.1","name":"code","pid":"2207","port":"34797"}
]
*/
type ListeningApplication struct {
	Address string `json:"address`
	Name    string `json:"name"`
	Pid     string `json:"pid"`
	Port    string `josn:"port"`
}

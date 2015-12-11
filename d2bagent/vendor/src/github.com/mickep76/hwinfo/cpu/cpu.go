package cpu

// CPU information.
type CPU struct {
	Model          string `json:"model"`
	Flags          string `json:"flags"`
	Logical        int    `json:"logical"`
	Physical       int    `json:"physical"`
	Sockets        int    `json:"sockets"`
	CoresPerSocket int    `json:"cores_per_socket"`
	ThreadsPerCore int    `json:"threads_per_core"`
}

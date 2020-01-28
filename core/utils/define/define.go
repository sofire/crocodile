package define

// Role Admin or Normal User
type Role uint8

const (
	// NormalUser define normal user
	NormalUser Role = iota + 1 // 普通用户 只对自已创建的主机或者主机组具有操作权限
	// AdminUser define admin user
	AdminUser // 管理员 具有所有操作
)

func (r Role) String() string {
	switch r {
	case AdminUser:
		return "Admin"
	case NormalUser:
		return "Normal"
	default:
		return "Unknow Role type"
	}
}

// TaskType task type
// shell
// api
// TaskDataType
type TaskType uint8

const (
	// Shell Rum Command
	Shell TaskType = iota + 1
	// API run http req
	API
)

// RunMode crocodile run mode
// run crocodile as server or client
type RunMode uint8

const (
	// Server run crocodile as server
	Server RunMode = iota + 1
	// Client run crocodile as client
	Client
)

// TaskRespType task type (parent task,master task, child task)
// TaskRunType
type TaskRespType uint8

const (
	// MasterTask task as master run
	MasterTask TaskRespType = iota + 1
	// ParentTask task as a task's parent task run
	ParentTask
	// ChildTask task as a task's child task run
	ChildTask
)

func (tasktype TaskRespType) String() string {
	switch tasktype {
	case MasterTask:
		return "master"
	case ChildTask:
		return "child"
	case ParentTask:
		return "parent"
	default:
		return "unknow"
	}
}

// Task Status

// GetTaskid get task id in post
type GetTaskid struct {
	ID string `json:"id"`
}

// 定义结构体
type common struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	CreateTime string `json:"create_time,omitempty"` // 创建时间
	UpdateTime string `json:"update_time,omitempty"` // 最后一次更新时间
	Remark     string `json:"remark"`                // 备注
}

// User user msg
type User struct {
	Role      Role   `json:"role" binding:"gte=1,lte=2"`     // 用户类型: 2 管理员 1 普通用户
	Forbid    int    `json:"forbid" binding:"gte=0,lte=1"`   // 禁止用户: 0 未禁止 1 已禁止
	Password  string `json:"password,omitempty"`             // 用户密码
	Email     string `json:"email" binding:"required,email"` // 用户邮箱 日后任务的通知信息会发送给此邮件
	DingPhone string `json:"dingphone"`                      // dingding phone
	Slack     string `json:"slack"`                          // slack user name
	Telegram  string `json:"telegram"`                       // telegram bot chat id
	WeChat    string `json:"wechat"`                         // wechat id
	common
}

// HostGroup define hostgroup
type HostGroup struct {
	HostsID     []string `json:"addrs"`        // 主机host
	CreateByUID string   `json:"create_byuid"` // 创建人ID
	CreateBy    string   `json:"create_by"`    // 创建人ID
	common
}

// Host worker host
type Host struct {
	common
	Addr               string   `json:"addr"`     // 主机IP
	HostName           string   `json:"hostname"` // 主机名
	Online             int      `json:"online"`   // 主机是否在线 0 not online,1 online
	Weight             int      `json:"weight"`   // 主机权重
	RunningTasks       []string `json:"running_tasks"`
	Version            string   `json:"version"` // 版本号
	Stop               int      `json:"stop"`    // 1 为不能运行 0 为可以运行
	LastUpdateTimeUnix int64    `json:"last_updatetimeunix"`
	LastUpdateTime     string   `json:"last_updatetime"`
}

// Task define Task
type Task struct {
	TaskType TaskType    `json:"task_type"` // 任务类型
	TaskData interface{} `json:"task_data"` // 任务数据
	// TODO 改为 0 后只是在调度循环中不再检测下次运行时间
	Run               int         `json:"run"`                              // 0 为不能运行 1 为可以运行 如果这个任务作为别的任务父任务或者子任务会忽略这个字段
	ParentTaskIds     []string    `json:"parent_taskids"`                   // 父任务 运行任务前先运行父任务 以父或子任务运行时 任务不会执行自已的父子任务，防止循环依赖
	ParentRunParallel int         `json:"parent_runparallel"`               // 是否以并行运行父任务 0否 1是
	ChildTaskIds      []string    `json:"child_taskids"`                    // 子任务 运行结束后运行子任务
	ChildRunParallel  int         `json:"child_runparallel"`                // 是否以并行运行子任务 0否 1是
	CreateBy          string      `json:"create_by"`                        // 创建人
	CreateByUID       string      `json:"create_byuid"`                     // 创建人ID
	HostGroup         string      `json:"host_group"`                       // 执行计划
	HostGroupID       string      `json:"host_groupid" binding:"required"`  // 主机组ID
	Cronexpr          string      `json:"cronexpr" binding:"required"`      // 执行任务表达式
	Timeout           int         `json:"timeout"`                          // 任务超时时间 (s)
	AlarmUserIds      []string    `json:"alarm_userids" binding:"required"` // 报警用户 多个用户
	RoutePolicy       RoutePolicy `json:"route_policy" binding:"required"`  // how to select a run worker from hostgroup
	ExpectCode        int         `json:"expect_code"`                      // expect task return code. if not set 0 or 200
	ExpectContent     string      `json:"expect_content"`                   // expect task return content. if not set do not check
	AlarmStatus       int         `json:"alatm_status" binding:"required"`  // alarm when task run success or fail or all all:-2 failed: -1 success: 1
	common
}

// RoutePolicy set a task hot to select run worker
type RoutePolicy uint8

const (
	// Random get host by random
	Random RoutePolicy = iota + 1
	// RoundRobin get host by order
	RoundRobin
	// Weight get host by host weight
	Weight
	// LeastTask get host by host LeastTask
	LeastTask
)

// RunTask running task message
type RunTask struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	StartTimeStr string       `json:"start_timestr"`
	StartTime    int64        `json:"start_time"`
	RunTime      int          `json:"run_time"`
	Tasktype     TaskRespType `json:"tasktype"`
	TaskTypeStr  string       `json:"task_typestr"` // 1 主任务 2 父任务 3 子任务
	RunHost      string       `json:"run_host"`
}

// TaskResp run task resp message
type TaskResp struct {
	TaskID      string       `json:"task_td"`
	LogData     string       `json:"resp_data"`    // task run log data
	Code        int          `json:"code"`         // return code
	TaskType    TaskRespType `json:"task_type"`    // 1 主任务 2 父任务 3 子任务
	TaskTypeStr string       `json:"task_typestr"` // 1 主任务 2 父任务 3 子任务
	RunHost     string       `json:"run_host"`     // task run host
}

// Log task log
type Log struct {
	RunByTaskID    string       `json:"run_bytaskId"`
	StartTime      int64        `json:"start_time"`      // ms
	StartTimeStr   string       `json:"start_timestr"`   //
	EndTime        int64        `json:"end_timeunix"`    // ms
	EndTimeStr     string       `json:"end_timestr"`     //
	TotalRunTime   int          `json:"total_runtime"`   // ms
	Status         int          `json:"status"`          // 任务运行结果 -1 失败 1 成功
	TaskResps      []*TaskResp  `json:"task_resps"`      // 任务执行过程日志
	ErrCode        int          `json:"err_code"`        // err code
	ErrMsg         string       `json:"err_msg"`         // 错误原因
	ErrTasktype    TaskRespType `json:"err_tasktype"`    // err task type
	ErrTaskTypeStr string       `json:"err_tasktypestr"` // 1 主任务 2 父任务 3 子任务
	ErrTaskID      string       `json:"err_taskid"`      // task failed id
	ErrTask        string       `json:"err_task"`        // task failed id
}

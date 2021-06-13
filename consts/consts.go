package consts

const (
	ON  = 1
	OFF = -1
)

// parmas
const (
	METHOD    = "method"
	PARAMS    = "params"
	CODE      = "code"
	ERROR     = "error"
	DATA      = "data"
	DATA_SIZE = "data_size"
	SLOT      = "slot"
	FLAGS     = "flags"
)

// chan
const (
	DEV_DATA         = "dev.data"         // 数据
	DEV_DATA_GET     = "dev.data.get"     // 读数据
	DEV_DATA_SET     = "dev.data.set"     // 写数据
	DEV_DATA_SIMPLE  = "dev.data.sim"     // 数据(简单模式)
	DEV_META         = "dev.meta"         // 元数据
	DEV_META_GET     = "dev.meta.get"     // 读元数据
	DEV_META_SET     = "dev.meta.set"     // 写元数据
	DEV_HIS_DATA     = "dev.his.data"     // 历史数据
	DEV_HIS_DATA_ACK = "dev.his.data.ack" // 历史数据反馈
)

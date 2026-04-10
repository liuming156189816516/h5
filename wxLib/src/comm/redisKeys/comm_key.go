package redisKeys

import (
	"comm/comm"
	"fmt"
)

//ttl
const (
	TTL_DAY   = 24 * 3600
	TTL_WEEK  = 7 * 24 * 3600
	TTL_MONTH = 31 * 24 * 3600
)

//token 信息 永久存
func GetUserTokenInfo() string {
	return fmt.Sprintf("%s_h_comm_user_token_info", comm.GetMgoDBName())
}

//inc
func GetAdminIncInfo() string {
	return fmt.Sprintf("%s_h_comm_user_inc_info", comm.GetMgoDBName())
}

//菜单
func GetAdminMenuInfo() string {
	return fmt.Sprintf("%s_o_webface_all_admin_menu_key", comm.GetMgoDBName())
}

//角色
func GetAdminRoleInfo() string {
	return fmt.Sprintf("%s_h_webface_all_admin_role_key", comm.GetMgoDBName())
}

//缓存用户信息
func GetAdminUserInfo() string {
	return fmt.Sprintf("%s_h_webface_all_admin_user_key", comm.GetMgoDBName())
}

//锁key
func GetLockKey(lock string) string {
	return fmt.Sprintf("%s_o_comm_lock_key_%s", comm.GetMgoDBName(), lock)
}

//正在工作的列表
func GetDoTaskList() string {
	return fmt.Sprintf("%s_h_task_do_task_list_key", comm.GetMgoDBName())
}

//缓存数据包内容
func GetDataPackListKey( keyStr string) string {
	return fmt.Sprintf("%s_u_data_pack_list_key_1_%s", comm.GetMgoDBName(), keyStr)
}

//缓存数据包异常内容
func GetDataPackListErrKey( keyStr string) string {
	return fmt.Sprintf("%s_u_data_pack_list_key_err_%s", comm.GetMgoDBName(), keyStr)
}

//缓存数据包内容-所有
func GetDataPackListKey2( keyStr string) string {
	return fmt.Sprintf("%s_u_data_pack_list_key_2_%s", comm.GetMgoDBName(), keyStr)
}

//一个账号下数据内容去重
func GetAllDataPackListKey() string {
	return fmt.Sprintf("%s_u_data_pack_list_key_3", comm.GetMgoDBName())
}

//获取任务状态
func GetScheduleKey() string {
	return fmt.Sprintf("%s_u_get_schedule_key", comm.GetMgoDBName())
}

//ip使用次数
func GetIpUserNumKey() string {
	return fmt.Sprintf("%s_ip_user_num_key", comm.GetMgoDBName())
}

//系统配置
func GetSysConfig() string {
	return fmt.Sprintf("%s_sys_config", comm.GetMgoDBName())
}

// 验证码信息
func GetAdamCodeInfo(uuid string) string {
	return fmt.Sprintf("%s_admin_code_key_%s", comm.GetMgoDBName(), uuid)
}

// app用户信息
func GetAppUserInfoKey() string {
	return fmt.Sprintf("%s_app_user_info_key", comm.GetMgoDBName())
}

// app用户余额
func AppUserBalanceKey() string {
	return fmt.Sprintf("%s_app_user_balance_key", comm.GetMgoDBName())
}

// 任务配置
func GetTaskConfigKey() string {
	return fmt.Sprintf("%s_task_config_key", comm.GetMgoDBName())
}

// app 任务编码
func AppUserTaskNoKey() string {
	return fmt.Sprintf("%s_app_user_task_no_key", comm.GetMgoDBName())
}

//提现列表
func GetAllWithdrawListKey() string {
	return fmt.Sprintf("%s_all_withdraw_list_key", comm.GetMgoDBName())
}

//当天新增活跃用户
func GetTodayNewActiveUserNumKey(timeStr, tuid string) string {
	return fmt.Sprintf("%s_today_new_active_user_num_key_%s", comm.GetUserMgoDBName(tuid), timeStr)
}

//当天活跃用户
func GetTodayActiveUserNumKey(timeStr, tuid string) string {
	return fmt.Sprintf("%s_today_active_user_num_key_%s", comm.GetUserMgoDBName(tuid), timeStr)
}

//当天新增用户
func GetTodayNewUserNumKey(timeStr, tuid string) string {
	return fmt.Sprintf("%s_active_user_num_key_%s", comm.GetUserMgoDBName(tuid), timeStr)
}

//登录人数
func GetTodayLoginNumKey(timeStr, tuid string) string {
	return fmt.Sprintf("%s_login_num_key_%s", comm.GetUserMgoDBName(tuid), timeStr)
}

//站内信列表
func GetAllMessageListKey() string {
	return fmt.Sprintf("%s_all_message_list_key", comm.GetMgoDBName())
}

//用户站内信
func GetUserMessageKey(id string) string {
	return fmt.Sprintf("%s_user_message_key", comm.GetUserMgoDBName(id))
}

//系统配置
func GetSysConfigKey() string {
	return fmt.Sprintf("%s_all_sys_config_key", comm.GetMgoDBName())
}

//渠道访问量
func GetJudaoTongji1Key(liveCode string) string {
	return fmt.Sprintf("%s_qudao_key_1", liveCode)
}

//渠道点击量
func GetJudaoTongji2Key(liveCode string) string {
	return fmt.Sprintf("%s_qudao_key_2", liveCode)
}

//渠道访问量
func GetJudaoTongji3Key() string {
	return "qudao_key_3"
}

// 账号对应socks5代理
func GetAccountSocks5Key() string {
	return fmt.Sprintf("vps_account_socks5")
}

// fb上报
func GetFbReportKey() string {
	return fmt.Sprintf("fb_report_key")
}

// check任务状态
func GetTaskStatusKey() string {
	return fmt.Sprintf("%s_task_status_key", comm.GetMgoDBName())
}

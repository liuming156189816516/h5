package redisKeys

import (
	"comm/comm"
	"fmt"
)

//账号的状态
func GetAccountStatusKey() string {
	return fmt.Sprintf("%s_u_account_status_key", comm.GetMgoDBName())
}

//账号的ip
func GetAccountProxyIpKey() string {
	return fmt.Sprintf("%s_u_account_proxy_ip_key", comm.GetMgoDBName())
}

//账号列表
func GetAccountListInfoKey() string {
	return fmt.Sprintf("%s_account_list_info_key", comm.GetMgoDBName())
}

//监控号对应的消息id
func GetMonitorMsgIdKey(timeStr string) string {
	return fmt.Sprintf("%s_monitor_msg_id_key", comm.GetUserMgoDBName(timeStr))
}

//监控号发现对应的消息被删除了
func GetMonitorMsgDelKey(timeStr string) string {
	return fmt.Sprintf("%s_monitor_msg_del_key", comm.GetUserMgoDBName(timeStr))
}
//监控号发现对应的消息收到了
func GetMonitorMsgReceiveKey(timeStr string) string {
	return fmt.Sprintf("%s_monitor_msg_receive_key", comm.GetUserMgoDBName(timeStr))
}
//真机消息发送状态
func GetMonitorMsgSendStatusKey(timeStr string) string {
	return fmt.Sprintf("%s_monitor_msg_send_status_key", comm.GetUserMgoDBName(timeStr))
}
//真机消息手机号对应的ws号
func GetMonitorMsgPhoneWsKey(timeStr string) string {
	return fmt.Sprintf("%s_monitor_msg_phone_ws_key", comm.GetUserMgoDBName(timeStr))
}

//登录列表
func GetAllTaskLoginList() string {
	return fmt.Sprintf("%s_l_task_all_task_login_list_key", comm.GetMgoDBName())
}

//消息回调列表
func GetAllTaskMessageResultList() string {
	return fmt.Sprintf("%s_l_task_all_task_message_result_list_key", comm.GetMgoDBName())
}

//处理接听电话回调
func GetAllTaskAcceptCallList() string {
	return fmt.Sprintf("%s_l_task_all_task_accept_call_list_key", comm.GetMgoDBName())
}

//处理炸群任务
func GetAllTaskZhaGroupList() string {
	return fmt.Sprintf("%s_l_task_all_task_zha_group_list_key", comm.GetMgoDBName())
}

//任务列表
func GetAllTaskList() string {
	return fmt.Sprintf("%s_l_task_all_task_list_key", comm.GetMgoDBName())
}


//账号的uid
func GetAccountUidKey() string {
	return fmt.Sprintf("%s_account_uid_key", comm.GetMgoDBName())
}

//真机私发检测任务
func GetAiApprovalTaskKey() string {
	return fmt.Sprintf("%s_ai_approval_task_key", comm.GetMgoDBName())
}

//真机私发获取数据包ws账号任务
func GetAiDataPackWsTaskKey() string {
	return fmt.Sprintf("%s_ai_data_pack_ws_task_key", comm.GetMgoDBName())
}

func GetAiUuidIncInfo() string {
	return fmt.Sprintf("%s_h_ai_uuid_inc_info", comm.GetMgoDBName())
}

// 检测qrcode
func GetCheckQrcodeTaskKey(uuid string) string {
	return fmt.Sprintf("check_qrcode_tasky_key_%s", uuid)
}

//拉群任务结算列表
func GetAllCreateGroupSettlementListKey() string {
	return fmt.Sprintf("%s_all_create_group_settlement_list_key", comm.GetMgoDBName())
}


//拉粉任务数据去重
func GetGroupLinkPackListKey() string {
	return fmt.Sprintf("%s_u_group_link_pack_list_key", comm.GetMgoDBName())
}

//拉群任务数据去重
func GetCreateGroupPackListKey() string {
	return fmt.Sprintf("%s_u_create_group_pack_list_key", comm.GetMgoDBName())
}

//所有推广号,避免推广号重复
func GetAllAccountListKey() string {
	return fmt.Sprintf("%s_all_account_list_key", comm.GetMgoDBName())
}

//ai私发上传图片任务
func GetAiUploadTaskWsKey() string {
	return fmt.Sprintf("%s_ai_upload_task_ws_key", comm.GetMgoDBName())
}

//ai提交审批任务
func GetAiApprovalTaskWsKey() string {
	return fmt.Sprintf("%s_ai_approval_task_ws_key", comm.GetMgoDBName())
}

//ai获取审批结果任务
func GetAiTaskBatchResultWsKey() string {
	return fmt.Sprintf("%s_ai_task_batch_result_ws_key", comm.GetMgoDBName())
}

//挂机时长
func GetAccountDurationKey() string {
	return fmt.Sprintf("%s_account_Duration", comm.GetMgoDBName())
}

//发消消息明细
func GetSendMsgTaskInfoKey(account string) string {
	return fmt.Sprintf("%s_send_msg_task_info_key", comm.GetUserMgoDBName(account))
}

//自动发消消息明细
func GetAutoSendMsgTaskInfoKey(account string) string {
	return fmt.Sprintf("%s_auto_send_msg_task_info_key", comm.GetUserMgoDBName(account))
}

//发消消息phone-lid
func GetSendMsgPhoneLidKey(account string) string {
	return fmt.Sprintf("%s_send_msg_phone_lid_key", comm.GetUserMgoDBName(account))
}

//acc_target对应的taskid
func GetAccTargetTaskIdKey() string {
	return fmt.Sprintf("%s_acc_target_task_id_key", comm.GetMgoDBName())
}

//计数
func GetSendMsgTaskCountKey() string {
	return fmt.Sprintf("%s_send_msg_task_count_key", comm.GetMgoDBName())
}

//计数
func GetSendMsgTaskInfoCountKey() string {
	return fmt.Sprintf("%s_send_msg_task_info_count_key", comm.GetMgoDBName())
}

//获取所有群发任务
func GetAllSendMsgTaskList() string {
	return fmt.Sprintf("%s_all_send_msg_task_list_key", comm.GetMgoDBName())
}

//私发
func GetAllTaskSendmsgList() string {
	return fmt.Sprintf("%s_l_task_all_task_sendmsg_list_key", comm.GetMgoDBName())
}

//检测账号是否可以发消息
func GetCheckAccountListKey() string {
	return fmt.Sprintf("%s_check_account_list_key", comm.GetMgoDBName())
}

//拉粉奖励任务数据去重
func GetGroupLinkPackAwardListKey() string {
	return fmt.Sprintf("%s_group_link_pack_award_list_key", comm.GetMgoDBName())
}

//获取所有自动群发任务
func GetAutoAllSendMsgTaskList() string {
	return fmt.Sprintf("%s_auto_all_send_msg_task_list_key", comm.GetMgoDBName())
}




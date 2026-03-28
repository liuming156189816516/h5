package tableName
import (
	"fmt"
	"strings"
)

// 邀请码
func GetTableInviteCodeList() string {
	return "u_invite_code_list"
}

// 轮播图
func GetTableCarouselList() string {
	return "u_carousel_list"
}

// 提现审批
func GetTableWithdrawApprovalList() string {
	return "u_withdraw_approval_list"
}

// 提现账户
func GetTableWithdrawCardList() string {
	return "u_withdraw_card_list"
}

// 人工修正
func GetTableWithdrawReviseList() string {
	return "u_withdraw_revision_list"
}

// 账单明细
func GetTableBillRecordList() string {
	return "u_bill_record"
}

// 账单明细历史
func GetTableBillRecordHisList(ym string) string {
	return fmt.Sprintf("u_bill_record_his_%s", ym)
}

// check账单明细历史
func CheckTableBillRecordHisList(tb string) bool {
	return strings.Index(tb, "u_bill_record_his_") == 0
}

// 任务明细
func GetTableTaskRecordList() string {
	return "u_task_record"
}

// 数据统计
func GetTableDataStatisList() string {
	return "u_data_statis"
}

// 建群任务详情
func GetTableCreateGroupInfo() string {
	return "u_create_group_info"
}


// 轮播图
func GetTableContactList() string {
	return "u_contact_list"
}

// 发送消息任务
func GetTableSendMsgTaskListInfo() string {
	return "u_send_msg_task"
}

// 发送消息详情
func GetTableSendMsgInfoListInfo() string {
	return "u_send_msg_info"
}

// 抽奖记录
func GetTableLotteryRecordList() string {
	return "u_lottery_record"
}
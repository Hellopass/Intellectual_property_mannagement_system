package models

// 审批状态常量定义
// 用于管理著作的审批流程状态
const (
	ApprovalStatusInProgress = 0 // 审批进行中（存在未完成的审批步骤）
	ApprovalStatusApproved   = 1 // 审批已通过（所有步骤完成且通过）
	ApprovalStatusRejected   = 2 // 审批已驳回（任一环节驳回）
)

// 审批步骤常量定义
// 表示当前所处的审批阶段
const (
	ApprovalStepInitial = 1 // 初审阶段（第一个审批环节）
	ApprovalStepFinal   = 2 // 终审阶段（最终审批环节）
)

// 特殊说明：
// 审批流程逻辑：
// 1. 新建著作默认状态为 ApprovalStatusInProgress
// 2. 审批步骤从 ApprovalStepInitial 开始
// 3. 当初审通过后，进入 ApprovalStepFinal 阶段
// 4. 任一环节驳回即终止流程，状态变为 ApprovalStatusRejected
// 5. 终审通过后状态变为 ApprovalStatusApproved

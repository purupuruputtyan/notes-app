package apperror

// クライアント向けエラーコード（採番はこのファイルのみで管理する）。
const (
	CodeUserNickNameRequired     = "ERR_001"
	CodeUserNickNameTooLong      = "ERR_002"
	CodeUserInvalidEmail         = "ERR_003"
	CodeUserPasswordTooShort     = "ERR_004"
	CodeUserInvalidPassword      = "ERR_005"
	CodeUserNotFound             = "ERR_006"
	CodeUserEmailAlreadyExists   = "ERR_007"
	CodeUserNickNameAlreadyTaken = "ERR_008"
	CodeUserInvalidLogin         = "ERR_009"
	CodeNoteOwnerNotFound        = "ERR_010"
)

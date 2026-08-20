package models

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

/** ====================================================================================
 * 🏁 AuthType
 * =====================================================================================
 */

// UserAuthTypeSystem 自有平台(手机登录)
var UserAuthTypeSystem = "system"

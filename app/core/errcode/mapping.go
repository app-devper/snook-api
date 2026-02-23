package errcode

import "net/http"

type CodeInfo struct {
	HttpStatus  int
	Description string
}

var codeMap = map[string]CodeInfo{
	// ─── Auth / Middleware ───────────────────────────────────────────────────
	AU_UNAUTHORIZED_001: {http.StatusUnauthorized, "missing or invalid authorization header"},
	AU_UNAUTHORIZED_002: {http.StatusUnauthorized, "token invalid or expired"},
	AU_UNAUTHORIZED_003: {http.StatusUnauthorized, "system mismatch"},
	AU_UNAUTHORIZED_004: {http.StatusUnauthorized, "clientId mismatch"},
	AU_UNAUTHORIZED_005: {http.StatusUnauthorized, "session invalid"},

	// ─── Table (TB) ─────────────────────────────────────────────────────────
	TB_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	TB_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	TB_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Table Session (TS) ─────────────────────────────────────────────────
	TS_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	TS_BAD_REQUEST_002: {http.StatusBadRequest, "operation failed"},
	TS_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Booking (BK) ───────────────────────────────────────────────────────
	BK_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	BK_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	BK_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Menu Category (MC) ─────────────────────────────────────────────────
	MC_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	MC_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	MC_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Menu Item (MI) ─────────────────────────────────────────────────────
	MI_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	MI_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	MI_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Table Order (TO) ───────────────────────────────────────────────────
	TO_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	TO_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	TO_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Payment (PY) ───────────────────────────────────────────────────────
	PY_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	PY_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	PY_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Creditor (CR) ──────────────────────────────────────────────────────
	CR_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	CR_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	CR_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Promotion (PM) ─────────────────────────────────────────────────────
	PM_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	PM_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	PM_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Expense (EX) ───────────────────────────────────────────────────────
	EX_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	EX_BAD_REQUEST_002: {http.StatusBadRequest, "create/update/delete failed"},
	EX_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Setting (SE) ───────────────────────────────────────────────────────
	SE_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request body"},
	SE_BAD_REQUEST_002: {http.StatusBadRequest, "upsert failed"},
	SE_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Dashboard (DA) ─────────────────────────────────────────────────────
	DA_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request or missing params"},
	DA_BAD_REQUEST_002: {http.StatusBadRequest, "query failed"},
	DA_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── Report (RP) ────────────────────────────────────────────────────────
	RP_BAD_REQUEST_001: {http.StatusBadRequest, "invalid request or missing params"},
	RP_BAD_REQUEST_002: {http.StatusBadRequest, "report generation failed"},
	RP_INTERNAL_001:    {http.StatusInternalServerError, "internal server error"},

	// ─── System (SY) ────────────────────────────────────────────────────────
	SY_NOT_FOUND_001: {http.StatusNotFound, "route not found"},
	SY_FORBIDDEN_001: {http.StatusForbidden, "invalid request, restricted endpoint"},
	SY_FORBIDDEN_002: {http.StatusForbidden, "don't have permission"},
	SY_INTERNAL_001:  {http.StatusInternalServerError, "internal server error"},
}

func GetCodeInfo(code string) (CodeInfo, bool) {
	info, ok := codeMap[code]
	return info, ok
}

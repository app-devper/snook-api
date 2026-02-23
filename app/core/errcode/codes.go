package errcode

// ─── Auth / Middleware ───────────────────────────────────────────────────────
const (
	AU_UNAUTHORIZED_001 = "AU-401-001" // missing / invalid authorization header
	AU_UNAUTHORIZED_002 = "AU-401-002" // token invalid or expired
	AU_UNAUTHORIZED_003 = "AU-401-003" // system mismatch
	AU_UNAUTHORIZED_004 = "AU-401-004" // clientId mismatch
	AU_UNAUTHORIZED_005 = "AU-401-005" // session invalid
)

// ─── Table (TB) ─────────────────────────────────────────────────────────────
const (
	TB_BAD_REQUEST_001 = "TB-400-001" // invalid request body
	TB_BAD_REQUEST_002 = "TB-400-002" // create/update/delete failed
	TB_INTERNAL_001    = "TB-500-001" // internal server error
)

// ─── Table Session (TS) ─────────────────────────────────────────────────────
const (
	TS_BAD_REQUEST_001 = "TS-400-001" // invalid request body
	TS_BAD_REQUEST_002 = "TS-400-002" // create/update/delete failed
	TS_INTERNAL_001    = "TS-500-001" // internal server error
)

// ─── Booking (BK) ───────────────────────────────────────────────────────────
const (
	BK_BAD_REQUEST_001 = "BK-400-001" // invalid request body
	BK_BAD_REQUEST_002 = "BK-400-002" // create/update/delete failed
	BK_INTERNAL_001    = "BK-500-001" // internal server error
)

// ─── Menu Category (MC) ─────────────────────────────────────────────────────
const (
	MC_BAD_REQUEST_001 = "MC-400-001" // invalid request body
	MC_BAD_REQUEST_002 = "MC-400-002" // create/update/delete failed
	MC_INTERNAL_001    = "MC-500-001" // internal server error
)

// ─── Menu Item (MI) ─────────────────────────────────────────────────────────
const (
	MI_BAD_REQUEST_001 = "MI-400-001" // invalid request body
	MI_BAD_REQUEST_002 = "MI-400-002" // create/update/delete failed
	MI_INTERNAL_001    = "MI-500-001" // internal server error
)

// ─── Table Order (TO) ───────────────────────────────────────────────────────
const (
	TO_BAD_REQUEST_001 = "TO-400-001" // invalid request body
	TO_BAD_REQUEST_002 = "TO-400-002" // create/update/delete failed
	TO_INTERNAL_001    = "TO-500-001" // internal server error
)

// ─── Payment (PY) ───────────────────────────────────────────────────────────
const (
	PY_BAD_REQUEST_001 = "PY-400-001" // invalid request body
	PY_BAD_REQUEST_002 = "PY-400-002" // create/update/delete failed
	PY_INTERNAL_001    = "PY-500-001" // internal server error
)

// ─── Creditor (CR) ──────────────────────────────────────────────────────────
const (
	CR_BAD_REQUEST_001 = "CR-400-001" // invalid request body
	CR_BAD_REQUEST_002 = "CR-400-002" // create/update/delete failed
	CR_INTERNAL_001    = "CR-500-001" // internal server error
)

// ─── Promotion (PM) ─────────────────────────────────────────────────────────
const (
	PM_BAD_REQUEST_001 = "PM-400-001" // invalid request body
	PM_BAD_REQUEST_002 = "PM-400-002" // create/update/delete failed
	PM_INTERNAL_001    = "PM-500-001" // internal server error
)

// ─── Expense (EX) ───────────────────────────────────────────────────────────
const (
	EX_BAD_REQUEST_001 = "EX-400-001" // invalid request body
	EX_BAD_REQUEST_002 = "EX-400-002" // create/update/delete failed
	EX_INTERNAL_001    = "EX-500-001" // internal server error
)

// ─── Setting (SE) ───────────────────────────────────────────────────────────
const (
	SE_BAD_REQUEST_001 = "SE-400-001" // invalid request body
	SE_BAD_REQUEST_002 = "SE-400-002" // upsert failed
	SE_INTERNAL_001    = "SE-500-001" // internal server error
)

// ─── Dashboard (DA) ─────────────────────────────────────────────────────────
const (
	DA_BAD_REQUEST_001 = "DA-400-001" // invalid request / missing params
	DA_BAD_REQUEST_002 = "DA-400-002" // query failed
	DA_INTERNAL_001    = "DA-500-001" // internal server error
)

// ─── Report (RP) ────────────────────────────────────────────────────────────
const (
	RP_BAD_REQUEST_001 = "RP-400-001" // invalid request / missing params
	RP_BAD_REQUEST_002 = "RP-400-002" // report generation failed
	RP_INTERNAL_001    = "RP-500-001" // internal server error
)

// ─── System (SY) ────────────────────────────────────────────────────────────
const (
	SY_NOT_FOUND_001 = "SY-404-001" // route not found
	SY_FORBIDDEN_001 = "SY-403-001" // invalid request / restricted endpoint
	SY_FORBIDDEN_002 = "SY-403-002" // no permission
	SY_INTERNAL_001  = "SY-500-001" // panic recovery / internal server error
)

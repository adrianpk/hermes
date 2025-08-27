package am

const (
	CSRFFieldName = "aquamarine.csrf.token"
	TimeFormat    = "2006-01-02T15:04:05Z07:00"
)

const (
	NoSlug = ""
	Slug   = "?slug=%s"
)

const (
	MsgGetAllItems = "%s retrieved successfully"
	MsgGetItem     = "%s retrieved successfully"
	MsgCreateItem  = "%s created successfully"
	MsgUpdateItem  = "%s updated successfully"
	MsgDeleteItem  = "%s deleted successfully"

	ErrInvalidID            = "Invalid %s ID"
	ErrCannotGetResources   = "Could not get %s from database"
	ErrCannotGetResource    = "Could not get %s from database"
	ErrCannotCreateResource = "Could not create %s in database"
	ErrCannotUpdateResource = "Could not update %s in database"
	ErrCannotDeleteResource = "Could not delete %s from database"
)

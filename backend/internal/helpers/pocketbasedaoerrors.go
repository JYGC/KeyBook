package helpers

func IsNoRowsResult(err error) bool {
	return err.Error() == "sql: no rows in result set"
}

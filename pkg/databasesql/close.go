package databasesql

//Close ...
func (d DatabaseSQL) Close() error {
	err := d.Db.Close()
	return err
}

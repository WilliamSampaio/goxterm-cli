package database

func CreateInfoTable() Migration {
	return Migration{
		Id:   1,
		Name: "2025_08_13_create_info_table",
		Up: `CREATE TABLE info (
			app_version TEXT(10)
		);`,
		Down: `DROP TABLE info;`,
	}
}

package g

func GetSqlData(sql string)(results []map[string]string,b bool)  {
	results, err := Engine.QueryString(sql)

	if err!=nil{
		Logger().Error("错误运行 sql :%s", sql)
		return results,false
	}

	if len(results) == 0  {
		//results[0]["UNIXSTMAP"] = "0"
		return results ,false
	}


	return results,true
}


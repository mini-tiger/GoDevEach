package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-oci8"
	"os"
	"time"
	"fmt"
)

func main() {
	os.Setenv("NLS_LANG", "")
	//if len(os.Args) != 2 {
	//	log.Fatalln(os.Args[0] + " user/password@host:port/sid")
	//}

	db, err := sql.Open("oci8", "test1/test1@192.168.43.236:1521/orcl")
	//db, err := sql.Open("oci8", "test1/test1@1.119.132.155:1521/orcl")
	//fmt.Printf("%+v\n",db)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	//rows, err := db.Query("select vehicle_serial from tf_op_vehicle where train_serial = '3D55FBC5EAD147A390CA3AFD775308C4' group by vehicle_serial")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var data string
	//	rows.Scan(&data)
	//	for i:=0;i<11;i++ {
	//		result, err := db.Exec(fmt.Sprintf("insert into SYS1 (ID, TRAIN_SERIAL, VEHICLE_SERIAL, VEHICLE_ORDER, PART_CODE, IMG_NAME, FAULT_TYPE, X1POS, Y1POS, X2POS, Y2POS, IS_DISPOSE, DISPOSE_RESULT, DISPOSE_MAN, DISPOSE_DT, IMG_INDEX, FAULT_SERIAL)"+
	//			"values ('C7183C33594D4742BC200E7A6BF0DBAB', '3D55FBC5EAD147A390CA3AFD775308C4', '%s',"+
	//			" 31, 9, '33_4_%s', 15, 937, 55, 325, 213, 0, null, null, null, 13, null)", data,strconv.Itoa(i)))
	//		fmt.Println(result, err)
	//	}
	//}
	//if err = rows.Err(); err != nil {
	//	log.Fatalln(err)
	//}

	t1:=time.Now()
	t2:=t1.Add(time.Duration(-10)*time.Second)
	sql11:=fmt.Sprintf("update tf_op_train set pass_time = to_date('%s', 'YYYY/MM/DD HH24:MI:SS') " +
		"where train_serial = '3D55FBC5EAD147A390CA3AFD775308C4'",t2.Format("2006-01-02 15:04:05"))


	result,err:=db.Exec(sql11)
	fmt.Println(result,err)
	//
	//t3:=t1.Add(time.Duration(-10)*time.Second)
	//sql11=fmt.Sprintf("update tf_op_train set pass_time = to_date('%s', 'YYYY/MM/DD HH24:MI:SS') " +
	//	"where train_serial = 'F50E2841AD1E4E5781A6C37BBCBD753D'",t3.Format("2006-01-02 15:04:05"))
	//
	//
	//result,err=db.Exec(sql11)
	//fmt.Println(result,err)

	//t4:=t1.Add(time.Duration(30)*time.Second)
	//sql11=fmt.Sprintf("update tf_op_train set pass_time = to_date('%s', 'YYYY/MM/DD HH24:MI:SS') " +
	//	"where train_serial = '32F6ACC0B0A243FC974BC2BCECF2AA51'",t4.Format("2006-01-02 15:04:05"))
	//
	//
	//result,err=db.Exec(sql11)
	//fmt.Println(result,err)
}


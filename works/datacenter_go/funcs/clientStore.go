package authfunc

import (
	"context"
	"datacenter/modules"
	"errors"

	"github.com/go-oauth2/oauth2/v4"
	"gorm.io/gorm"
	"strings"
)

type ClientStore struct {
	db                *gorm.DB
	TableName         string
	initTableDisabled bool
}

// ClientStoreItem data item
//type ClientStoreItem struct {
//	ID     string `db:"id"`
//	Secret string `db:"secret"`
//	Domain string `db:"domain"`
//	Data   string `db:"data"`
//}

// NewClientStore creates PostgreSQL store instance
func NewClientStore(db *gorm.DB, tablename string, options ...ClientStoreOption) (*ClientStore, error) {

	store := &ClientStore{
		db:        db,
		TableName: tablename,
	}

	for _, o := range options {
		o(store)
	}

	var err error
	//fmt.Printf("%+v\n",store)
	if !store.initTableDisabled {
		err = store.initTable()
		//fmt.Println(err.Error())
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return nil, err
			}
		}
	}

	return store, nil
}

func (s *ClientStore) initTable() error {

	//	query := fmt.Sprintf(`
	//	CREATE TABLE IF NOT EXISTS %s (
	//		id VARCHAR(255) NOT NULL PRIMARY KEY,
	//		secret VARCHAR(255) NOT NULL,
	//		domain VARCHAR(255) NOT NULL,
	//		data TEXT NOT NULL
	//	  );
	//`, s.tableName)
	//
	//	stmt, err := s.db.Prepare(query)
	//	s.db.cr
	//	if err != nil {
	//		return err
	//	}
	//	_, err = stmt.Exec()
	//	if err != nil {
	//		return err
	//	}

	dm := s.db.Migrator()
	//fmt.Println(s.db.AutoMigrate(&modules.OauthClientDetails{}))
	//fmt.Println(dm.HasTable(s.tableName))
	//fmt.Printf(s.tableName)

	if !dm.HasTable(s.TableName) {
		db := s.db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;")
		dm = db.Migrator()
		return dm.CreateTable(&modules.OauthClientDetails{})
	}

	return nil
}

//func (s *ClientStore) toClientInfo(data string) (oauth2.ClientInfo, error) {
//	var cm modules.OauthClientDetails
//	err := jsoniter.Unmarshal([]byte(data), &cm)
//	return &cm, err
//}

// GetByID retrieves and returns client information by id
func (s *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	if id == "" {
		return nil, nil
	}

	var clientDetail modules.OauthClientDetails
	db := modules.MysqlDb.Table(s.TableName).Select("*").Where("client_id=?", id).First(&clientDetail)

	if db.Error != nil {
		return nil, db.Error
	}
	if db.RowsAffected == 0 {
		return &clientDetail, errors.New("Not Found")
	}

	return &clientDetail, nil
}

// Create creates and stores the new client information
func (s *ClientStore) Create(info *modules.OauthClientDetails) error {
	//data, err := jsoniter.Marshal(info)

	db := s.db.Table(s.TableName).Create(info)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (s *ClientStore) GetDetailByID(id string) (clientDetail *modules.OauthClientDetails, db *gorm.DB) {
	if id == "" {
		return
	}

	//var cd modules.OauthClientDetails
	db = modules.MysqlDb.Table(s.TableName).Select("*").Where("client_id=?", id).First(&clientDetail)
	return
}

func (s *ClientStore) GetDetailByWhere(where string) (clientDetail *modules.OauthClientDetails, db *gorm.DB) {

	//var cd modules.OauthClientDetails
	db = modules.MysqlDb.Table(s.TableName).Select("*").Where(where).First(&clientDetail)
	return
}

func (s *ClientStore) CreateDetail(info *modules.OauthClientDetails) *gorm.DB {
	//data, err := jsoniter.Marshal(info)

	return s.db.Table(s.TableName).Create(info)

}

func (s *ClientStore) Save(info *modules.OauthClientDetails) *gorm.DB {
	//data, err := jsoniter.Marshal(info)

	return s.db.Table(s.TableName).Save(info)

}

func (s *ClientStore) Delete(info *modules.OauthClientDetails) *gorm.DB {
	//data, err := jsoniter.Marshal(info)

	return s.db.Table(s.TableName).Delete(info)

}

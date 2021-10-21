package modules

import (
	"context"
	"database/sql"
	"datacenter/g"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

// StoreItem data item
//type StoreItem struct {
//	ID              int64  `db:"id,primarykey,autoincrement"`
//	ExpiredAt       int64  `db:"expired_at"`
//	AccessExpiredAt int64  `db:"access_expired_at"`
//	Code            string `db:"code,size:255"`
//	Access          string `db:"access,size:255"`
//	Refresh         string `db:"refresh,size:255"`
//	ClientID        string `db:"clientId,size:255"`
//	UserName        string `db:"userName,size:255"`
//	Data            string `db:"data,size:2048"`
//}

// NewConfig create mysql configuration instance
func NewConfig(dsn string) *Config {
	return &Config{
		DSN:          dsn,
		MaxLifetime:  time.Hour * 2,
		MaxOpenConns: 50,
		MaxIdleConns: 25,
	}
}

// Config mysql configuration
type Config struct {
	DSN          string
	MaxLifetime  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

// NewDefaultStore create mysql store instance
func NewDefaultStore(config *Config) *Store {
	return NewStore(config, "", 0)
}
func NewStoreEntry(s *Store) *StoreEntry {

	return &StoreEntry{s}
}

// NewStore create mysql store instance,
// config mysql configuration,
// tableName table name (default oauth2_token),
// GC time interval (in seconds, default 600)
func NewStore(config *Config, tableName string, gcInterval int) *Store {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	return NewStoreWithDB(db, tableName, gcInterval)
}

// NewStoreWithDB create mysql store instance,
// db sql.DB,
// tableName table name (default oauth2_token),
// GC time interval (in seconds, default 600)
func NewStoreWithDB(db *sql.DB, tableName string, gcInterval int) *Store {
	store := &Store{
		db:        &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Encoding: "UTF8", Engine: "MyISAM"}},
		tableName: "oauth2_token_keep2",
		stdout:    os.Stderr,
	}
	if tableName != "" {
		store.tableName = tableName
	}

	interval := 600
	if gcInterval > 0 {
		interval = gcInterval
	}
	store.ticker = time.NewTicker(time.Second * time.Duration(interval))

	table := store.db.AddTableWithName(models.Token{}, store.tableName)
	table.AddIndex("idx_code", "Btree", []string{"Code"})
	table.AddIndex("idx_access", "Btree", []string{"Access"})
	table.AddIndex("idx_refresh", "Btree", []string{"Refresh"})
	table.AddIndex("idx_clientId", "Btree", []string{"ClientID"})
	table.AddIndex("idx_userName", "Btree", []string{"UserID"})

	err := store.db.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}

	// 由于mysql_gorm.go 默认varchar 255字符， 唯一索引过长
	store.db.Exec(fmt.Sprintf("alter table %s modify column  ClientID varchar(128)", store.tableName))
	store.db.Exec(fmt.Sprintf("alter table %s modify column  UserID varchar(200)", store.tableName))
	//idx := table.AddIndex("idx_UserClient", "Btree", []string{"UserID", "ClientID"})
	//idx.Unique = true
	//
	//store.db.CreateIndex() // 只有新建表的时候 有作用,重复建立索引，报错中止

	go store.gc()
	return store
}

// Store mysql token store
type Store struct {
	tableName string
	db        *gorp.DbMap
	stdout    io.Writer
	ticker    *time.Ticker
}
type StoreEntry struct {
	*Store
}

// SetStdout set error output
func (s *Store) SetStdout(stdout io.Writer) *Store {
	s.stdout = stdout
	return s
}

// Close close the store
func (s *Store) Close() {
	s.ticker.Stop()
	s.db.Db.Close()
}

func (s *Store) gc() {
	for range s.ticker.C {
		s.clean()
	}
}

func (s *Store) clean() {
	now := time.Now().Unix()
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE UNIX_TIMESTAMP(AccessCreateAt)+AccessExpiresIn/1000000000<=? OR (Code='' AND Access='' AND Refresh='')", s.tableName)
	n, err := s.db.SelectInt(query, now)
	if err != nil || n == 0 {
		if err != nil {
			s.errorf(err.Error())
		}
		return
	}

	result, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE UNIX_TIMESTAMP(AccessCreateAt)+AccessExpiresIn/1000000000<=? OR (Code='' AND Access='' AND Refresh='')", s.tableName), now)
	if err != nil {
		s.errorf(err.Error())
	}
	rn, _ := result.RowsAffected()
	if rn > 0 && g.GetConfig().IsDebug() {
		g.GetLog().Debug("TokenStore Clean GC Row: %d\n", rn)
	}

}

func (s *Store) errorf(format string, args ...interface{}) {
	if s.stdout != nil {
		buf := fmt.Sprintf("[OAUTH2-MYSQL-ERROR]: "+format, args...)
		s.stdout.Write([]byte(buf))
	}
}

// Create create and store the new token information
func (s *Store) Create(ctx context.Context, info oauth2.TokenInfo) error {
	//buf, _ := jsoniter.Marshal(info)
	//item := &StoreItem{
	//	Data: string(buf),
	//}
	//fmt.Printf("%+v\n",info.(*models.Token))
	//item.UserName = info.GetUserID()
	//item.ClientID = info.GetClientID()
	//
	//if code := info.GetCode(); code != "" {
	//	item.Code = code
	//	item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	//} else {
	//	item.Access = info.GetAccess()
	//	item.AccessExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()
	//
	//	if refresh := info.GetRefresh(); refresh != "" {
	//		item.Refresh = info.GetRefresh()
	//		item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
	//	}
	//
	//}

	return s.db.Insert(info.(*models.Token))
}

// RemoveByCode delete the authorization code
func (s *Store) RemoveByCode(ctx context.Context, code string) error {
	query := fmt.Sprintf("UPDATE %s SET Code='' WHERE Code=? LIMIT 1", s.tableName)
	_, err := s.db.Exec(query, code)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

// RemoveByAccess use the access token to delete the token information
func (s *Store) RemoveByAccess(ctx context.Context, access string) error {
	query := fmt.Sprintf("UPDATE %s SET Access='', ClientID='' ,UserID=''  WHERE Access=? LIMIT 1", s.tableName)
	_, err := s.db.Exec(query, access)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

// RemoveByRefresh use the refresh token to delete the token information
func (s *Store) RemoveByRefresh(ctx context.Context, refresh string) error {
	query := fmt.Sprintf("UPDATE %s SET Refresh='' , ClientID='' ,UserID='' WHERE Refresh=? LIMIT 1", s.tableName)
	_, err := s.db.Exec(query, refresh)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

//func (s *Store) toTokenInfo(data string) oauth2.TokenInfo {
//	var tm models.Token
//	jsoniter.Unmarshal([]byte(data), &tm)
//	return &tm
//}

// GetByCode use the authorization code for token information data
func (s *Store) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	if code == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE Code=? LIMIT 1", s.tableName)
	var item models.Token
	err := s.db.SelectOne(&item, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	//return s.toTokenInfo(item.Data), nil
	return &item, err
}

// GetByAccess use the access token for token information data
func (s *Store) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	if access == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE Access=? LIMIT 1", s.tableName)
	var item models.Token
	err := s.db.SelectOne(&item, query, access)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	//return s.toTokenInfo(item.Data), nil
	return &item, err
}

// GetByRefresh use the refresh token for token information data
func (s *Store) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	if refresh == "" {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE refresh=? LIMIT 1", s.tableName)
	var item models.Token
	err := s.db.SelectOne(&item, query, refresh)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	//return s.toTokenInfo(item.Data), nil
	return &item, err
}

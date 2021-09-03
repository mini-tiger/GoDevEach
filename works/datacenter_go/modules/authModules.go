package modules

import "datacenter/g"

/**
 * @Author: Tao Jun
 * @Description: modules
 * @File:  authModules
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午4:16
 */

type User struct {
	Id         int64  `json:"id" gorm:"primary_key"`
	Name       string `json:"username"`
	Password   string `json:"password"`
	Roles      string `json:"roles"`
	Status     string `json:"status"`
	Source     string `json:"source"`
	UpdateTime string `json:"update_time" gorm:"->"` // 只读
	Client     string `json:"client"`
}

func (User) TableName() string {
	return "user"
}

type OauthClientDetails struct {
	AccessTokenValidity   int    `json:"access_token_validity" gorm:"type:int(11)" `
	RefreshTokenValidity  int    `json:"refresh_token_validity" gorm:"type:int(11)"`
	ClientId              string `json:"client_id" gorm:"primary_key;type:varchar(250)"`
	ResourceIds           string `json:"resource_ids"`
	ClientSecret          string `json:"client_secret"`
	Scope                 string `json:"scope"`
	AuthorizedGrantTypes  string `json:"authorized_grant_types"`
	WebServerRedirectUri  string `json:"web_server_redirect_uri"`
	Authorities           string `json:"authorities"`
	AdditionalInformation string `json:"additional_information"`
	Autoapprove           string `json:"autoapprove"`
}

func (OauthClientDetails) TableName() string {
	//return "oauth_client_details"
	return g.GetConfig().ClientTableName
}

// GetID client id
func (c *OauthClientDetails) GetID() string {
	return c.ClientId
}

// GetSecret client domain
func (c *OauthClientDetails) GetSecret() string {
	return c.ClientSecret
}

// GetDomain client domain
func (c *OauthClientDetails) GetDomain() string {
	return c.WebServerRedirectUri
}

// GetUserID user id  默认用户 与 clientid 一样
func (c *OauthClientDetails) GetUserID() string {
	return c.ClientId
}

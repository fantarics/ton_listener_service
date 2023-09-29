package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

func (Transaction) TableName() string {
	return "Transactions"
}

type Transaction struct {
	TxHash             string `gorm:"primaryKey;column:tx_hash;not null"`
	Hash               string `gorm:"column:hash"`
	DestinationAddress string `gorm:"column:destination_address"`
	WalletID           string `gorm:"column:wallet_id;foreignKey:WalletAddress"`
	UserID             uint64 `gorm:"column:user_id;foreignKey:UserID"`

	IncomingValue   uint64    `gorm:"column:incoming_value"`
	SendedAmount    uint64    `gorm:"column:sended_amount;default:0;"`
	TransactionType string    `gorm:"column:transaction_type"`
	CreatedTime     time.Time `gorm:"column:created_time"`
}

func (Network) TableName() string {
	return "Network"
}

type Network struct {
	ID        uint64 `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	ShortName string `gorm:"column:short_name"`

	Coins []Coin `gorm:"foreignKey:NetworkID"`
}

func (Coin) TableName() string {
	return "Coin"
}

type Coin struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name"`
	ShortName   string    `gorm:"short_name"`
	Digits      uint64    `gorm:"column:digits"`
	Price       float64   `gorm:"column:price"`
	Image       string    `gorm:"column:image"`
	LastUpdated time.Time `gorm:"column:last_updated"`
	NetworkID   uint64    `gorm:"column:network_id;foreignKey:NetworkID"`

	Balance []Balance `gorm:"foreignKey:CoinID"`

	Network Network `gorm:"foreignKey:NetworkID"`
}

func (Balance) TableName() string {
	return "Balance"
}

type Balance struct {
	ID        uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint64 `gorm:"column:user_id;foreignKey:UserID"`
	NetworkID uint64 `gorm:"column:network_id;foreignKey:NetworkID"`
	CoinID    uint64 `gorm:"column:coin_id;foreignKey:CoinID"`
	Balance   uint64 `gorm:"column:balance"`

	User User `gorm:"foreignKey:UserID"`
	Coin Coin `gorm:"foreignKey:CoinID"`
}

func (Log) TableName() string {
	return "Log"
}

type Log struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Action      string    `gorm:"column:string"`
	JSONObject  []byte    `gorm:"type:jsonb;column:json_object"`
	CreatedDate time.Time `gorm:"column:created_date"`

	UserID uint64 `gorm:"column:user_id;foreignKey:UserID"`
	User   User   `gorm:"foreignKey:UserID"`
}

func (Wallet) TableName() string {
	return "Wallet"
}

type Wallet struct {
	Address          string `gorm:"primaryKey;unique"`
	UserID           uint64 `gorm:"column:user_id;foreignKey:UserID"`
	MinWithdrawLimit uint64 `gorm:"column:min_withdraw_limit"`
	WithdrawalFee    uint64 `gorm:"column:min_withdraw_limit"`

	NetworkID uint64 `gorm:"column:network_id;foreignKey:NetworkID"`
	CoinID    uint64 `gorm:"column:coin_id;foreignKey:CoinID"`

	PublicKey   string         `gorm:"column:min_withdraw_limit"`
	PrivateKey  string         `gorm:"column:min_withdraw_limit"`
	Mnemonics   pq.StringArray `gorm:"column:mnemonics;type:varchar(255)[]"`
	IsActive    bool           `gorm:"column:is_active"`
	CreatedDate time.Time      `gorm:"column:created_date"`
	UpdatedDate time.Time      `gorm:"column:updated_date"`

	User User `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "User"
}

type User struct {
	ID           uint64         `gorm:"column:id"`
	Username     string         `gorm:"column:username"`
	FirstName    string         `gorm:"column:first_name"`
	ReferralCode string         `gorm:"column:referral_code"`
	InvitedBy    sql.NullString `gorm:"column:invited_by"`
	IsAdmin      sql.NullBool   `gorm:"column:is_admin"`
	IsVerified   sql.NullBool   `gorm:"column:is_verified"`
	IsActive     sql.NullBool   `gorm:"column:is_active"`
	NativeLang   string         `gorm:"column:native_lang"`
	BotLang      sql.NullString `gorm:"column:bot_lang"`
	CreatedDate  time.Time      `gorm:"column:created_date"`

	Wallets  []Wallet  `gorm:"foreignKey:UserID"`
	Balances []Balance `gorm:"foreignKey:UserID"`
	Logs     []Log     `gorm:"foreignKey:UserID"`
}

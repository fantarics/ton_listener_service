package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetWallets(coin string) (wallets []Wallet, err error)
	GetCoin(name string) (coin Coin, err error)
	CreateTrx(trx *Transaction) (err error)
	AddBalance(userID, coinID uint64, amount int64) (err error)
	CreateLog(log *Log) (err error)
	GetUser(userId uint64) (user User, err error)
}

type dbSQL struct {
	db *gorm.DB
}

func NewSQL(db *gorm.DB) Repository {
	return &dbSQL{db: db}
}

func (db *dbSQL) GetWallets(name string) (wallets []Wallet, err error) {

	coin, err := db.GetCoin(name)
	if err != nil {
		return
	}

	err = db.db.Model(&Wallet{CoinID: coin.ID}).Preload(clause.Associations).Find(&wallets).Error
	if err != nil {
		return
	}
	return wallets, nil
}

func (db *dbSQL) GetCoin(name string) (coin Coin, err error) {

	err = db.db.Model(&Coin{ShortName: name}).First(&coin).Error
	if err != nil {
		return
	}

	return coin, nil
}

func (db *dbSQL) CreateTrx(trx *Transaction) (err error) {
	tx := db.db.Begin()
	err = tx.Create(trx).Error
	if err != nil {
		return err
	}
	return tx.Commit().Error
}

func (db *dbSQL) AddBalance(userID, coinID uint64, amount int64) (err error) {
	err = db.db.Model(&Balance{}).Where("user_id = ?", userID).Where("coin_id = ?", coinID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return
	}
	return nil
}

func (db *dbSQL) CreateLog(log *Log) (err error) {
	tx := db.db.Begin()
	err = tx.Create(log).Error
	if err != nil {
		return
	}
	return tx.Commit().Error

}

func (db *dbSQL) GetUser(userId uint64) (user User, err error) {
	err = db.db.Model(&User{ID: userId}).First(&user).Error
	if err != nil {
		return
	}
	return user, err
}

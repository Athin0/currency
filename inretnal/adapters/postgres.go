package adapters

import (
	"context"
	"currency/inretnal/model"
	"database/sql"
	_ "encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type PostgresDB struct {
	Client *sql.DB
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres connect error : (%v)", err)
	}
	fmt.Println(db)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{Client: db}, nil
}

func InitDB() (*PostgresDB, error) {
	viper.AddConfigPath("../currency/db") //change
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("error in reading config: %v", err)
		return nil, err
	}
	db, err := NewPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("error creating db: %v \n", err)
		return nil, err
	}
	return db, nil
}

func (db *PostgresDB) InsertCurrencies(ctx context.Context, currencies []model.Currency, time time.Time) error {
	tx, err := db.Client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	for _, currency := range currencies {
		value := strings.Replace(currency.Value, ",", ".", -1)
		_, err = tx.Exec(
			"INSERT INTO Currencies(ID,NumCode, CharCode,Nom, Name, Value,Date ) values ($1, $2, $3,$4, $5,$6,$7)",
			currency.ID, currency.NumCode, currency.CharCode, currency.Nom, currency.Name, value, time,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (db *PostgresDB) GetLastDate(ctx context.Context) (time.Time, error) {
	var t time.Time
	q := db.Client.QueryRowContext(ctx,
		"select max(date) from currencies")
	if q.Err() != nil {
		return time.Time{}, q.Err()
	}
	err := q.Scan(&t)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (db *PostgresDB) GetMax(ctx context.Context) (*model.ResponseCurrency, error) {
	err := db.Client.QueryRowContext(ctx,
		`SELECT Name, Value, Date
			FROM Currencies
			WHERE Date >= NOW() - INTERVAL '90 days'
			ORDER BY Value DESC, Date ASC
			LIMIT 1;`,
	)
	if err.Err() != nil {
		log.Printf(err.Err().Error())
		return nil, err.Err()
	}
	var ans model.ResponseCurrency
	err.Scan(&ans.Name, &ans.Value, &ans.Date)
	return &ans, nil
}

func (db *PostgresDB) GetMin(ctx context.Context) (*model.ResponseCurrency, error) {
	err := db.Client.QueryRowContext(ctx,
		`SELECT Name, Value, Date
			FROM Currencies
			WHERE Date >= NOW() - INTERVAL '90 days'
			ORDER BY Value ASC, Date ASC
			LIMIT 1;`,
	)
	if err.Err() != nil {
		log.Printf(err.Err().Error())
		return nil, err.Err()
	}
	var ans model.ResponseCurrency
	err.Scan(&ans.Name, &ans.Value, &ans.Date)
	return &ans, nil
}

func (db *PostgresDB) GetAvg(ctx context.Context) ([]model.ResponseCurrencyAvg, error) {
	rows, err := db.Client.QueryContext(ctx,
		`SELECT Name, avg(Value)
				FROM Currencies
				WHERE Date >= NOW() - INTERVAL '90 days'
				GROUP BY Name;`,
	)
	if err != nil {
		return nil, err
	}
	var arr []model.ResponseCurrencyAvg
	for rows.Next() {
		var ans model.ResponseCurrencyAvg
		rows.Scan(&ans.Name, &ans.Value)
		arr = append(arr, ans)
	}
	return arr, nil
}

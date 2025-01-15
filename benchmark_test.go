package benchmark

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-juicedev/juice"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"testing"
)

var (
	engine *juice.Engine
	gormDB *gorm.DB
	db     *sql.DB
)

func init() {
	setupTestDB()
}

func setupTestDB() {
	var err error
	cfg, err := juice.NewXMLConfiguration("juice.xml")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	engine, err = juice.New(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create juice engine: %v", err))
	}
	db = engine.DB()
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create gorm db: %v", err))
	}
	initTestTable()
}

func truncateAndRecreateTable() error {
	_, err := db.Exec("DROP TABLE IF EXISTS `tbl_user`")
	if err != nil {
		return err
	}

	query := `
CREATE TABLE tbl_user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_age (age),
    INDEX idx_email (email)
);`
	if _, err = db.Exec(query); err != nil {
		return err
	}
	return nil
}

func initTestTable() {
	if err := truncateAndRecreateTable(); err != nil {
		panic(err)
	}
}

func prepareTestData(b *testing.B, count int) {
	b.Helper()
	for i := 0; i < count; i++ {
		user := &JuiceUser{
			Name:  "test" + strconv.Itoa(i),
			Age:   18,
			Email: "test" + strconv.Itoa(i) + "@example.com",
		}
		_, err := db.Exec(
			"INSERT INTO tbl_user(`name`, `age`, `email`) VALUES (?,?,?)",
			user.Name, user.Age, user.Email,
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserCreate(b *testing.B) {
	b.Run("STD_DB", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			user := &JuiceUser{
				Name:  "test" + strconv.Itoa(n),
				Age:   18,
				Email: "test" + strconv.Itoa(n) + "@example.com",
			}
			query := "INSERT INTO tbl_user(`name`, `age`, `email`) VALUES (?,?,?)"
			result, err := db.Exec(query, user.Name, user.Age, user.Email)
			if err != nil {
				b.Fatal(err)
			}

			id, err := result.LastInsertId()
			if err != nil {
				b.Fatal(err)
			}
			if id != int64(n+1) {
				b.Fatalf("expected id %d, got %d", n+1, id)
			}
		}
	})

	b.Run("Juice", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		ctx := juice.ContextWithManager(context.Background(), engine)
		userRepo := NewUserRepository()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			user := &JuiceUser{
				Name:  "test" + strconv.Itoa(n),
				Age:   18,
				Email: "test" + strconv.Itoa(n) + "@example.com",
			}
			_, err := userRepo.Create(ctx, user)
			if err != nil {
				b.Fatal(err)
			}
			if user.ID != n+1 {
				b.Fatalf("expected id %d, got %d", n+1, user.ID)
			}
		}
	})

	b.Run("GORM", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			user := &GormUser{
				Name:  "test" + strconv.Itoa(n),
				Age:   18,
				Email: "test" + strconv.Itoa(n) + "@example.com",
			}
			if err := gormDB.Create(user).Error; err != nil {
				b.Fatal(err)
			}
			if user.ID != uint(n+1) {
				b.Fatalf("expected id %d, got %d", n+1, user.ID)
			}
		}
	})
}

func timerScope[T any](b *testing.B, f func() T) T {
	b.StopTimer()
	defer b.StartTimer()
	return f()
}

func BenchmarkBatchCreate(b *testing.B) {

	b.Run("STD_DB", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {

			values := strings.Repeat("(?,?,?),", 1000)
			values = values[:len(values)-1]

			query := "INSERT INTO tbl_user(`name`, `age`, `email`) VALUES " + values

			insertValues := timerScope(b, func() []interface{} {
				var users []JuiceUser
				items := make([]interface{}, 0, 3000)
				for i := 0; i < 1000; i++ {
					users = append(users, JuiceUser{
						Name:  "test" + strconv.Itoa(i),
						Age:   18,
						Email: "test" + strconv.Itoa(i) + "@example.com",
					})
					items = append(items, users[i].Name, users[i].Age, users[i].Email)
				}
				return items
			})

			result, err := db.Exec(query, insertValues...)
			if err != nil {
				b.Fatal(err)
			}
			_, err = result.LastInsertId()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Juice", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		ctx := juice.ContextWithManager(context.Background(), engine)
		userRepo := NewUserRepository()

		b.ResetTimer()

		for n := 0; n < b.N; n++ {

			users := timerScope(b, func() []*JuiceUser {
				var users = make([]*JuiceUser, 0, 1000)
				for i := 0; i < 1000; i++ {
					users = append(users, &JuiceUser{
						Name:  "test" + strconv.Itoa(i),
						Age:   18,
						Email: "test" + strconv.Itoa(i) + "@example.com",
					})
				}
				return users
			})

			_, err := userRepo.BatchCreate(ctx, users)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GORM", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {

			users := timerScope(b, func() []*GormUser {
				var users = make([]*GormUser, 0, 1000)
				for i := 0; i < 1000; i++ {
					users = append(users, &GormUser{
						Name:  "test" + strconv.Itoa(i),
						Age:   18,
						Email: "test" + strconv.Itoa(i) + "@example.com",
					})
				}
				return users
			})

			if err := gormDB.Create(users).Error; err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUserQueryAll(b *testing.B) {
	const dataCount = 1000

	b.Run("STD_DB", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		prepareTestData(b, dataCount)

		query := "SELECT * FROM tbl_user"
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			rows, err := db.Query(query)
			if err != nil {
				b.Fatal(err)
			}
			users := make([]JuiceUser, 0, dataCount)
			for rows.Next() {
				var user JuiceUser
				if err := rows.Scan(
					&user.ID, &user.Name, &user.Age, &user.Email,
					&user.CreatedAt, &user.UpdatedAt,
				); err != nil {
					rows.Close()
					b.Fatal(err)
				}
				users = append(users, user)
			}
			rows.Close()
			if len(users) != dataCount {
				b.Fatalf("expected %d users, got %d", dataCount, len(users))
			}
		}
	})

	b.Run("Juice", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		prepareTestData(b, dataCount)

		ctx := juice.ContextWithManager(context.Background(), engine)
		userRepo := NewUserRepository()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			users, err := userRepo.QueryAll(ctx)
			if err != nil {
				b.Fatal(err)
			}
			if len(users) != dataCount {
				b.Fatalf("expected %d users, got %d", dataCount, len(users))
			}
		}
	})

	b.Run("GORM", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		prepareTestData(b, dataCount)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			var users []GormUser
			if err := gormDB.Find(&users).Error; err != nil {
				b.Fatal(err)
			}
			if len(users) != dataCount {
				b.Fatalf("expected %d users, got %d", dataCount, len(users))
			}
		}
	})
}

func BenchmarkUserQueryWithLimit(b *testing.B) {
	const dataCount = 1000

	b.Run("STD_DB", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		prepareTestData(b, dataCount)

		query := "SELECT * FROM tbl_user LIMIT ?"
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			rows, err := db.Query(query, dataCount)
			if err != nil {
				b.Fatal(err)
			}
			users := make([]JuiceUser, 0, dataCount)
			for rows.Next() {
				var user JuiceUser
				if err := rows.Scan(
					&user.ID, &user.Name, &user.Age, &user.Email,
					&user.CreatedAt, &user.UpdatedAt,
				); err != nil {
					rows.Close()
					b.Fatal(err)
				}
				users = append(users, user)
			}
			rows.Close()
			if len(users) != dataCount {
				b.Fatalf("expected %d users, got %d", dataCount, len(users))
			}
		}
	})

	b.Run("Juice", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		prepareTestData(b, dataCount)

		ctx := juice.ContextWithManager(context.Background(), engine)
		userRepo := NewUserRepository()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			users, err := userRepo.QueryWithLimit(ctx, dataCount)
			if err != nil {
				b.Fatal(err)
			}
			if len(users) != dataCount {
				b.Fatalf("expected %d users, got %d", dataCount, len(users))
			}
		}
	})

	b.Run("GORM", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}
		prepareTestData(b, dataCount)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			var users = make([]GormUser, 0, dataCount)
			if err := gormDB.Limit(dataCount).Find(&users).Error; err != nil {
				b.Fatal(err)
			}
			if len(users) != dataCount {
				b.Fatalf("expected %d users, got %d", dataCount, len(users))
			}
		}
	})
}

// run with batchSize

func BenchmarkUserBatchCreate(b *testing.B) {
	const batchSize = 100
	const times = 10

	b.Run("STD_DB", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		users := make([]JuiceUser, 0, batchSize*times)

		for i := range batchSize * times {
			users = append(users, JuiceUser{
				Name:  "test" + strconv.Itoa(i),
				Age:   18,
				Email: "test" + strconv.Itoa(i) + "@example.com",
			})
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			for i := 0; i < times; i++ {
				query := "INSERT INTO tbl_user(`name`, `age`, `email`) VALUES " + strings.Repeat("(?,?,?),", batchSize)
				query = query[:len(query)-1]
				values := make([]interface{}, 0, batchSize)
				for j := 0; j < batchSize; j++ {
					values = append(values, users[i*batchSize+j].Name, users[i*batchSize+j].Age, users[i*batchSize+j].Email)
				}
				if _, err := db.Exec(query, values...); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("Juice", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		ctx := juice.ContextWithManager(context.Background(), engine)
		userRepo := NewUserRepository()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {

			users := timerScope(b, func() []*JuiceUser {
				users := make([]*JuiceUser, 0, batchSize*times)

				for i := range batchSize * times {
					users = append(users, &JuiceUser{
						Name:  "test" + strconv.Itoa(i),
						Age:   18,
						Email: "test" + strconv.Itoa(i) + "@example.com",
					})
				}
				return users
			})

			if _, err := userRepo.BatchCreateWithBatchSize(ctx, users); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GORM", func(b *testing.B) {
		if err := truncateAndRecreateTable(); err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			users := timerScope(b, func() []*GormUser {
				var users = make([]*GormUser, 0, 1000)
				for i := 0; i < 1000; i++ {
					users = append(users, &GormUser{
						Name:  "test" + strconv.Itoa(i),
						Age:   18,
						Email: "test" + strconv.Itoa(i) + "@example.com",
					})
				}
				return users
			})

			if err := gormDB.CreateInBatches(users, batchSize).Error; err != nil {
				b.Fatal(err)
			}
		}
	})
}

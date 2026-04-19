package cmd

import (
	"fmt"
	apiModels "github.com/mmtaee/ocserv-dashboard/api/internal/models"
	"github.com/mmtaee/ocserv-dashboard/common/models"
	"github.com/mmtaee/ocserv-dashboard/common/pkg/config"
	"github.com/mmtaee/ocserv-dashboard/common/pkg/database"
	"github.com/mmtaee/ocserv-dashboard/common/pkg/logger"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"reflect"
)

var dbLoaderCmd = &cobra.Command{
	Use:   "db-loader",
	Short: "Run the db-loader to load old Sqlite to PostgreSQL",
	Run: func(cmd *cobra.Command, args []string) {
		loader()
	},
}

func init() {
	rootCmd.AddCommand(dbLoaderCmd)
}

func loader() {
	config.Init(false, "", 0)

	database.ConnectPostgres()
	defer database.ClosePostgres()
	pgDB := database.GetPostgres()

	database.ConnectSQLite(true)
	defer database.CloseSQLite()
	sqliteDB := database.GetSQLite()

	var oldSqliteModels = []interface{}{
		// system
		&apiModels.System{},

		// users
		&apiModels.User{},
		&apiModels.UserToken{},

		// ocserv
		&models.OcservGroup{},
		&models.OcservUser{},
		&models.OcservUserTrafficStatistics{},
	}

	for _, model := range oldSqliteModels {
		tableName := getTableName(pgDB, model)
		logger.Info("Start migrating: %v", tableName)

		// 1. Create schema
		if err := pgDB.AutoMigrate(model); err != nil {
			logger.Fatal("error migrating: %v", err)
		}

		// 2. Copy data
		if err := migrateTable(sqliteDB, pgDB, model); err != nil {
			logger.Fatal("error migrating: %v", err)
		}

		// 3. Fix sequence (PostgreSQL)
		if err := resetSequence(pgDB, model); err != nil {
			logger.Error("sequence warning: %v", err)
		}

		logger.Info("Migration for table (%s) complete", tableName)
	}

}

func migrateTable(sqliteDB, pgDB *gorm.DB, model interface{}) error {
	batchSize := 100

	slicePtr := createSlice(model)

	return sqliteDB.Model(model).
		FindInBatches(slicePtr, batchSize, func(tx *gorm.DB, batch int) error {

			data := tx.Statement.Dest

			if err := pgDB.Create(data).Error; err != nil {
				return err
			}

			fmt.Printf("  → batch %d inserted\n", batch)
			return nil
		}).Error
}

func createSlice(model interface{}) interface{} {
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(modelType)
	return reflect.New(sliceType).Interface()
}

func getTableName(db *gorm.DB, model interface{}) string {
	stmt := &gorm.Statement{DB: db}
	_ = stmt.Parse(model)
	return stmt.Schema.Table
}

func resetSequence(db *gorm.DB, model interface{}) error {
	table := getTableName(db, model)

	query := fmt.Sprintf(`
		SELECT setval(
			pg_get_serial_sequence('%s', 'id'),
			COALESCE(MAX(id), 1),
			true
		) FROM %s;
	`, table, table)

	return db.Exec(query).Error
}

package repository

import (
	"errors"
	"fmt"
	"time"
	"wereserve/models"

	"gorm.io/gorm"
)

type TableRepository struct {
	DB *gorm.DB
}


func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{DB: db}
}


func (r *TableRepository) IsTableExists(tableName string) (bool, error) {
	var tableExists bool
	err := r.DB.Raw("SELECT EXISTS(SELECT 1 FROM tables WHERE table_name = ?)",tableName).Scan(&tableExists).Error
	if err != nil {
		return false, err
	}
	return tableExists, nil
}

func (r *TableRepository) GetAllTables() ([]models.Table, error) {
	var tables []models.Table
	err := r.DB.Raw("SELECT * FROM tables").Scan(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
} 

func (r *TableRepository) GetTableByID(id int) (*models.Table, error) {
	var table models.Table
    err := r.DB.First(&table, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("table with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to fetch table with ID %d: %w", id, err)
    }
    return &table, nil
}

func (r *TableRepository) GetTableByStatus(status string) ([]models.Table, error) {
	var tables []models.Table
	// err := r.DB.Raw("SELECT * FROM tables WHERE status = $1", status).Scan(&tables).Error
	err := r.DB.Where("status = ?",status).Find(&tables).Error
	if err != nil {
		// Handle error database
		return nil, fmt.Errorf("failed to fetch tables: %w", err)
	}

	return tables, nil 
}


func (r *TableRepository) CreateTable(table *models.Table) error {
	result := r.DB.FirstOrCreate(&table, models.Table{TableName: table.TableName})
	if result.Error != nil {
		return result.Error
	} 

	if result.RowsAffected == 0 {
		return errors.New("meja dengan nomor tersebut telah digunakan")
	}

	return nil
}

func (r *TableRepository) DeleteTable(id int) error {
	sqlQuery := `DELETE FROM tables WHERE id = $1`
	result := r.DB.Exec(sqlQuery, id)

	if result.Error != nil {
		return result.Error
	}

	// validate if err
	if result.RowsAffected == 0 {
		return fmt.Errorf("tables with id %d not found", id)
	}

	return nil
}

// update Table
func (r *TableRepository) UpdateTable(id int, table *models.Table) error {
	var currentTable models.Table
	err := r.DB.First(&currentTable, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("table tidak di temukan")
		}
		return err
	}

	// Validasi table
	if table.TableName != "" && table.TableName != currentTable.TableName {
		tableExists, err := r.IsTableExists(table.TableName)
		if err != nil {
			return err
		}

		if tableExists {
			return errors.New("table name sudah terdaftar")
		}
	}

	// update filed yang diisi
	updates := make(map[string]interface{})
	if table.TableName != "" {
		updates["table_name"] = table.TableName
	}

	if table.Capacity != 0 {
		updates["capacity"] = table.Capacity
	}

	if table.Status != "" {
		updates["status"] = table.Status
	}
	updates["updated_at"] = time.Now()

	// eksekusi query update
	err = r.DB.Model(&models.Table{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
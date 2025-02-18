package services

import (
	"errors"
	"wereserve/models"
	"wereserve/repository"
)

type TableService struct {
	tableRepo *repository.TableRepository
}

func NewTableService(tableRepo *repository.TableRepository) *TableService {
	return &TableService{tableRepo: tableRepo}
}


// get All Table
func (s *TableService) GetAllTable() ([]models.Table, error) {
	tables, err := s.tableRepo.GetAllTables()
	if err != nil {
		return nil, err
	}
	return tables, nil
}

// get UserById
func (s *TableService) GetTableById(id int) (*models.Table, error) {
	table, err := s.tableRepo.GetTableByID(id)
	if err != nil {
		return nil, err
	}

	return table, nil
}

// get UserByStatus
func (s *TableService) GetTableByStatus(status string) ([]models.Table, error) {
	tables, err := s.tableRepo.GetTableByStatus(status)
		if err != nil {
			return nil, err
		}
	return tables, nil
}


// Create Table 
func (s *TableService) CreateTable(table *models.Table) error {
	//Cek Table Number Exist
	_, err := s.tableRepo.IsTableExists(table.TableName)
	if err != nil {
		return err
	}

	err = s.tableRepo.CreateTable(table)
	if err != nil {
		return err
	}

	return nil
}

// Create Delete Table
func (s *TableService) DeleteTable(id int) error {
	//Cek Table Exist
	_, err := s.tableRepo.GetTableByID(id)
	if err != nil {
		return errors.New("table id tidak di temukan")
	}

	//Delete
	err = s.DeleteTable(id)
	if err != nil {
		return err
	}

	return nil
}


func (s *TableService) UpdateTable(id int, table models.Table) error {
	//Validasi Input min 1 field yang disi
	if table.TableName == "" && table.Capacity == 0 && table.Status == "" {
		return errors.New("minimal satu field harus diisi")
	}

	err := s.tableRepo.UpdateTable(id, &models.Table{
		TableName: table.TableName,
		Capacity:  table.Capacity,
		Status:    table.Status,
	})

	if err != nil {
		return err
	}

	return nil
}
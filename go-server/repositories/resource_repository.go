package repositories

import (
	"fmt"
	"server/database"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResourceRepository struct{}

func (r *ResourceRepository) GetAll(resource string) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(fmt.Sprintf("SELECT * FROM %s", resource))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	result := []map[string]interface{}{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			// Handling nil
			if val != nil {
				rowMap[colName] = *val
			} else {
				rowMap[colName] = nil
			}
		}
		result = append(result, rowMap)
	}
	return result, nil
}

func (r *ResourceRepository) Create(resource string, data map[string]interface{}) (int, error) {
	keys := []string{}
	values := []interface{}{}
	placeholders := []string{}
	i := 1
	for k, v := range data {
		keys = append(keys, k)
		values = append(values, v)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", resource, strings.Join(keys, ", "), strings.Join(placeholders, ", "))
	var newID int
	err := database.DB.QueryRow(query, values...).Scan(&newID)
	return newID, err
}

func (r *ResourceRepository) Update(resource string, id string, data map[string]interface{}) error {
	keys := []string{}
	values := []interface{}{}
	i := 1
	for k, v := range data {
		keys = append(keys, fmt.Sprintf("%s = $%d", k, i))
		values = append(values, v)
		i++
	}
	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", resource, strings.Join(keys, ", "), i)
	_, err := database.DB.Exec(query, values...)
	return err
}

func (r *ResourceRepository) Delete(resource string, id string) error {
	_, err := database.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = $1", resource), id)
	return err
}

func (r *ResourceRepository) GetEmployeeNames() ([]gin.H, error) {
	rows, err := database.DB.Query("SELECT id, name FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var empList []gin.H
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err == nil {
			empList = append(empList, gin.H{"id": id, "name": name})
		}
	}
	return empList, nil
}

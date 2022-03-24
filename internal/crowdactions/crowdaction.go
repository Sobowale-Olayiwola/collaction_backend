package crowdaction

import (
	"context"
	"fmt"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
)

type CrowdactionRepository interface {
	GetDBItem(tableName string, pk string, crowdactionId string) (*models.CrowdactionData, error)
	Query(tableName string, filterCond string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error)
}
type Service interface {
	GetCrowdactionById(ctx context.Context, crowdactionId string) (*models.CrowdactionData, error)
	GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error)
}
type crowdaction struct {
	dynamodb CrowdactionRepository
}

const (
	KeyDateStart      = "date_start"
	KeyDateEnd        = "date_end"
	KeyDateJoinBefore = "date_limit_join"
)

func NewCrowdactionService(dynamodb CrowdactionRepository) Service {
	return &crowdaction{dynamodb: dynamodb}
}

/**
	GET Crowdaction by Id
**/
func (e *crowdaction) GetCrowdactionById(ctx context.Context, crowdactionID string) (*models.CrowdactionData, error) {
	fmt.Println("GetCrowdaction calling internal:", crowdactionID)
	item, err := e.dynamodb.GetDBItem(constants.TableName, utils.PKCrowdaction, crowdactionID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("crowdaction not found")
	}
	return item, err
}

/**
	GET Crowdaction by Status
**/
func (e *crowdaction) GetCrowdactionsByStatus(ctx context.Context, status string, startFrom *utils.PrimaryKey) ([]models.CrowdactionData, error) {
	crowdactions := []models.CrowdactionData{} // empty crowdaction array

	switch status {
	case "joinable":
		items, err := e.dynamodb.Query(constants.TableName, KeyDateJoinBefore, startFrom)
		return items, err
	case "active":
		items, err := e.dynamodb.Query(constants.TableName, KeyDateStart, startFrom)
		return items, err
	case "ended":
		items, err := e.dynamodb.Query(constants.TableName, KeyDateEnd, startFrom)
		return items, err
	default:
		return crowdactions, nil
	}
}

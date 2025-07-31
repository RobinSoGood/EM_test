package service

import (
	"github.com/RobinSoGood/EM_test/internal/models"

	"github.com/rs/zerolog/log"
)

type SubStorage interface {
	SaveSub(models.Sub) (string, error)
	GetSubs() ([]models.Sub, error)
	GetSub(string) (models.Sub, error)
	SetDeleteSubStatus(string) error
	GetTotalPriceByPeriod(models.SumRequest) (int, error)
}
type SubService struct {
	stor SubStorage
}

func NewSubService(stor SubStorage) SubService {
	return SubService{stor: stor}
}

func (ss *SubService) AddSub(sub models.Sub) (string, error) {
	sid, err := ss.stor.SaveSub(sub)
	if err != nil {
		log.Error().Err(err).Msg("save sub failed")
		return ``, err
	}
	return sid, nil
}

func (ss *SubService) GetSub(sid string) (models.Sub, error) {
	return ss.stor.GetSub(sid)
}

func (ss *SubService) GetSubs() ([]models.Sub, error) {
	return ss.stor.GetSubs()
}

func (ss *SubService) SetDeleteStatus(sid string) error {
	return ss.stor.SetDeleteSubStatus(sid)
}

func (ss *SubService) GetTotalPriceByPeriod(sum models.SumRequest) (int, error) {
	total, err := ss.stor.GetTotalPriceByPeriod(sum)
	if err != nil {
		log.Error().Err(err).Msg("failed to get total price by period")
		return 0, err
	}
	return total, nil
}

package dal

import (
	"context"

	"github.com/phuslu/log"
)

func GetAllBonus(ctx context.Context) ([]Bonus, error) {
	var bonuses []Bonus
	err := DB.NewSelect().Model(&bonuses).Scan(ctx)
	if err == nil {
		log.Warn().Msgf("Retrieved all bonuses: #%d", len(bonuses))
	}
	return bonuses, err
}

func GetBonusByName(ctx context.Context, name string) (*Bonus, error) {
	var bonus Bonus
	log.Info().Msgf("GetBonusByName name: %s", name)
	err := DB.NewSelect().Model(&bonus).Where("ename = ?", name).Scan(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get bonus with name: %s", name)
		return nil, err
	}
	return &bonus, nil
}

func CreateBonus(ctx context.Context, bonus *Bonus) error {
	_, err := DB.NewInsert().Model(bonus).Exec(ctx)
	if err == nil {
		log.Info().Msgf("Created bonus: %v", bonus)
	} else {
		log.Error().Err(err).Msgf("Failed to create bonus: %v", bonus)
	}
	return err
}

func UpdateBonus(ctx context.Context, bonus *Bonus) error {
	_, err := DB.NewUpdate().Model(bonus).Where("ename = ?", bonus.Ename).Exec(ctx)
	if err == nil {
		log.Info().Msgf("Updated bonus: %v", bonus)
	} else {
		log.Error().Err(err).Msgf("Failed to update bonus: %v", bonus)
	}
	return err
}

func DeleteBonus(ctx context.Context, name string) error {
	_, err := DB.NewDelete().Model(&Bonus{}).Where("ename = ?", name).Exec(ctx)
	if err == nil {
		log.Info().Msgf("Deleted bonus with name: %s", name)
	} else {
		log.Error().Err(err).Msgf("Failed to delete bonus with name: %s", name)
	}
	return err
}

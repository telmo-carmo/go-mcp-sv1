package dal

import (
	"context"
	"encoding/json"
	"net/http"

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

func BonusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodGet:
		name := r.PathValue("name")
		if name == "" {
			bonuses, err := GetAllBonus(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(bonuses)
			return
		}

		bonus, err := GetBonusByName(ctx, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(bonus)
	case http.MethodPost:
		var bonus Bonus
		if err := json.NewDecoder(r.Body).Decode(&bonus); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := CreateBonus(ctx, &bonus); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		var bonus Bonus
		if err := json.NewDecoder(r.Body).Decode(&bonus); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := UpdateBonus(ctx, &bonus); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		//name := r.URL.Query().Get("name")
		name := r.PathValue("name")
		if err := DeleteBonus(ctx, name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

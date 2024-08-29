package activity_locations

import (
	"context"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ActivityLocation struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	ZipCode   string  `json:"zip_code"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Service struct {
	sync.Mutex
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
	}
}

func (service *Service) Create(location ActivityLocation) (ActivityLocation, error) {
	service.Lock()
	defer service.Unlock()

	var id int
	err := service.db.QueryRow(
		context.Background(),
		"INSERT INTO activity_locations (name, address, city, state, zip_code, country, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		location.Name, location.Address, location.City, location.State, location.ZipCode, location.Country, location.Latitude, location.Longitude,
	).Scan(&id)
	if err != nil {
		return ActivityLocation{}, err
	}

	location.ID = id
	return location, nil
}

func (service *Service) ReadAll() ([]ActivityLocation, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, name, address, city, state, zip_code, country, latitude, longitude FROM activity_locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []ActivityLocation
	for rows.Next() {
		var location ActivityLocation
		if err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.City, &location.State, &location.ZipCode, &location.Country, &location.Latitude, &location.Longitude); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (service *Service) Read(id string) (ActivityLocation, bool, error) {
	service.Lock()
	defer service.Unlock()

	var location ActivityLocation
	err := service.db.QueryRow(context.Background(), "SELECT id, name, address, city, state, zip_code, country, latitude, longitude FROM activity_locations WHERE id = $1", id).Scan(
		&location.ID, &location.Name, &location.Address, &location.City, &location.State, &location.ZipCode, &location.Country, &location.Latitude, &location.Longitude)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ActivityLocation{}, false, nil
		}
		return ActivityLocation{}, false, err
	}

	return location, true, nil
}

func (service *Service) Update(id string, location ActivityLocation) (ActivityLocation, bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "UPDATE activity_locations SET name = $1, address = $2, city = $3, state = $4, zip_code = $5, country = $6, latitude = $7, longitude = $8 WHERE id = $9",
		location.Name, location.Address, location.City, location.State, location.ZipCode, location.Country, location.Latitude, location.Longitude, id)
	if err != nil {
		return ActivityLocation{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return ActivityLocation{}, false, nil
	}

	location.ID, _ = strconv.Atoi(id)
	return location, true, nil
}

func (service *Service) Delete(id string) (bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "DELETE FROM activity_locations WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

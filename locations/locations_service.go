package locations

import (
	"context"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Location struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	City      string   `json:"city"`
	State     string   `json:"state"`
	ZipCode   string   `json:"zip_code"`
	Country   string   `json:"country"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
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

func (service *Service) Create(location Location) (Location, error) {
	service.Lock()
	defer service.Unlock()

	var id int
	err := service.db.QueryRow(
		context.Background(),
		"INSERT INTO locations (name, address, city, state, zip_code, country, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		location.Name, location.Address, location.City, location.State, location.ZipCode, location.Country, location.Latitude, location.Longitude,
	).Scan(&id)
	if err != nil {
		return Location{}, err
	}

	location.ID = id
	return location, nil
}

func (service *Service) ReadAll() ([]Location, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, name, address, city, state, zip_code, country, latitude, longitude FROM locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		if err := rows.Scan(&location.ID, &location.Name, &location.Address, &location.City, &location.State, &location.ZipCode, &location.Country, &location.Latitude, &location.Longitude); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (service *Service) Read(id string) (Location, bool, error) {
	service.Lock()
	defer service.Unlock()

	var location Location
	err := service.db.QueryRow(context.Background(), "SELECT id, name, address, city, state, zip_code, country, latitude, longitude FROM locations WHERE id = $1", id).Scan(
		&location.ID, &location.Name, &location.Address, &location.City, &location.State, &location.ZipCode, &location.Country, &location.Latitude, &location.Longitude)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Location{}, false, nil
		}
		return Location{}, false, err
	}

	return location, true, nil
}

func (service *Service) Update(id string, location Location) (Location, bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "UPDATE locations SET name = $1, address = $2, city = $3, state = $4, zip_code = $5, country = $6, latitude = $7, longitude = $8 WHERE id = $9",
		location.Name, location.Address, location.City, location.State, location.ZipCode, location.Country, location.Latitude, location.Longitude, id)
	if err != nil {
		return Location{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return Location{}, false, nil
	}

	location.ID, _ = strconv.Atoi(id)
	return location, true, nil
}

func (service *Service) Delete(id string) (bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "DELETE FROM locations WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

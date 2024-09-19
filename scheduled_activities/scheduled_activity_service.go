package scheduled_activities

import (
	"context"
	"fmt"
	"friendsocial/activity_participants"
	"friendsocial/user_activity_preferences"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ScheduledActivity struct {
	ID                       int       `json:"id"`
	ActivityID               int       `json:"activity_id"`
	IsActive                 bool      `json:"is_active"`
	ScheduledAt              time.Time `json:"scheduled_at"` // New field for scheduled_at
	UserActivityPreferenceID *int      `json:"user_activity_preference_id"`
}

type Service struct {
	sync.Mutex
	db       *pgxpool.Pool
	services *map[string]interface{}
}

func NewService(db *pgxpool.Pool, services *map[string]interface{}) *Service {
	return &Service{
		db:       db,
		services: services,
	}
}

// Create a new scheduled activity
func (service *Service) Create(scheduledActivity ScheduledActivity) (ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	var id int
	err := service.db.QueryRow(
		context.Background(),
		"INSERT INTO scheduled_activities (activity_id, is_active, scheduled_at, user_activity_preference_id) VALUES ($1, $2, $3, $4) RETURNING id",
		scheduledActivity.ActivityID, scheduledActivity.IsActive, scheduledActivity.ScheduledAt, scheduledActivity.UserActivityPreferenceID,
	).Scan(&id)
	if err != nil {
		// Log the error and the values being inserted
		fmt.Printf("Error inserting scheduled activity: %v\n", err)
		fmt.Printf("Values: ActivityID: %d, IsActive: %t, ScheduledAt: %v, UserActivityPreferenceID: %v\n",
			scheduledActivity.ActivityID, scheduledActivity.IsActive, scheduledActivity.ScheduledAt, scheduledActivity.UserActivityPreferenceID)
		return ScheduledActivity{}, fmt.Errorf("failed to insert scheduled activity: %w", err)
	}

	scheduledActivity.ID = id
	return scheduledActivity, nil
}

func (service *Service) CreateMultiple(
	activityID int,
	selectedDates []string,
	scheduledActivitiesStartTime string,
	scheduledActivitiesEndTime string,
	timeZone string,
) ([]ScheduledActivity, error) {

	var scheduledActivities []ScheduledActivity

	// Parse the start and end times
	startTimeParsed, err := time.Parse(time.RFC3339, scheduledActivitiesStartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	endTimeParsed, err := time.Parse(time.RFC3339, scheduledActivitiesEndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %v", err)
	}

	// Load the time zone
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, fmt.Errorf("invalid time zone: %v", err)
	}

	for _, dateStr := range selectedDates {
		// Parse the date
		date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}

		// Skip past dates
		if date.Before(time.Now().In(loc)) {
			continue
		}

		// Combine date and start time
		scheduledAt := time.Date(
			date.Year(), date.Month(), date.Day(),
			startTimeParsed.Hour(), startTimeParsed.Minute(), startTimeParsed.Second(), 0,
			loc,
		)

		// Check availability
		available, err := service.checkIfUserIsAvailable(scheduledAt, startTimeParsed, endTimeParsed, loc)
		if err != nil {
			return nil, err
		}
		if !available {
			continue
		}

		scheduledActivity := ScheduledActivity{
			ActivityID:               activityID,
			IsActive:                 true,
			ScheduledAt:              scheduledAt,
			UserActivityPreferenceID: nil,
		}

		newScheduledActivity, err := service.Create(scheduledActivity)
		if err != nil {
			return nil, fmt.Errorf("failed to create scheduled activity for date %s: %w", dateStr, err)
		}

		scheduledActivities = append(scheduledActivities, newScheduledActivity)
	}

	return scheduledActivities, nil
}

func (service *Service) checkIfUserIsAvailable(
	date time.Time,
	desiredStartTime time.Time,
	desiredEndTime time.Time,
	loc *time.Location,
) (bool, error) {

	// Adjust desiredStartTime and desiredEndTime to the specific date
	desiredStart := time.Date(
		date.Year(), date.Month(), date.Day(),
		desiredStartTime.Hour(), desiredStartTime.Minute(), desiredStartTime.Second(), 0,
		loc,
	)
	desiredEnd := time.Date(
		date.Year(), date.Month(), date.Day(),
		desiredEndTime.Hour(), desiredEndTime.Minute(), desiredEndTime.Second(), 0,
		loc,
	)

	// Retrieve scheduled activities on the date
	rows, err := service.db.Query(
		context.Background(),
		"SELECT id, activity_id, is_active, scheduled_at, user_activity_preference_id FROM scheduled_activities WHERE DATE(scheduled_at) = $1",
		date.Format("2006-01-02"),
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(
			&scheduledActivity.ID,
			&scheduledActivity.ActivityID,
			&scheduledActivity.IsActive,
			&scheduledActivity.ScheduledAt,
			&scheduledActivity.UserActivityPreferenceID,
		); err != nil {
			return false, err
		}

		// Get estimated time for the activity
		estimatedDuration, err := service.getEstimatedTime(scheduledActivity.ActivityID)
		if err != nil {
			return false, err
		}

		activityStart := scheduledActivity.ScheduledAt
		activityEnd := activityStart.Add(estimatedDuration)

		// Check for time overlap
		if activityStart.Before(desiredEnd) && activityEnd.After(desiredStart) {
			return false, nil // Not available
		}
	}

	return true, nil // Available
}

func (service *Service) getEstimatedTime(activityID int) (time.Duration, error) {
	var estimatedTimeInSeconds float64
	err := service.db.QueryRow(
		context.Background(),
		"SELECT EXTRACT(EPOCH FROM estimated_time) FROM activities WHERE id = $1",
		activityID,
	).Scan(&estimatedTimeInSeconds)
	if err != nil {
		return 0, err
	}

	// Convert seconds to time.Duration
	return time.Duration(estimatedTimeInSeconds) * time.Second, nil
}

// Read all user activities for a specific user
func (service *Service) ReadAll() ([]ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, activity_id, is_active, scheduled_at, user_activity_preference_id FROM scheduled_activities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduledActivities []ScheduledActivity
	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt, &scheduledActivity.UserActivityPreferenceID); err != nil {
			return nil, err
		}
		scheduledActivities = append(scheduledActivities, scheduledActivity)
	}

	return scheduledActivities, nil
}

// Read a specific user activity by ID
func (service *Service) Read(id string) (ScheduledActivity, bool, error) {
	service.Lock()
	defer service.Unlock()

	var scheduledActivity ScheduledActivity
	err := service.db.QueryRow(context.Background(), "SELECT id, activity_id, is_active, scheduled_at, user_activity_preference_id FROM scheduled_activities WHERE id = $1", id).Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt, &scheduledActivity.UserActivityPreferenceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ScheduledActivity{}, false, nil
		}
		return ScheduledActivity{}, false, err
	}

	return scheduledActivity, true, nil
}

// Update an existing user activity
func (service *Service) Update(id string, scheduledActivity ScheduledActivity) (ScheduledActivity, bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "UPDATE scheduled_activities SET activity_id = $1, is_active = $2, scheduled_at = $3, user_activity_preference_id = $4 WHERE id = $5",
		scheduledActivity.ActivityID, scheduledActivity.IsActive, scheduledActivity.ScheduledAt, scheduledActivity.UserActivityPreferenceID, id)
	if err != nil {
		return ScheduledActivity{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return ScheduledActivity{}, false, nil
	}

	return scheduledActivity, true, nil
}

// Delete a user activity by ID
func (service *Service) Delete(id string) (bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "DELETE FROM scheduled_activities WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

// Get all active user activities for a specific user
func (service *Service) GetActiveScheduledActivities(userID string) ([]ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, activity_id, is_active, scheduled_at, user_activity_preference_id FROM scheduled_activities WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeScheduledActivities []ScheduledActivity
	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt, &scheduledActivity.UserActivityPreferenceID); err != nil {
			return nil, err
		}
		activeScheduledActivities = append(activeScheduledActivities, scheduledActivity)
	}

	return activeScheduledActivities, nil
}

// Get all inactive user activities for a specific user
func (service *Service) GetInactiveScheduledActivities(userID string) ([]ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, activity_id, is_active, scheduled_at, user_activity_preference_id FROM scheduled_activities WHERE is_active = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inactiveScheduledActivities []ScheduledActivity
	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt, &scheduledActivity.UserActivityPreferenceID); err != nil {
			return nil, err
		}
		inactiveScheduledActivities = append(inactiveScheduledActivities, scheduledActivity)
	}

	return inactiveScheduledActivities, nil
}

func (s *Service) DeclineRepeatedActivity(userID int, scheduledActivityID int) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Get the user_activity_preference_id for the scheduled activity
	var userActivityPreferenceID int
	err = tx.QueryRow(context.Background(),
		"SELECT user_activity_preference_id FROM scheduled_activities WHERE id = $1",
		scheduledActivityID).Scan(&userActivityPreferenceID)
	if err != nil {
		return fmt.Errorf("failed to get user_activity_preference_id: %v", err)
	}

	// Delete all activity participants for this user and all scheduled activities linked to the same user_activity_preference
	_, err = tx.Exec(context.Background(),
		`DELETE FROM activity_participants
		 WHERE user_id = $1 AND scheduled_activity_id IN (
			 SELECT id FROM scheduled_activities
			 WHERE user_activity_preference_id = $2 AND scheduled_at >= NOW()
		 )`,
		userID, userActivityPreferenceID)
	if err != nil {
		return fmt.Errorf("failed to delete activity participants: %v", err)
	}

	return tx.Commit(context.Background())
}

// CreateRepeatingScheduledActivity creates new, repeating, scheduled activities up to 6 months in advance, based on the user's activity preference
func (s *Service) CreateRepeatingScheduledActivity(preference user_activity_preferences.UserActivityPreference, startTime string, timeZone string) ([]ScheduledActivity, error) {
	s.Lock()
	defer s.Unlock()

	// Start a transaction
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	now := time.Now()
	sixMonthsLater := now.AddDate(0, 6, 0)

	// Parse days of week
	daysOfWeek := []time.Weekday{}
	for _, day := range strings.Split(preference.DaysOfWeek, ",") {
		dayInt, err := strconv.Atoi(strings.TrimSpace(day))
		if err != nil {
			return nil, fmt.Errorf("invalid day of week: %v", err)
		}
		daysOfWeek = append(daysOfWeek, time.Weekday(dayInt))
	}

	scheduledActivities := []ScheduledActivity{}

	for currentDate := now; currentDate.Before(sixMonthsLater); currentDate = currentDate.AddDate(0, 0, 1) {
		if !containsWeekday(daysOfWeek, currentDate.Weekday()) {
			continue
		}

		if !shouldScheduleActivity(preference, currentDate, now) {
			continue
		}

		loc, err := time.LoadLocation(timeZone)
		if err != nil {
			return nil, fmt.Errorf("invalid time zone: %v", err)
		}
		startTimeParsed, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time format: %v", err)
		}
		scheduledAt := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), startTimeParsed.Hour(), startTimeParsed.Minute(), startTimeParsed.Second(), 0, loc)

		newActivity := ScheduledActivity{
			ActivityID:               preference.ActivityID,
			IsActive:                 true,
			ScheduledAt:              scheduledAt,
			UserActivityPreferenceID: &preference.ID,
		}

		scheduledActivityService := (*s.services)["scheduled_activities"].(*Service)
		createdScheduledActivity, err := scheduledActivityService.Create(newActivity)
		if err != nil {
			return nil, fmt.Errorf("failed to create scheduled activity: %v", err)
		}

		scheduledActivities = append(scheduledActivities, createdScheduledActivity)
	}

	// Fetch participants for the user activity preference
	rows, err := tx.Query(context.Background(),
		"SELECT user_id FROM user_activity_preferences_participants WHERE user_activity_preference_id = $1",
		preference.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch participants: %v", err)
	}
	defer rows.Close()

	var participantUserIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan participant user ID: %v", err)
		}
		participantUserIDs = append(participantUserIDs, userID)
	}

	// Get the activity participant service
	activityParticipantService := (*s.services)["activity_participants"].(*activity_participants.Service)

	// Create activity participants for each scheduled activity
	for _, scheduledActivity := range scheduledActivities {
		for _, userID := range participantUserIDs {
			participant := activity_participants.ActivityParticipant{
				UserID:              userID,
				ScheduledActivityID: scheduledActivity.ID,
				InviteStatus:        "Pending", // Or whatever default status you want
			}
			_, err := activityParticipantService.Create(participant)
			if err != nil {
				return nil, fmt.Errorf("failed to create activity participant: %v", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return scheduledActivities, nil
}

func containsWeekday(days []time.Weekday, day time.Weekday) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

func shouldScheduleActivity(preference user_activity_preferences.UserActivityPreference, currentDate, startDate time.Time) bool {
	switch preference.FrequencyPeriod {
	case "week":
		weeksSinceStart := int(currentDate.Sub(startDate).Hours() / 168) // 168 hours in a week
		return weeksSinceStart%preference.Frequency == 0
	case "month":
		monthsSinceStart := int(currentDate.Sub(startDate).Hours() / 730) // Approximate hours in a month
		return monthsSinceStart%preference.Frequency == 0
	default:
		return false
	}
}

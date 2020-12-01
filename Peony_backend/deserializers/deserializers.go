package deserializers

import (
	"Peony/Peony_backend/models/entity"
	"encoding/json"
	"errors"
	"time"
)

type RawInfo struct {
	CourseNumber string
	School       string
	FieldTitle   string
	FieldContent string
	Origin       string
	StartTime    string
	EndTime      string
}

func InfoDeserializer(b []byte) (entity.Info, error) {
	var raw_info RawInfo
	err := json.Unmarshal(b, &raw_info)
	if err != nil {
		return entity.Info{}, errors.New("REQUEST BODY MALFORMED.")
	}
	var new_info entity.Info
	start_time, err := time.Parse(time.RFC3339, raw_info.StartTime)
	if err != nil {
		return entity.Info{}, errors.New("START_TIME PARSE ERROR.")
	}
	end_time, err := time.Parse(time.RFC3339, raw_info.EndTime)
	if err != nil {
		return entity.Info{}, errors.New("END_TIME PARSE ERROR.")
	}

	new_info.CourseNumber = raw_info.CourseNumber
	new_info.School = raw_info.School
	new_info.FieldTitle = raw_info.FieldTitle
	new_info.FieldContent = raw_info.FieldContent
	new_info.Origin = raw_info.Origin
	new_info.StartTime = start_time
	new_info.EndTime = end_time

	return new_info, nil
}

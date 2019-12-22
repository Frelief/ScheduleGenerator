package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	year   = "Y"
	fall   = "F"
	spring = "S"
)

type lecture struct {
	MeetingDay        string `json:"meetingDay"`
	MeetingStartTime  string `json:"meetingStartTime"`
	MeetingEndTime    string `json:"meetingEndTime"`
	MeetingScheduleID string `json:"meetingScheduleId"`
	AssignedRoom1     string `json:"assignedRoom1"`
	AssignedRoom2     string `json:"assignedRoom2"`
}

type instructors struct {
	InstructorID string `json:"instructorId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

type meetings struct {
	Schedule            map[string]lecture     `json:"schedule"`
	Instructors         map[string]instructors `json:"instructors"`
	MeetingID           string                 `json:"meetingId"`
	TeachingMethod      string                 `json:"teachingMethod"`
	SectionNumber       string                 `json:"sectionNumber"`
	Subtitle            string                 `json:"subtitle"`
	Cancel              string                 `json:"cancel"`
	Waitlist            string                 `json:"waitlist"`
	Online              string                 `json:"online"`
	EnrollmentCapacity  string                 `json:"enrollmentCapacity"`
	ActualEnrolment     string                 `json:"actualEnrolment"`
	ActualWaitlist      string                 `json:"actualWaitlist"`
	EnrollmentIndicator string                 `json:"enrollmentIndicator"`
	MeetingStatusNotes  string                 `json:"meetingStatusNotes,omitempty"`
	EnrollmentControls  string                 `json:"enrollmentControls"` // This is actually a list of things, we need to expand if we want to use anything from here
}

type course struct {
	CourseID                 string              `json:"courseId"`
	Org                      string              `json:"org"`
	OrgName                  string              `json:"orgName"`
	CourseTitle              string              `json:"courseTitle"`
	Code                     string              `json:"code"`
	CourseDescription        string              `json:"courseDescription"`
	Prerequisite             string              `json:"prerequisite"`
	Corequisite              string              `json:"corequisite"`
	Exclusion                string              `json:"exclusion"`
	RecommendedPreparation   string              `json:"recommendedPreparation"`
	Section                  string              `json:"section"`
	Session                  string              `json:"session"`
	WebTimetableInstructions string              `json:"webTimetableInstructions"`
	BreadthCategories        string              `json:"breadthCategories"`
	DistributionCategories   string              `json:"distributionCategories"`
	Meetings                 map[string]meetings `json:"meetings"`
}

func getCourseInfo(courseName string, session string) (map[string]meetings, error) {
	var file string = "./Courses/courses" + session + ".json"
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	calendarMap := make(map[string]course)
	json.Unmarshal(byteValue, &calendarMap)
	for course, value := range calendarMap {
		if strings.Contains(course, courseName) {
			return value.Meetings, nil
		}
	}
	return nil, errors.New("Course not found")
}

func main() {
	fmt.Println(getCourseInfo("CSC207", fall))
}

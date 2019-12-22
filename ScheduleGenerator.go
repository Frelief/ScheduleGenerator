package main

import (
	"fmt"
)

const (
	Year   = "Y"
	Fall   = "F"
	Spring = "S"
)

type Lecture struct {
	meetingDay string `json.meetingDay`
	meetingStartTime string `json.meetingStartTime`
	meetingEndTime string `json.meetingEndTime`
	meetingScheduleId string `json.meetingScheduleId`
	assignedRoom1 string `json.assignedRoom1`
	assignedRoom2 string `json.assignedRoom2`
}

type Instructors struct {
	instructorID string `json.instructorId`
	firstName string `json.firstName`
	lastName string `json.lastName`
}

type Meetings struct {
	schedule []Lecture `json.schedule`
	instructors []Instructors `json.instructors`
	meetingID string `json.meetingId`
	teachingMethod string `json.teachingMethod`
	sectionNumber string `json.sectionNumber`
	
	"teachingMethod": "LEC",
	"sectionNumber": "0101",
	"subtitle": "",
	"cancel": "",
	"waitlist": "N",
	"online": "",
	"enrollmentCapacity": "300",
	"actualEnrolment": "170",
	"actualWaitlist": "0",
	"enrollmentIndicator": "E",
	"meetingStatusNotes": null,
	"enrollmentControls": []
}

type Course struct {
	courseID string `json:courseId`
	org string `json:org`
	orgName string `json:orgName`
	courseTitle string `json:courseTitle`
	code string `json:code`
	courseDescription string `json:courseDescription`
	prerequisite string `json.prerequisite`
	corequisite string `json.corequisite`
	exclusion string `json.exclusion`
	recommendedPreparation string `json.recommendedPreparation`
	section string `json.section`
	session string `json.session`
	webTimetableInstructions string `json.webTimetableInstructions`
	breadthCategories string `json.breadthCategories`
	distributionCategories string `json.distributionCategories`
	meetings
        "section": "Y",
        "session": "20199",
        "webTimetableInstructions": "This course is only open to students admitted to the Vic One program (http:\/\/www.vic.utoronto.ca\/Future_Students\/vicone.htm).",
        "breadthCategories": "",
        "distributionCategories": "",
        "meetings": {
            "LEC-0101": {
                "schedule": {
                    "WE-163117": {
                        "meetingDay": "WE",
                        "meetingStartTime": "16:00",
                        "meetingEndTime": "18:00",
                        "meetingScheduleId": "163117",
                        "assignedRoom1": "BT 101",
                        "assignedRoom2": "BT 101"
                    }
                },
                "instructors": [],
                "meetingId": "103406",
                "teachingMethod": "LEC",
                "sectionNumber": "0101",
                "subtitle": "",
                "cancel": "",
                "waitlist": "N",
                "online": "",
                "enrollmentCapacity": "300",
                "actualEnrolment": "170",
                "actualWaitlist": "0",
                "enrollmentIndicator": "E",
                "meetingStatusNotes": null,
                "enrollmentControls": []
            }
        }
}

func getCourseInfo(course string, session string) string {
	var file string
	switch session {
	case Year:
		file = "/Courses/coursesY.json"
	case Fall:
		file = "/Courses/coursesF.json"
	case Spring:
		file = "/Courses/coursesS.json"
	}
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	return file
}

func main() {
	fmt.Println(getCourseInfo("CSC207", Fall))
}

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
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

type meeting struct {
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
	CourseID                 string             `json:"courseId"`
	Org                      string             `json:"org"`
	OrgName                  string             `json:"orgName"`
	CourseTitle              string             `json:"courseTitle"`
	Code                     string             `json:"code"`
	CourseDescription        string             `json:"courseDescription"`
	Prerequisite             string             `json:"prerequisite"`
	Corequisite              string             `json:"corequisite"`
	Exclusion                string             `json:"exclusion"`
	RecommendedPreparation   string             `json:"recommendedPreparation"`
	Section                  string             `json:"section"`
	Session                  string             `json:"session"`
	WebTimetableInstructions string             `json:"webTimetableInstructions"`
	BreadthCategories        string             `json:"breadthCategories"`
	DistributionCategories   string             `json:"distributionCategories"`
	Meetings                 map[string]meeting `json:"meetings"`
}

type schedule struct {
	Courses  []string
	Semester string
	Classes  map[string]meeting
}

func (sched schedule) checkConflict(newMeeting meeting) bool {
	for _, meeting := range sched.Classes {
		for _, newLecture := range newMeeting.Schedule {
			for _, lecture := range meeting.Schedule {
				lectureStart, _ := time.Parse("15:04", lecture.MeetingStartTime)
				lectureEnd, _ := time.Parse("15:04", lecture.MeetingEndTime)
				newLectureStart, _ := time.Parse("15:04", newLecture.MeetingStartTime)
				newLectureEnd, _ := time.Parse("15:04", newLecture.MeetingEndTime)
				if lecture.MeetingDay == newLecture.MeetingDay && newLectureStart.Sub(lectureEnd) < 0 && newLectureEnd.Sub(lectureStart) > 0 {
					return false
				}
			}
		}
	}
	return true
}

func getCourseInfo(courseName string, session string) (map[string]meeting, error) {
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

func getAllCourses() []map[string]meeting {
	fmt.Println("How many courses do you wish to take?")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	numCourse, _ := strconv.Atoi(strings.Split(text, "\n")[0])
	coursesInfoArray := make([]map[string]meeting, numCourse)
	for index := 0; index < numCourse; index++ {
		fmt.Println("Course number", index)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		temp := strings.Split(text, "\n")
		var err error
		coursesInfoArray[index], err = getCourseInfo(temp[0], fall)
		if err != nil {
			fmt.Println(err)
			index--
		}
	}
	return coursesInfoArray
}

func buildAllSchedules(coursesInfoArray []map[string]meeting) {
	numCourse := len(coursesInfoArray)
	for index := 0; index < numCourse; index++ {

	}
}

func main() {
	testCourse, _ := getCourseInfo("CSC207", "F")
	var sched = map[string]meeting{"LEC-0101": testCourse["LEC-0101"]}
	fmt.Println(sched)
	testSchedule := schedule{Classes: sched}

	fmt.Println(testSchedule.checkConflict(testCourse["LEC-0101"]))

}

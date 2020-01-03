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
	lec    = "LEC"
	tut    = "TUT"
	pra    = "PRA"
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
	Courses []string
	Classes map[string]meeting
}

func createSchedule() schedule {
	sched := schedule{}
	sched.Classes = make(map[string]meeting)
	return sched
}

func (sched schedule) copySchedule() schedule {
	copiedSchedule := createSchedule()
	copy(copiedSchedule.Courses, sched.Courses)
	for k, v := range sched.Classes {
		copiedSchedule.Classes[k] = v
	}
	copiedSchedule.Courses = sched.Courses
	return copiedSchedule
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

func (sched *schedule) addClass(courseName string, newMeeting meeting) {
	sched.Courses = append(sched.Courses, courseName)
	sched.Classes[courseName] = newMeeting
}

/*
Meetings can be either Lectures, Tutorials or Practicals. In each case we need to enroll in one of each
*/
func separateMeetingTypes(meetings map[string]meeting) (map[string]meeting, map[string]meeting, map[string]meeting, error) {
	lectureMap := make(map[string]meeting)
	tutorialMap := make(map[string]meeting)
	practicalMap := make(map[string]meeting)

	for meetingName, meeting := range meetings {
		if meeting.TeachingMethod == lec {
			lectureMap[meetingName] = meeting
		} else if meeting.TeachingMethod == tut {
			tutorialMap[meetingName] = meeting
		} else if meeting.TeachingMethod == pra {
			practicalMap[meetingName] = meeting
		} else {
			return nil, nil, nil, errors.New("Unexpected teaching method: " + meeting.TeachingMethod)
		}
	}
	return lectureMap, tutorialMap, practicalMap, nil
}

func getCourseInfo(courseName string, session string) ([]map[string]map[string]meeting, error) {
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
	var courseArray []map[string]map[string]meeting
	for course, value := range calendarMap {
		if strings.Contains(course, courseName) {
			lectureMap, tutorialMap, practicalMap, err := separateMeetingTypes(value.Meetings)
			if err != nil {
				return nil, err
			}
			if len(lectureMap) != 0 {
				courseMap := make(map[string]map[string]meeting)
				courseMap[course+lec] = lectureMap
				courseArray = append(courseArray, courseMap)
			}
			if len(tutorialMap) != 0 {
				courseMap := make(map[string]map[string]meeting)
				courseMap[course+tut] = tutorialMap
				courseArray = append(courseArray, courseMap)
			}
			if len(practicalMap) != 0 {
				courseMap := make(map[string]map[string]meeting)
				courseMap[course+pra] = practicalMap
				courseArray = append(courseArray, courseMap)
			}
			return courseArray, nil
		}
	}
	return nil, errors.New("Course not found")
}

func getAllCourses() []map[string]map[string]meeting {
	fmt.Println("How many courses do you wish to take?")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	numCourse, _ := strconv.Atoi(strings.Split(text, "\n")[0])
	// coursesInfoArray := make([]map[string]map[string]meeting, numCourse)
	var coursesInfoArray []map[string]map[string]meeting
	for index := 0; index < numCourse; index++ {
		fmt.Println("Course number", index)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		temp := strings.Split(text, "\n")
		meetingMap, err := getCourseInfo(temp[0], fall)
		if err != nil {
			fmt.Println(err)
			index--
		} else {
			coursesInfoArray = append(coursesInfoArray, meetingMap...)
		}
	}
	return coursesInfoArray
}

func buildAllSchedules(sched schedule, coursesInfoArray []map[string]map[string]meeting) []schedule {
	var schedArray []schedule
	for courseName, meetings := range coursesInfoArray[0] {
		for _, class := range meetings {
			schedCopy := sched.copySchedule()
			if schedCopy.checkConflict(class) {
				schedCopy.addClass(courseName, class)
				if len(coursesInfoArray) > 1 {
					schedArray = append(schedArray, buildAllSchedules(schedCopy, coursesInfoArray[1:])...)
				} else {
					schedArray = append(schedArray, schedCopy)
				}
			}
		}
	}
	return schedArray
}

func main() {
	testSchedule := createSchedule()
	// buildAllSchedules(testSchedule, getAllCourses())
	for index, element := range buildAllSchedules(testSchedule, getAllCourses()) {
		// fmt.Println(buildAllSchedules(testSchedule, getAllCourses()))
		fmt.Println("Schedule", index)
		fmt.Println(element.Courses)
		fmt.Println(element.Classes)
	}
}

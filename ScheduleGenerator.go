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

// ListOfSchedules is a list of schedules
var ListOfSchedules []schedule

//TODO: make this pointer instead of object
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

//TODO: implement schedule check via 263 tutorial
/*
* schedule.checkNoConflict checks whether the meeting newMeeting would conflict with the meetings
* already inside the schedule. Returns false if there is a conflict.
 */
func (sched schedule) checkNoConflict(newMeeting meeting) bool {
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

/*
* schedule.addClass adds a meeting newMeeting to the schedule, with name courseName.
*
 */
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
		switch meeting.TeachingMethod {
		case lec:
			lectureMap[meetingName] = meeting
		case tut:
			tutorialMap[meetingName] = meeting
		case pra:
			practicalMap[meetingName] = meeting
		default:
			return nil, nil, nil, errors.New("Unexpected teaching method: " + meeting.TeachingMethod)
		}
	}
	return lectureMap, tutorialMap, practicalMap, nil
}

/*
* getCourseInfo returns the information for the course courseName during the session session.
* returns an error if it cannot find the course.
 */
func getCourseInfo(courseName string, session string) ([]map[string]map[string]meeting, error) {
	if session != fall && session != spring && session != year {
		return nil, errors.New("Unknown semester: " + session)
	}
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

			//Maps Meeting Types to their respective Meeting Maps
			meetingTypesToMeetingMaps := map[string](map[string]meeting){lec: lectureMap, tut: tutorialMap, pra: practicalMap}

			for meetingType, meetingMap := range meetingTypesToMeetingMaps {
				if len(meetingMap) != 0 {
					courseMap := make(map[string]map[string]meeting)
					courseMap[course+meetingType] = meetingMap
					courseArray = append(courseArray, courseMap)
				}
			}

			return courseArray, nil
		}
	}
	return nil, errors.New("Course not found")
}

/*
* getAllCourses returns a map of all the courses the user wants to take based on input
*
 */
func getAllCourses() []map[string]map[string]meeting {
	fmt.Println("How many courses do you wish to take?")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	numCourse, _ := strconv.Atoi(strings.Split(text, "\n")[0])
	// coursesInfoArray := make([]map[string]map[string]meeting, numCourse)
	var coursesInfoArray []map[string]map[string]meeting
	for index := 0; index < numCourse; index++ {
		fmt.Println("Course code", index)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		courseCode := strings.Split(text, "\n")

		fmt.Println("Which semester? F S Y")
		text, _ = reader.ReadString('\n')
		semester := strings.Split(text, "\n")
		meetingMap, err := getCourseInfo(courseCode[0], semester[0])
		if err != nil {
			fmt.Println(err)
			index--
		} else {
			coursesInfoArray = append(coursesInfoArray, meetingMap...)
		}
	}
	return coursesInfoArray
}

//TODO: merge into buildAllSchedules?
func testMethod(sched schedule, schedArray []schedule, class meeting, courseName string, coursesInfoArray []map[string]map[string]meeting) []schedule {
	schedCopy := sched.copySchedule()
	if schedCopy.checkNoConflict(class) {
		schedCopy.addClass(courseName, class)
		if len(coursesInfoArray) > 1 {
			schedArray = append(schedArray, buildAllSchedules(schedCopy, coursesInfoArray[1:])...)
		} else {
			schedArray = append(schedArray, schedCopy)
		}
	}
	return schedArray
}

/*
* buildAllSchedules buils all possible schedules without any conflict
*
 */
func buildAllSchedules(sched schedule, coursesInfoArray []map[string]map[string]meeting) []schedule {
	var schedArray []schedule
	for courseName, meetings := range coursesInfoArray[0] {
		for _, class := range meetings {
			temp := testMethod(sched, schedArray, class, courseName, coursesInfoArray)
			schedArray = temp
		}
	}
	return schedArray
}

/*
* addToListOfSchedules appends schedule to ListOfSchedules public array
 */
func addToListOfSchedules(sched schedule) {
	ListOfSchedules = append(ListOfSchedules, sched)
}

/*
* addToSchedule makes all the possible schedules?
 */
func addToSchedule(sched schedule, courseInfoArray []map[string]map[string]meeting) {
	fmt.Println(len(courseInfoArray))
	if len(courseInfoArray) == 0 {
		addToListOfSchedules(sched)
	} else {
		nextClass := courseInfoArray[0]
		newCourseInfoArray := courseInfoArray[1:]
		for courseName, meetings := range nextClass {
			for _, class := range meetings {
				newSchedule := sched.copySchedule()
				if newSchedule.checkNoConflict(class) {
					newSchedule.addClass(courseName, class)
					go addToSchedule(newSchedule, newCourseInfoArray)
				}
			}
		}
	}
}

func main() {
	testSchedule := createSchedule()
	// buildAllSchedules(testSchedule, getAllCourses())
	fmt.Println(ListOfSchedules)
	addToSchedule(testSchedule, getAllCourses())
	fmt.Println(ListOfSchedules)
	/*
		for index, element := range ListOfSchedules {
			// fmt.Println(buildAllSchedules(testSchedule, getAllCourses()))
			fmt.Println("Schedule", index)
			fmt.Println(element.Courses)
			fmt.Println(element.Classes)
		}
	*/
}

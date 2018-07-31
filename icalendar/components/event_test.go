package components

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/gtrevg/caldav-go/icalendar"
	"github.com/gtrevg/caldav-go/icalendar/values"
	. "github.com/gtrevg/check"
)

type EventSuite struct{}

var _ = Suite(new(EventSuite))

func TestEvent(t *testing.T) { TestingT(t) }

func (s *EventSuite) TestMissingEndMarshal(c *C) {
	now := time.Now().UTC()
	event := NewEvent("test", now)
	_, err := icalendar.Marshal(event)
	c.Assert(err, ErrorMatches, "end date or duration must be set")
}

func (s *EventSuite) TestBasicWithDurationMarshal(c *C) {
	now := time.Now().UTC()
	event := NewEventWithDuration("test", now, time.Hour)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDURATION:PT1H\r\nEND:VEVENT"
	fdate := now.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, fdate, fdate))
}

func (s *EventSuite) TestBasicWithEndMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	event := NewEventWithEnd("test", now, end)
	enc, err := icalendar.Marshal(event)
	c.Assert(err, IsNil)
	tmpl := "BEGIN:VEVENT\r\nUID:test\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate))
}

func (s *EventSuite) TestFullEventMarshal(c *C) {
	now := time.Now().UTC()
	end := now.Add(time.Hour)
	oneDay := time.Hour * 24
	oneWeek := oneDay * 7
	event := NewEventWithEnd("1:2:3", now, end)
	uri, _ := url.Parse("http://gtrevg.com/some/attachment.ics")
	event.Attachment = values.NewUrl(*uri)
	event.Attendees = []*values.AttendeeContact{
		values.NewAttendeeContact("Jon Azoff", "jon@gtrevg.com"),
		values.NewAttendeeContact("Matthew Davie", "matthew@gtrevg.com"),
	}
	event.Categories = values.NewCSV("vinyasa", "level 1")
	event.Comments = values.NewComments("Great class, 5 stars!", "I love this class!")
	event.ContactInfo = values.NewCSV("Send us an email!", "<jon@gtrevg.com>")
	event.Created = event.DateStart
	event.Description = "An all-levels class combining strength and flexibility with breath"
	ex1 := values.NewDateTime(now.Add(oneWeek))
	ex2 := values.NewDateTime(now.Add(oneWeek * 2))
	event.ExceptionDateTimes = values.NewExceptionDateTimes(ex1, ex2)
	event.Geo = values.NewGeo(37.747643, -122.445400)
	event.LastModified = event.DateStart
	event.Location = values.NewLocation("Dolores Park")
	event.Organizer = values.NewOrganizerContact("Jon Azoff", "jon@gtrevg.com")
	event.Priority = 1
	event.RecurrenceId = event.DateStart
	r1 := values.NewDateTime(now.Add(oneWeek + oneDay))
	r2 := values.NewDateTime(now.Add(oneWeek*2 + oneDay))
	event.RecurrenceDateTimes = values.NewRecurrenceDateTimes(r1, r2)
	event.AddRecurrenceRules(values.NewRecurrenceRule(values.WeekRecurrenceFrequency))
	uri, _ = url.Parse("matthew@gtrevg.com")
	event.RelatedTo = values.NewUrl(*uri)
	event.Resources = values.NewCSV("yoga mat", "towel")
	event.Sequence = 1
	event.Status = values.TentativeEventStatus
	event.Summary = "Jon's Super-Sweaty Vinyasa 1"
	event.TimeTransparency = values.OpaqueTimeTransparency
	uri, _ = url.Parse("http://student.gtrevg.com/san-francisco/jonathan-azoff/vinyasa-1")
	event.Url = values.NewUrl(*uri)
	enc, err := icalendar.Marshal(event)
	if err != nil {
		c.Fatal(err.Error())
	}
	tmpl := "BEGIN:VEVENT\r\nUID:1:2:3\r\nDTSTAMP:%sZ\r\nDTSTART:%sZ\r\nDTEND:%sZ\r\nCREATED:%sZ\r\n" +
		"DESCRIPTION:An all-levels class combining strength and flexibility with breath\r\n" +
		"GEO:37.747643 -122.445400\r\nLAST-MODIFIED:%sZ\r\nLOCATION:Dolores Park\r\n" +
		"ORGANIZER;CN=\"Jon Azoff\":MAILTO:jon@gtrevg.com\r\nPRIORITY:1\r\nSEQUENCE:1\r\nSTATUS:TENTATIVE\r\n" +
		"SUMMARY:Jon's Super-Sweaty Vinyasa 1\r\nTRANSP:OPAQUE\r\n" +
		"URL;VALUE=URI:http://student.gtrevg.com/san-francisco/jonathan-azoff/vinyasa-1\r\n" +
		"RECURRENCE-ID:%sZ\r\nRRULE:FREQ=WEEKLY\r\nATTACH;VALUE=URI:http://gtrevg.com/some/attachment.ics\r\n" +
		"ATTENDEE;CN=\"Jon Azoff\":MAILTO:jon@gtrevg.com\r\nATTENDEE;CN=\"Matthew Davie\":MAILTO:matthew@gtrevg.com\r\n" +
		"CATEGORIES:vinyasa,level 1\r\nCOMMENT:Great class, 5 stars!\r\nCOMMENT:I love this class!\r\n" +
		"CONTACT:Send us an email!,<jon@gtrevg.com>\r\nEXDATE:%s,%s\r\nRDATE:%s,%s\r\n" +
		"RELATED-TO;VALUE=URI:matthew@gtrevg.com\r\nRESOURCES:yoga mat,towel\r\nEND:VEVENT"
	sdate := now.Format(values.DateTimeFormatString)
	edate := end.Format(values.DateTimeFormatString)
	c.Assert(enc, Equals, fmt.Sprintf(tmpl, sdate, sdate, edate, sdate, sdate, sdate, ex1, ex2, r1, r2))
}

func (s *EventSuite) TestQualifiers(c *C) {
	now := time.Now().UTC()
	event := NewEventWithDuration("test", now, time.Hour)
	c.Assert(event.IsRecurrence(), Equals, false)
	event.RecurrenceId = values.NewDateTime(now)
	c.Assert(event.IsRecurrence(), Equals, true)
	c.Assert(event.IsOverride(), Equals, false)
	event.DateStart = values.NewDateTime(now.Add(time.Hour))
	c.Assert(event.IsRecurrence(), Equals, true)
	c.Assert(event.IsOverride(), Equals, true)
}

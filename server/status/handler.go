package status

import (
	"fmt"
	"github.com/peter-mount/golib/rest"
	"strings"
	"time"
)

// Generates the JSON rest response for external handling of statuses
func (s *Status) jsonHandler(r *rest.Rest) error {
	r.Status(200).
		JSON().
		Value(s.Snapshot())

	return nil
}

// Generates the JSON rest response for external handling of statuses
func (s *Status) xmlHandler(r *rest.Rest) error {
	r.Status(200).
		XML().
		Value(s.Snapshot())

	return nil
}

func (s *Status) htmlHandler(r *rest.Rest) error {

	var a = []string{
		"<html><head><title>Status</title>",
		"<style type=\"text/css\">",
		"body {font-family: arial, helvetica, sans-serif;font-size: 12px;font-weight: normal;color: black;background: white;}",
		"table.tbl {border-collapse:collapse;border-style:none;}",
		".titre{background:#20d0d0;color:#000000;font-weight:bold;text-align:centre;}",
		"th.pxname {background: #b00040;color: #ffff40;font-weight: bold;border-style: solid solid none solid;padding: 2px 3px;white-space: nowrap;}",
		"table.tbl th {border-width: 1px;border-style: solid solid solid solid;border-color: gray;}",
		"th, td {font-size: 10px;}",
		"table.tbl td {text-align: right;border-width: 1px 1px 1px 1px;border-style: solid solid solid solid;padding: 2px 3px;border-color: gray;white-space: nowrap;}",
		"table.tbl td.ac {text-align: left;}",
		"a.px:visited {color: #ffff40;text-decoration: none;}",
		"a.px:link {color: #ffff40;text-decoration: none;}",
		"table.tbl th.empty {border-style: none;empty-cells: hide;background: white;}",
		"table.tbl th {border-width: 1px;border-style: solid solid solid solid;border-color: gray;}",
		".active_up {background: #c0ffc0;}",
		".active_down {background: #ff9090;}",
		".backend {background: #e8e8d0;}",
		"</style>",
		"<meta http-equiv=\"refresh\" content=\"10\"/>",
		"</head><body>",
	}

	for _, cat := range s.Snapshot().Entries {

		a = append(a,
			"<table class=\"tbl\" width=\"100%\">",
			"<tr class=\"titre\">",
			"<th class=\"pxname\" width=\"10%\">", cat.Title, "</th>",
			"<th class=\"empty\" width=\"90%\"></th>",
			"</tr></table>",
			"<table class=\"tbl\" width=\"100%\">",
			"<tr class=\"titre\">",
			"<th rowspan=\"2\"></th>",
			"<th colspan=\"3\">Time</th>",
			"<th colspan=\"2\">Local</th>",
			"<th colspan=\"2\">Remote</th>",
			"<th rowspan=\"2\">Secure</th>",
			"<th colspan=\"2\">Bytes</th>",
			"</tr>",
			"<tr class=\"titre\">",
			"<th>Connected</th><th>Duration</th><th>Idle</th>",
			"<th>Interface</th><th>Port</th>",
			"<th>Interface</th><th>Port</th>",
			"<th>In</th><th>Out</th>",
			"</tr>",
		)

		for _, con := range cat.Entries {
			a = append(a, "<tr class=\"active_up\"><td class=\"ac\">")
			if con.Name == "" {
				a = append(a, fmt.Sprintf("%d", con.ID))
			} else {
				a = append(a, con.Name)
			}
			a = append(a, "</td>")

			a = append(a,
				"<td>",
				con.Stats.Time.Format(time.RFC3339),
				"</td><td>",
				con.Stats.Duration.String(),
				//now.Sub(con.Time).Truncate(time.Second).String(),
				"</td><td>",
				con.Stats.Idle.String(),
				"</td>",
			)

			a = con.Local.Append(a)
			a = con.Remote.Append(a)
			a = append(a, fmt.Sprintf(
				"<td>%v</td><td>%d</td><td>%d</td>",
				con.Secure,
				con.Stats.BytesIn,
				con.Stats.BytesOut,
			))
			a = append(a, "</tr>")
		}

		a = append(a,
			"<tr class=\"backend\"><td class=\"ac\">", cat.Title, "</td>",
			"<td>", cat.Stats.Time.Format(time.RFC3339), "</td>",
			"<td>", time.Now().Sub(cat.Stats.Time).Truncate(time.Second).String(), "</td>",
			"<td>", cat.Stats.Idle.Truncate(time.Second).String(), "</td>",
			"<td colspan=\"5\"></td>",
		)
		a = append(a, fmt.Sprintf(
			"<td>%d</td><td>%d</td>",
			cat.Stats.BytesIn,
			cat.Stats.BytesOut,
		))
		a = append(a, "</tr></table><p></p>")

	}

	a = append(a, "</body></html>")

	r.Status(200).HTML()
	r.Writer().Write([]byte(strings.Join(a, "")))

	return nil
}

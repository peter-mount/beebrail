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
		"table.tbl td.ac {text-align: center;}",
		"a.px:visited {color: #ffff40;text-decoration: none;}",
		"a.px:link {color: #ffff40;text-decoration: none;}",
		"table.tbl th.empty {border-style: none;empty-cells: hide;background: white;}",
		"table.tbl th {border-width: 1px;border-style: solid solid solid solid;border-color: gray;}",
		".active_up {background: #c0ffc0;}",
		".active_down {background: #ff9090;}",
		".backend {background: #e8e8d0;}",
		"</style>",
		"</head><body>",
	}

	for _, cat := range s.Snapshot().Entries {
		// Title row for category
		a = append(a,
			"<table class=\"tbl\">",
			"<tr class=\"titre\"><th class=\"pxname\" width=\"20%\">",
			"<a name=\"", cat.Name, "\"></a>",
			"<a class=\"px\" href=\"#", cat.Name, "\">", cat.Title, "</a>",
		)

		if cat.Port > 0 {
			a = append(a, fmt.Sprintf(" %d", cat.Port))
		}

		a = append(a,
			"</th><th class=\"empty\" width=\"80%\"></th>",
			"</tr>",
			"</table>",
		)

		// Now the table containing the connections
		a = append(a, "<table class=\"tbl\">",
			"<tr class=\"titre\">",
			"<th rowspan=\"2\"></th>",
			"<th colspan=\"2\">Time</th>",
			"<th colspan=\"2\">Local</th>",
			"<th colspan=\"2\">Remote</th>",
			"</tr>",
			"<tr class=\"titre\">",
			"<th>Connected</th><th>Duration</th>",
			"<th>Interface</th><th>Port</th>",
			"<th>Interface</th><th>Port</th>",
			"</tr>",
		)
		for _, con := range cat.Entries {
			a = append(a, "<tr class=\"active_up\"><td class=\"ac\">")
			if con.Name == "" {
				a = append(a, fmt.Sprintf("connection%d", con.ID))
			} else {
				a = append(a, con.Name)
			}
			a = append(a, "</td>")

			a = append(a,
				"<td>",
				con.Time.Format(time.RFC3339),
				"</td><td>",
				con.Duration.String(),
				"</td>",
			)

			a = con.Local.Append(a)
			a = con.Remote.Append(a)
			a = append(a, "</tr>")
		}
		a = append(a, "</table>")

		// Category separator
		a = append(a, "<p></p>")
	}

	a = append(a,
		"</table>",
		"</body></html>",
	)

	r.Status(200).HTML()
	r.Writer().Write([]byte(strings.Join(a, "")))

	return nil
}

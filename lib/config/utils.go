package config

import "time"

func parseTimeFormat(fmt string) (ts_format string) {
	switch fmt {
	case "ANSIC":
		{
			ts_format = time.ANSIC
		}
	case "UnixDate":
		{
			ts_format = time.UnixDate
		}
	case "RubyDate":
		{
			ts_format = time.RubyDate
		}
	case "RFC822":
		{
			ts_format = time.RFC822
		}
	case "RFC822Z":
		{
			ts_format = time.RFC822Z
		}
	case "RFC850":
		{
			ts_format = time.RFC850
		}
	case "RFC1123":
		{
			ts_format = time.RFC1123
		}
	case "RFC1123Z":
		{
			ts_format = time.RFC1123Z
		}
	case "RFC3339":
		{
			ts_format = time.RFC3339
		}
	case "RFC3339Nano":
		{
			ts_format = time.RFC3339Nano
		}
	case "Kitchen":
		{
			ts_format = time.Kitchen
		}
	case "Stamp":
		{
			ts_format = time.Stamp
		}
	case "StampMilli":
		{
			ts_format = time.StampMilli
		}
	case "StampMicro":
		{
			ts_format = time.StampMicro
		}
	case "StampNano":
		{
			ts_format = time.StampNano
		}
	default:
		{
			ts_format = fmt
		}
	}
	return
}

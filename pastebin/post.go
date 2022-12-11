package pastebin

type ExpiryTime int

type Visibility int

const (
	ExpireNever ExpiryTime = iota
	Expire10Mins
	Expire1Hour
	Expire1Day
	Expire1Week
	Expire2Weeks
	Expire1Month
	Expire6Months
	Expire1Year
)

const (
	VisibilityPublic Visibility = iota
	VisibilityUnlisted
	VisibilityPrivate
)

func (e ExpiryTime) String() string {
	switch e {
	case ExpireNever:
		return "N"
	case Expire10Mins:
		return "10M"
	case Expire1Hour:
		return "1H"
	case Expire1Day:
		return "1D"
	case Expire1Week:
		return "1W"
	case Expire2Weeks:
		return "2W"
	case Expire1Month:
		return "1M"
	case Expire6Months:
		return "6M"
	case Expire1Year:
		return "1Y"
	default:
		return ""
	}
}

func ExpirationFromString(time string) ExpiryTime {
	switch time {
	case ExpireNever.String(), "never":
		return ExpireNever
	case Expire10Mins.String(), "10mins":
		return Expire10Mins
	case Expire1Hour.String(), "1hour", "hour":
		return Expire1Hour
	case Expire1Day.String(), "1day", "day":
		return Expire1Day
	case Expire1Week.String(), "1week", "week":
		return Expire1Week
	case Expire2Weeks.String(), "2weeks":
		return Expire2Weeks
	case Expire1Month.String(), "1month", "month":
		return Expire1Month
	case Expire6Months.String(), "6months":
		return Expire6Months
	case Expire1Year.String(), "1year", "year":
		return Expire1Year
	default:
		return ExpiryTime(-1)
	}
}

func (v Visibility) String() string {
	switch v {
	case VisibilityPublic:
		return "0"
	case VisibilityUnlisted:
		return "1"
	case VisibilityPrivate:
		return "2"
	default:
		return ""
	}
}

func VisibilityFromString(v string) Visibility {
	switch v {
	case VisibilityPublic.String(), "public", "pub":
		return VisibilityPublic
	case VisibilityUnlisted.String(), "unlisted":
		return VisibilityUnlisted
	case VisibilityPrivate.String(), "private", "priv":
		return VisibilityPrivate
	default:
		return Visibility(-1)
	}
}

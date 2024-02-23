package version

var (
	version        = "0.0.0"
	buildVersion   = "unknown"
	buildTime      = "unknown"
	goVersion      = "unknown"
	lastCommitTime = "unknown"
	goos           = "unknown"
	goarch         = "unknown"
)

func Version() string {
	return version
}
func BuildVersion() string {
	return buildVersion
}
func BuildTime() string {
	return buildTime
}
func GoVersion() string {
	return goVersion
}
func LastCommitTime() string {
	return lastCommitTime
}
func Goos() string {
	return goos
}
func Goarch() string {
	return goarch
}

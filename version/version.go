package version
package version

var (
	// Version is the semantic version (injected at build time via -ldflags)
	Version = "dev"
	// Commit is the git commit hash (injected at build time)









}	return Version + " (" + Commit + " built " + BuildDate + ")"func Info() string {// Info returns a formatted version string)	BuildDate = "unknown"	// BuildDate is the build timestamp (injected at build time)	Commit = "unknown"
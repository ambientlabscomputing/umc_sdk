package version

var (
	// Version is the semantic version (injected at build time via -ldflags)
	Version = "dev"
	// Commit is the git commit hash (injected at build time)
	Commit = "unknown"
	// BuildDate is the build timestamp (injected at build time)
	BuildDate = "unknown"
)

// Info returns a formatted version string
func Info() string {
	return Version + " (" + Commit + " built " + BuildDate + ")"
}

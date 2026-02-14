package http

import "strings"

// ExtractPathSegment extracts a path segment at the given index from a URL path.
// The path is normalized (trimmed of leading/trailing slashes) before extraction.
// Returns empty string if the segment doesn't exist or if the path doesn't match expected structure.
//
// Example:
//
//	ExtractPathSegment("/users/123/roles/456", 0) -> "users"
//	ExtractPathSegment("/users/123/roles/456", 1) -> "123"
//	ExtractPathSegment("/users/123/roles/456", 2) -> "roles"
//	ExtractPathSegment("/users/123/roles/456", 3) -> "456"
func ExtractPathSegment(path string, index int) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return ""
	}
	segments := strings.Split(path, "/")
	if index < 0 || index >= len(segments) {
		return ""
	}
	return segments[index]
}

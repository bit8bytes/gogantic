package tools

// Input holds the data passed to a tool. It is intentionally kept as a struct
// rather than a plain string to remain extensible for future fields (e.g. metadata, context).
type Input struct {
	Content string
}

// Output holds the data returned by a tool. Like Input, it is a struct to allow
// additional fields to be added in the future without breaking existing tool implementations.
type Output struct {
	Content string
}

package set

// A simple set that can add remove and get n random elements
type IStringSet interface {
	// Add given string to set
	Add(v string)
	// Remove a string from set
	Remove(v string)
	// Get n random strings from the set
	RandomN(n int) []string
}

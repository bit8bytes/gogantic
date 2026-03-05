// Package roles defines the standard message roles used in large language model (LLM) conversations.
package roles

// Role represents a participant role in an LLM conversation.
// Common roles include System, User, and Assistant, which correspond
// to the standard roles supported by most LLM APIs.
type Role string

// Common LLM conversation roles.
const (
	// System defines instructions or context that guide the assistant's behavior.
	System Role = "system"
	// User represents a message from the human participant in the conversation.
	User Role = "user"
	// Assistant represents a message from the LLM in the conversation.
	Assistant Role = "assistant"
)

func (r Role) String() string {
	return string(r)
}

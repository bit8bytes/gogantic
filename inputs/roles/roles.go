package roles

type Role string

const (
	System    Role = "system"
	User      Role = "user"
	Assistant Role = "assistant"
)

func (r Role) String() string {
	return string(r)
}

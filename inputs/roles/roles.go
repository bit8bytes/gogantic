package roles

type Role string

const (
	System    Role = "system"
	User      Role = "user"
	Assistent Role = "assistent"
)

func (r Role) String() string {
	return string(r)
}

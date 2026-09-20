package ast

// Program represents a `gotiny` program containing a sequence of statements.
type Program struct {
	Statements []Stmt
}

func (p Program) String() string {
	var s string
	for i, stmt := range p.Statements {
		if i != 0 {
			s = s + "\n"
		}
		s += stmt.String() + ";"

	}
	return s
}

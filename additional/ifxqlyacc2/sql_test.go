package ifxqlyacc

import (
	"fmt"
	"strings"
	"testing"
)

type myVisitor struct {
}

func (mv *myVisitor) Visit(node Node) Visitor {
	Walk(mv, node)
	return mv
}
func TestYyParser(t *testing.T) {
	//s := NewScanner(strings.NewReader("SELECT value as a from myseries WHERE a = 'b"))
	tokenizer := &Tokenizer{
		query: Query{},
		//scanner:NewScanner(strings.NewReader("select *  From b where a = 1 order by time")),
	}
	tokenizer.scanner = NewScanner(strings.NewReader("select *  From b where a = 1 order by time"))
	yyParse(tokenizer)
	node := tokenizer.query.Statements[0]
	fmt.Printf("%v", node)
	fmt.Println()
	mv := &myVisitor{}
	Walk(mv, node)
}

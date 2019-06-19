%{
package ifxqlyacc
import (
    "time"
    "vitess.io/vitess/go/vt/log"
	"strconv"
	"fmt"
	"github.com/pkg/errors"
)
func setParseTree(yylex interface{},stmt Statement){
    yylex.(*Tokenizer).query.Statements = append(yylex.(*Tokenizer).query.Statements,stmt)
}
%}
%union{
    stmt                Statement
    stmts               Statements
    selStmt             *SelectStatement
    sdbStmt             *ShowDatabasesStatement
    cdbStmt             *CreateDatabaseStatement
    smmStmt             *ShowMeasurementsStatement
    str                 string
    query               Query
    field               *Field
    fields              Fields
    sources             Sources
    sortfs              SortFields
    sortf               *SortField
    ment                *Measurement
    dimens              Dimensions
    dimen               *Dimension
    int                 int
    int64               int64
    float64             float64
    expr                Expr
    tdur                time.Duration
    bool                bool
}
%token <str>    SELECT FROM WHERE AS GROUP BY ORDER LIMIT SHOW CREATE
%token <str>    DATABASES DATABASE MEASUREMENTS
%token <str>    COMMA SEMICOLON
%token <int>    MUL
%token <int>    EQ NEQ LT LTE GT GTE
%token <str>    IDENT
%token <int64>  INTEGER
%token <tdur>   DURATIONVAL
%token <str>    STRING
%token <bool>   DESC ASC
%token <float64> NUMBER
%left <int> AND OR
%type <stmt>                        STATEMENT
%type <sdbStmt>                     SHOW_DATABASES_STATEMENT
%type <cdbStmt>                     CREATE_DATABASE_STATEMENT
%type <selStmt>                     SELECT_STATEMENT
%type <smmStmt>                     SHOW_MEASUREMENTS_STATEMENT
%type <fields>                      COLUMN_NAMES
%type <field>                       COLUMN_NAME
%type <stmts>                       ALL_QUERIES
%type <sources>                     FROM_CLAUSE TABLE_NAMES
%type <ment>                        TABLE_NAME
%type <dimens>                      DIMENSION_NAMES GROUP_BY_CLAUSE
%type <dimen>                       DIMENSION_NAME
%type <expr>                        WHERE_CLAUSE CONDITION CONDITION_VAR OPERATION_EQUAL
%type <int>                         OPER LIMIT_INT
%type <sortfs>                      SORTFIELDS ORDER_CLAUSES
%type <sortf>                       SORTFIELD
%%
ALL_QUERIES:
        STATEMENT
        {
            setParseTree(yylex, $1)
        }
        | STATEMENT SEMICOLON
        {
            setParseTree(yylex, $1)
        }
        | STATEMENT SEMICOLON ALL_QUERIES
        {
            setParseTree(yylex, $1)
        }
STATEMENT:
    SELECT_STATEMENT
    {
        $$ = $1
    }
    |SHOW_DATABASES_STATEMENT
    {
        $$ = $1
    }
    |CREATE_DATABASE_STATEMENT
    {
        $$ = $1
    }
    |SHOW_MEASUREMENTS_STATEMENT
    {
        $$ = $1
    }
SELECT_STATEMENT:
    //SELECT COLUMN_NAMES
    //SELECT COLUMN_NAMES FROM_CLAUSE GROUP_BY_CLAUSE WHERE_CLAUSE ORDER_CLAUSES INTO_CLAUSE
    SELECT COLUMN_NAMES FROM_CLAUSE GROUP_BY_CLAUSE WHERE_CLAUSE ORDER_CLAUSES LIMIT_INT
    {
        sel := &SelectStatement{}
        sel.Fields = $2
        //sel.Target = $7
        sel.Sources = $3
        sel.Dimensions = $4
        sel.Condition = $5
        sel.SortFields = $6
        sel.Limit = $7
        $$ = sel
    }
COLUMN_NAMES:
    COLUMN_NAME
    {
        $$ = []*Field{$1}
    }
    |COLUMN_NAME COMMA COLUMN_NAMES
    {
        $$ = append($3,$1)
    }
COLUMN_NAME:
    MUL
    {
        $$ = &Field{Expr:&Wildcard{Type:$1}}
    }
    |IDENT
    {
        $$ = &Field{Expr:&VarRef{Val:$1}}
    }
    |IDENT AS IDENT
    {
        $$ = &Field{Expr:&VarRef{Val:$1},Alias:$3}
    }
FROM_CLAUSE:
    FROM TABLE_NAMES
    {
        $$ = $2
    }
    |
    {
        $$ = nil
    }
TABLE_NAMES:
    TABLE_NAME
    {
        $$ = []Source{$1}
    }
    |TABLE_NAME COMMA TABLE_NAMES
    {
        $$ = append($3,$1)
    }
TABLE_NAME:
    IDENT
    {
        $$ = &Measurement{Name:$1}
    }
GROUP_BY_CLAUSE:
    GROUP BY DIMENSION_NAMES
    {
        $$ = $3
    }
    |
    {
        $$ = nil
    }
DIMENSION_NAMES:
    DIMENSION_NAME
    {
        $$ = []*Dimension{$1}
    }
    |DIMENSION_NAME COMMA DIMENSION_NAMES
    {
        $$ = append($3,$1)
    }
DIMENSION_NAME:
    IDENT
    {
        $$ = &Dimension{Expr:&VarRef{Val:$1}}
    }
WHERE_CLAUSE:
    WHERE CONDITION
    {
        $$ = $2
    }
    |
    {
        $$ = nil
    }
CONDITION:
    OPERATION_EQUAL
    {
        $$ = $1
    }
    |CONDITION AND CONDITION
    {
        $$ = &BinaryExpr{Op:$2,LHS:$1,RHS:$3}
    }
    |CONDITION OR CONDITION
    {
        $$ = &BinaryExpr{Op:$2,LHS:$1,RHS:$3}
    }
OPERATION_EQUAL:
    CONDITION_VAR OPER CONDITION_VAR
    {
        $$ = &BinaryExpr{Op:$2,LHS:$1,RHS:$3}
    }
OPER:
    EQ
    {
        $$ = $1
    }
    |NEQ
    {
        $$ = $1
    }
    |LT
    {
        $$ =$1
    }
    |LTE
    {
        $$ = $1
    }
    |GT
    {
        $$ = $1
    }
    |GTE
    {
        $$ = $1
    }
CONDITION_VAR:
    IDENT
    {
        $$ = &VarRef{Val:$1}
    }
    |NUMBER
    {
        $$ = &NumberLiteral{Val:$1}
    }
    |INTEGER
    {
        $$ = &IntegerLiteral{Val:$1}
    }
    |DURATIONVAL
    {
        $$ = &DurationLiteral{Val:$1}
    }
    |STRING
    {
        $$ = &StringLiteral{Val:$1}
    }
ORDER_CLAUSES:
    ORDER BY SORTFIELDS
    {
        $$ = $3
    }
    |
    {
        $$ = nil
    }
SORTFIELDS:
    SORTFIELD
    {
        $$ = []*SortField{$1}
    }
    |SORTFIELD COMMA SORTFIELDS
    {
        $$ = append($3,$1)
    }
SORTFIELD:
    IDENT
    {
        $$ = &SortField{Name:$1}
    }
    |IDENT DESC
    {
        $$ = &SortField{Name:$1,Ascending:$2}
    }
    |IDENT ASC
    {
        $$ = &SortField{Name:$1,Ascending:$2}
    }
LIMIT_INT:
    LIMIT INTEGER
    {
        $$ = int($2)
    }
    |
    {
        $$ = 0
    }
SHOW_DATABASES_STATEMENT:
    SHOW DATABASES
    {
        $$ = &ShowDatabasesStatement{}
    }
CREATE_DATABASE_STATEMENT:
    CREATE DATABASE IDENT
    {
        $$ = &CreateDatabaseStatement{Name:$3}
    }
SHOW_MEASUREMENTS_STATEMENT:
    SHOW MEASUREMENTS WHERE_CLAUSE ORDER_CLAUSES LIMIT_INT
    {
        sms := &ShowMeasurementsStatement{}
        sms.Condition = $3
        sms.SortFields = $4
        sms.Limit = $5
        $$ = sms
    }
%%
type Tokenizer struct {
	query Query
	scanner *Scanner
}

func (tkn *Tokenizer) Lex(lval *yySymType) int{
	var typ int
	var val string

	for {
		typ, _, val  = tkn.scanner.Scan()
		if typ == EOF{
			return 0
		}
		if typ == MUL{
			lval.int = ILLEGAL
			//break
		}
		if typ >= EQ && typ <= GTE{
			//println("oo string is ",val)
			//lval.int,_ = strconv.Atoi(val)
			//println("oo is ",lval.int)
			lval.int = typ
			//break
		}
		if typ == NUMBER{
			lval.float64, _ = strconv.ParseFloat(val,64)
			//break
		}
		if typ == INTEGER{
			lval.int64, _ = strconv.ParseInt(val,10,64)
			//break
		}
		if typ == DURATIONVAL{
			time,_ := ParseDuration(val)
			lval.tdur = time
			//break
		}
		if typ == DESC{
			lval.bool = false
			//break
		}
		if typ == AND{
			lval.int = AND
			//break
		}
		if typ == OR{
			lval.int = OR
			//break
		}
		if typ == ASC{
			lval.bool = true
			//break
		}
		if typ !=WS{
			break
		}
	}
	lval.str = val
	return typ
}
func (tkn *Tokenizer) Error(err string){
	log.Fatal(err)
}

var ErrInvalidDuration = errors.New("invalid duration")
// ParseDuration parses a time duration from a string.
// This is needed instead of time.ParseDuration because this will support
// the full syntax that InfluxQL supports for specifying durations
// including weeks and days.
func ParseDuration(s string) (time.Duration, error) {
	// Return an error if the string is blank or one character
	if len(s) < 2 {
		return 0, ErrInvalidDuration
	}

	// Split string into individual runes.
	a := []rune(s)

	// Start with a zero duration.
	var d time.Duration
	i := 0

	// Check for a negative.
	isNegative := false
	if a[i] == '-' {
		isNegative = true
		i++
	}

	var measure int64
	var unit string

	// Parsing loop.
	for i < len(a) {
		// Find the number portion.
		start := i
		for ; i < len(a) && isDigit(a[i]); i++ {
			// Scan for the digits.
		}

		// Check if we reached the end of the string prematurely.
		if i >= len(a) || i == start {
			return 0, ErrInvalidDuration
		}

		// Parse the numeric part.
		n, err := strconv.ParseInt(string(a[start:i]), 10, 64)
		if err != nil {
			return 0, ErrInvalidDuration
		}
		measure = n

		// Extract the unit of measure.
		// If the last two characters are "ms" then parse as milliseconds.
		// Otherwise just use the last character as the unit of measure.
		unit = string(a[i])
		switch a[i] {
		case 'n':
			if i+1 < len(a) && a[i+1] == 's' {
				unit = string(a[i : i+2])
				d += time.Duration(n)
				i += 2
				continue
			}
			return 0, ErrInvalidDuration
		case 'u', 'µ':
			d += time.Duration(n) * time.Microsecond
		case 'm':
			if i+1 < len(a) && a[i+1] == 's' {
				unit = string(a[i : i+2])
				d += time.Duration(n) * time.Millisecond
				i += 2
				continue
			}
			d += time.Duration(n) * time.Minute
		case 's':
			d += time.Duration(n) * time.Second
		case 'h':
			d += time.Duration(n) * time.Hour
		case 'd':
			d += time.Duration(n) * 24 * time.Hour
		case 'w':
			d += time.Duration(n) * 7 * 24 * time.Hour
		default:
			return 0, ErrInvalidDuration
		}
		i++
	}

	// Check to see if we overflowed a duration
	if d < 0 && !isNegative {
		return 0, fmt.Errorf("overflowed duration %d%s: choose a smaller duration or INF", measure, unit)
	}

	if isNegative {
		d = -d
	}
	return d, nil
}
package parser

import "testing"

// TestNamedTupleColumn reads back a column type ClickHouse itself printed.
//
// SHOW CREATE TABLE returns a tuple in whichever form the column was declared, and
// the named form is a column list, not a list of types. Reading only the unnamed
// form stopped the parse at the first element's TYPE — "expected the last token
// kind is: )" pointing at `label String` — which made an entire real table
// unreadable to anything that parses its schema back.
func TestNamedTupleColumn(t *testing.T) {
	const ddl = "CREATE TABLE event.log (`el` Tuple(label String, role LowCardinality(String), path Array(String))) ENGINE = MergeTree ORDER BY tuple()"
	stmts, err := NewParser(ddl).ParseStmts()
	if err != nil {
		t.Fatalf("named tuple did not parse: %v", err)
	}
	ct, ok := stmts[0].(*CreateTable)
	if !ok {
		t.Fatalf("statement is %T, want *CreateTable", stmts[0])
	}
	if n := len(ct.TableSchema.Columns); n != 1 {
		t.Fatalf("columns = %d, want 1", n)
	}
}

// TestUnnamedTupleColumnStillParses is the other half: the lookahead must not
// send a bare list of types down the column-list grammar.
func TestUnnamedTupleColumnStillParses(t *testing.T) {
	const ddl = "CREATE TABLE t (`x` Tuple(String, Array(String), LowCardinality(String))) ENGINE = MergeTree ORDER BY tuple()"
	if _, err := NewParser(ddl).ParseStmts(); err != nil {
		t.Fatalf("unnamed tuple regressed: %v", err)
	}
}

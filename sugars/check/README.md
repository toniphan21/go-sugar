## check 

### Syntax

```ebnf
IdentifierLHS = IdentifierList ":=" .
CallExpr = SelectorPath CallSuffix .

Check = [ IdentifierLHS ] "check" CallExpr .
```

![syntax](./imgs/syntax.svg)

```
--- /tools/railroad-diagram

Diagram(
  Start({type:'complex'}),
  Optional('IdentifierLHS'),
  Stack("check"),
  Stack('CallExpr'),
  End({type:'complex'})
)
```

#### Components

##### *IdentifierList*
![IdentifierList](/lex/gn/imgs/identifier-list.svg)

##### *IdentifierLHS*
![IdentifierLHS](/lex/imgs/identifier-lhs.svg)

##### *SelectorPath*
![SelectorPath](/lex/imgs/selector-path.svg)

##### *CallSuffix*
![CallSuffix](/lex/imgs/call-suffix.svg)

##### *CallExpr*
![CallExpr](/lex/imgs/call-expr.svg)

#### Examples
```
func something() error {
    return nil
}

// sugar
check something()

// desugar
err := something()
if err != nil {
    return [<zero>, ...] err
}

// sugar
x := check strconv.Atoi("123")

// desugar
x, err := strconv.Atoi("123")
if err != nil {
    return [<zero>, ...] err
}
```

### Lexical state machine

![railroad-diagram](./imgs/lex.svg)

```
--- /tools/railroad-diagram

Diagram(
  Start({type:'complex'}),
  NonTerminal('doBegin'),
  Optional(Stack('IdentifierLHS', NonTerminal('doCollectIdentifiers'))),
  Comment('expect-check'),
  Stack(NonTerminal('doCollectCheckPos'), "check"),
  Comment('expect-expr'),
  Stack(NonTerminal('doCollectCheckEnd'), 'CallExpr', NonTerminal('doCollectEnd')),
  End({type:'complex'})
)
```

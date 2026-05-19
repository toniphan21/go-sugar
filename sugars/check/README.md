## check 

### Syntax

```ebnf
IdentifierLHS = IdentifierList ":=" .
CallExpr = SelectorPath CallSuffix .

Check = [ IdentifierLHS ] "check" CallExpr .
```

![syntax](./imgs/syntax.svg)

*syntax*

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

#### components

![IdentifierList](/lex/gn/imgs/identifier-list.svg)

*IdentifierList*

![IdentifierLHS](/lex/imgs/identifier-lhs.svg)

*IdentifierLHS*

![SelectorPath](/lex/imgs/selector-path.svg)

*SelectorPath*

![CallSuffix](/lex/imgs/call-suffix.svg)

*CallSuffix*

![CallExpr](/lex/imgs/call-expr.svg)

*CallExpr*

#### examples
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

## check 

### Syntax

the `check` sugar syntax is the sugar for error handling, it looks like `[var] [assign] check <expr>`

examples
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

```mermaid
stateDiagram-v2
    Idle --> Start: Boundary<br/>
    Idle --> Idle: Any<br/>nil
    Start --> ExpectCheck: IDENT("check")<br/>nil
    Start --> LHS: IDENT<br/>nil
    Start --> Idle: Any<br/>reset()
    LHS --> LHS: IDENT<br/>appendVariable()
    LHS --> LHS: COMMA<br/>nil
    LHS --> ExpectCheck: ASSIGN<br/>setOpAssign()
    LHS --> ExpectCheck: DEFINE<br/>setOpDefine()
    LHS --> Idle: Any<br/>reset()
    ExpectCheck --> Expr: IDENT("check")<br/>nil
    ExpectCheck --> Idle: Any<br/>reset()
    Expr --> Expr: IDENT<br/>appendOperand()
    Expr --> Expr: PERIOD<br/>appendOperand()
    Expr --> ExprIgnore: LPAREN<br/>nil
    Expr --> End: Boundary<br/>
    Expr --> Idle: Any<br/>reset()
    ExprIgnore --> ExprIgnore: any<br/>nil
    ExprIgnore --> Expr: RPAREN<br/>nil
    ExprIgnore --> End: Boundary<br/>nil
```

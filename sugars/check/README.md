## check 

### Syntax

```ebnf
CheckStmt         = [ CheckResult ] "check" Expression

CheckResult       = CheckShortVarDecl | CheckVarDecl | CheckAssignment .
CheckShortVarDecl = IdentifierList ":=" .
CheckVarDecl      = "var" IdentifierList [ Type ] "=" .
CheckAssignment   = ExpressionList "=" .
```

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
    Start --> Target: IDENT<br/>nil
    Start --> Idle: Any<br/>reset()
    Target --> Target: IDENT<br/>appendVariable()
    Target --> Target: COMMA<br/>nil
    Target --> ExpectCheck: ASSIGN<br/>setOpAssign()
    Target --> ExpectCheck: DEFINE<br/>setOpDefine()
    Target --> Idle: Any<br/>reset()
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

### WIP

````
full form

```ebnf
CheckStmt         = [ CheckResult ] "check" Expression [ "handle" CheckHandlerExpr ] .

CheckResult       = CheckShortVarDecl | CheckVarDecl | CheckAssignment .
CheckShortVarDecl = IdentifierList ":=" .
CheckVarDecl      = "var" IdentifierList [ Type ] "=" .
CheckAssignment   = ExpressionList "=" .

CheckHandlerExpr  = Expression . /* must have type func(error) error */
```
````

````
// group and re-use state machine - WIP

```mermaid
stateDiagram-v2
    state Target {
        [*] --> TargetCollect: IDENT
        TargetCollect --> TargetCollect: IDENT<br/>appendVariable()
        TargetCollect --> TargetCollect: COMMA
        TargetCollect --> [*]: ASSIGN<br/>setOpAssign()
        TargetCollect --> [*]: DEFINE<br/>setOpDefine()
    }
    
    Idle --> Start: Boundary<br/>
    Idle --> Idle: Any<br/>nil
    Start --> ExpectCheck: IDENT("check")<br/>nil
    Start --> Target: IDENT<br/>nil
    Start --> Idle: Any<br/>reset()
    Target --> ExpectCheck: ASSIGN<br/>setOpAssign()
    Target --> ExpectCheck: DEFINE<br/>setOpDefine()
    Target --> Idle: Any<br/>reset()
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
````

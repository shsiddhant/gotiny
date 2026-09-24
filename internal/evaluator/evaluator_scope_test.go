package evaluator

import (
	"fmt"
	"testing"

	"github.com/shsiddhant/gotiny/internal/checker"
	"github.com/shsiddhant/gotiny/internal/environment"
	"github.com/shsiddhant/gotiny/internal/objects"
	"github.com/shsiddhant/gotiny/internal/parser"
)

func TestEvalIfElseMutation(t *testing.T) {
	progString := `
        let x = 2000;
        if x > 1729 {
            x = x - 1729;
        } else {
            x = 1729 - x;
        }
    `

	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}

	env := environment.NewEnvironment()
	typEnv := checker.NewTypeEnvironment()

	_, err = EvalProgram(program, env, typEnv)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("variable 'x' not found in environment")
	}

	if got != objects.Int(271) {
		t.Errorf("got %d, expected 271", got)
	}
}

func TestEvalVariableShadowing(t *testing.T) {
	progString := `
        let x = 1712;
        if true {
            let x = -1729;
        }
    `

	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}
	env := environment.NewEnvironment()
	typEnv := checker.NewTypeEnvironment()

	_, err = EvalProgram(program, env, typEnv)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("variable 'x' not found in environment")
	}

	if got != objects.Int(1712) {
		t.Errorf("expected outer x to stay 1712, got %v", got)
	}
}

func TestEvalNestedMutation(t *testing.T) {
	progString := `
        let x = 1712;
        if true {
            if true {
                x = 1205;
            }
        }
    `

	program, err := parser.Parse(progString)
	if err != nil {
		t.Fatal(err)
	}

	env := environment.NewEnvironment()
	typeEnv := checker.NewTypeEnvironment()

	_, err = EvalProgram(program, env, typeEnv)
	if err != nil {
		t.Fatal(err)
	}

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("variable 'x' not found in environment")
	}

	if got != objects.Int(1205) {
		t.Errorf("got %v, expected 1205", got)
	}
}

func TestEvalFunctions(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected objects.Value
	}{
		{
			name: "Simple function call",
			source: `
        		fn add(a Int, b Int) Int {
            		return a + b;
        		}
        		add(1712, -1729);
    		`,
			expected: objects.Int(-17),
		},
		{
			name: "Negative arguments",
			source: `
        		fn difference(a Int, b Int) Int {
            		return a - b;
        		}
        		difference(-1712, 1205);
    		`,
			expected: objects.Int(-2917),
		},
		{
			name: "Multiple parameters and comparisons",
			source: `
		        fn check(a Int, b Int, c Int) Bool {
		            return a < b && b < c;
		        }

		        check(-1712, 1013, 1729);
		    `,
			expected: objects.Bool(true),
		},
		{
			name: "Function call inside expression",
			source: `
		        fn square(x Int) Int {
		            return x * x;
		        }

		        square(41) + square(2);
		    `,
			expected: objects.Int(1685),
		},
		{
			name: "Function result used in comparison",
			source: `
		        fn difference(a Int, b Int) Int {
		            return a - b;
		        }

		        difference(1729, 1013) > 1205;
		    `,
			expected: objects.Bool(false),
		},
		{
			name: "Function with conditional return",
			source: `
		        fn abs(x Int) Int {
		            if x < 0 {
		                return -x;
		            } else {
		                return x;
		            }
		        }

		        abs(-1712);
		    `,
			expected: objects.Int(1712),
		},
		{
			name: "Conditional mutation inside function",
			source: `
		        fn classify(x Int) Int {
		            let result = 0;

		            if x > 1205 {
		                result = 1729;
		            } else {
		                result = 1013;
		            }

		            return result;
		        }

		        classify(1712);
		    `,
			expected: objects.Int(1729),
		},
		{
			name: "Function can mutate global",
			source: `
		        let x = 1712;

		        fn update() {
		            x = x + 17;
		        }

		        update();
		        x;
		    `,
			expected: objects.Int(1729),
		},
		{
			name: "Function local mutation",
			source: `
		        fn f() Int {
		            let x = 1013;
		            x = x + 192;
		            return x;
		        }

		        f();
		    `,
			expected: objects.Int(1205),
		},
		{
			name: "Closure mutation is allowed",
			source: `
				fn outer() Int {
				    let x = 1013;

				    fn increment() {
				        x = x + 1;
				    }

				    increment();
				    increment();

				    return x;
				}
				outer();
			`,
			expected: objects.Int(1015),
		},
		{
			name: "Nested block can mutate function local",
			source: `
		        fn f() Int {
		            let x = 1116;

		            if true {
		                x = 1205;
		            }

		            return x;
		        }

		        f();
		    `,
			expected: objects.Int(1205),
		},
		{
			name: "Nested block shadowing",
			source: `
		        fn f(x Int) Int {
		            if true {
		                let x = 1205;
		                return x;
		            }

		            return x;
		        }

		        f(1712);
		    `,
			expected: objects.Int(1205),
		},
		{
			name: "Closure reads defining environment",
			source: `
				fn outer(x Int) Int {
					fn inner(y Int) Int {
						return x + y;
					}

					return inner(1205);
				}

				outer(1712);
		    `,
			expected: objects.Int(2917),
		},
		{
			name: "Closure ignores caller shadowing",
			source: `
		        fn outer(x Int) Int {
		            fn inner(y Int) Int {
		                return x + y;
		            }

		            if true {
		                let x = 1013;
		                return inner(1205);
		            }
					return 0;
		        }

		        outer(1712);
		    `,
			expected: objects.Int(2917),
		},
		{
			name: "Nested function calls",
			source: `
		        fn add(a Int, b Int) Int {
		            return a + b;
		        }

		        fn multiply(a Int, b Int) Int {
		            return a * b;
		        }

		        multiply(add(1712, -1729), 2);
		    `,
			expected: objects.Int(-34),
		},
		{
			name: "Recursive factorial",
			source: `
		        fn factorial(n Int) Int {
		            if n <= 1 {
		                return 1;
		            } else {
		                return n * factorial(n - 1);
		            }
		        }

		        factorial(6);
		    `,
			expected: objects.Int(720),
		},
		{
			name: "Recursive fibonacci",
			source: `
		        fn fib(n Int) Int {
		            if n <= 1 {
		                return n;
		            } else {
		                return fib(n - 1) + fib(n - 2);
		            }
		        }

		        fib(11);
		    `,
			expected: objects.Int(89),
		},
		{
			name: "Recursive function with negative input",
			source: `
		        fn fib(n Int) Int {
		            if n < 0 {
		                return 0;
		            }

		            if n <= 1 {
		                return n;
		            } else {
		                return fib(n - 1) + fib(n - 2);
		            }
		        }

		        fib(-1712);
		    `,
			expected: objects.Int(0),
		},
		{
			name: "Boolean return from function",
			source: `
	        	fn isLarge(x Int) Bool {
	            	return x > 1205;
	        	}

	        	isLarge(1729);
	    	`,
			expected: objects.Bool(true),
		},
		{
			name: "Boolean function with compound condition",
			source: `
		        fn interesting(x Int) Bool {
		            return x > 1013 && x < 1729 || x == 1116;
		        }

		        interesting(1205);
		    `,
			expected: objects.Bool(true),
		},
		{
			name: "Function result used with unary operator",
			source: `
		        fn isSmall(x Int) Bool {
		            return x < 1205;
		        }

		        !isSmall(1712);
		    `,
			expected: objects.Bool(true),
		},
		{
			name: "Function arguments are expressions",
			source: `
		        fn difference(a Int, b Int) Int {
		            return a - b;
		        }

		        difference(1729 - 1205, 1116 - 1013);
		    `,
			expected: objects.Int(421),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d:%s", i, tt.name), func(t *testing.T) {
			program, err := parser.Parse(tt.source)
			if err != nil {
				t.Fatalf("unexpected parser error: %v", err)
			}

			env := environment.NewEnvironment()
			typeEnv := checker.NewTypeEnvironment()

			got, err := EvalProgram(program, env, typeEnv)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Errorf("got %v, expected %v", got, tt.expected)
			}

		})
	}
}

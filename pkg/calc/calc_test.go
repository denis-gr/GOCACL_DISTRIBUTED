package calc

import (
	"testing"
)

func TestSplitExpression(t *testing.T) {
	cases := []struct {
		name       string
		expression string
		want       []string
	}{
		{
			name:       "simple expression",
			expression: "2+2",
			want:       []string{"2", "+", "2"},
		},
		{
			name:       "expression with parentheses",
			expression: "2*(3+4)",
			want:       []string{"2", "*", "(", "3", "+", "4", ")"},
		},
		{
			name:       "complex expression",
			expression: "2+2*(3+3*(1+2))",
			want:       []string{"2", "+", "2", "*", "(", "3", "+", "3", "*", "(", "1", "+", "2", ")", ")"},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := splitExpression(tc.expression)
			if len(got) != len(tc.want) {
				t.Errorf("splitExpression(%v) = %v; want %v", tc.expression, got, tc.want)
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("splitExpression(%v) = %v; want %v", tc.expression, got, tc.want)
					break
				}
			}
		})
	}
}

func TestPrecedence(t *testing.T) {
	cases := []struct {
		name string
		op   string
		want int
	}{
		{
			name: "addition",
			op:   "+",
			want: 1,
		},
		{
			name: "multiplication",
			op:   "*",
			want: 2,
		},
		{
			name: "unknown operator",
			op:   "^",
			want: 0,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := precedence(tc.op)
			if got != tc.want {
				t.Errorf("precedence(%v) = %v; want %v", tc.op, got, tc.want)
			}
		})
	}
}

func TestToRPN(t *testing.T) {
	cases := []struct {
		name       string
		expression []string
		want       []string
	}{
		{
			name:       "simple expression",
			expression: []string{"2", "+", "2"},
			want:       []string{"2", "2", "+"},
		},
		{
			name:       "expression with parentheses",
			expression: []string{"2", "*", "(", "3", "+", "4", ")"},
			want:       []string{"2", "3", "4", "+", "*"},
		},
		{
			name:       "complex expression",
			expression: []string{"2", "+", "2", "*", "(", "3", "+", "3", "*", "(", "1", "+", "2", ")", ")"},
			want:       []string{"2", "2", "3", "3", "1", "2", "+", "*", "+", "*", "+"},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := toRPN(tc.expression)
			if len(got) != len(tc.want) {
				t.Errorf("toRPN(%v) = %v; want %v", tc.expression, got, tc.want)
				return
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("toRPN(%v) = %v; want %v", tc.expression, got, tc.want)
					break
				}
			}
		})
	}
}

func TestCalc(t *testing.T) {
	cases := []struct {
		name   string
		values string
		want   float64
	}{
		{
			name:   "positive values",
			values: "2+2",
			want:   4.0,
		},
		{
			name:   "mixed values",
			values: "2-4",
			want:   -2.0,
		},
		{
			name:   "complex expression",
			values: "2+2*(3+3*(1+2))",
			want:   26.0,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, _ := Calc(tc.values)
			if got != tc.want {
				t.Errorf("Calc(%v) = %v; want %v", tc.values, got, tc.want)
			}
		})
	}
}

func TestCalc2(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		shouldFail bool
	}{
		{"", 0, true},
		{"2+2", 4, false},
		{"2*3+4", 10, false},
		{"2*(3+4)", 14, false},
		{"2/0", 0, true},
		{"2+2*(3+3*(1+2", 0, true},
		{"2+2*3+3*1+2", 13, false},
		{"2+2*3+3*(1+2)", 17, false},
		{"2+2*3+3*(1+2)+5/5", 18, false},
		{"2+2*3+3*(1+2)+5/0", 0, true},
		{"2+2*3+3*(1+2)+5/5-1", 17, false},
	}

	for _, test := range tests {
		result, err := Calc(test.expression)
		if test.shouldFail {
			if err == nil {
				t.Errorf("Calc(%q) = %v; expected an error", test.expression, result)
			}
		} else {
			if err != nil {
				t.Errorf("Calc(%q) returned an error: %v", test.expression, err)
			}
			if result != test.expected {
				t.Errorf("Calc(%q) = %v; want %v", test.expression, result, test.expected)
			}
		}
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"2+2", 4},
		{"2+2*2", 6},
		{"(2+2)*2", 8},
		{"2+2*(3+3*(1+2))", 20},
	}

	for _, test := range tests {
		result, err := Calc(test.expression)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != test.expected {
			t.Errorf("for expression %s, expected %v, got %v", test.expression, test.expected, result)
		}
	}
}

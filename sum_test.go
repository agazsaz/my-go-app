package main
import "testing"
func TestSum(t *testing.T) {
    tests := []struct { a, b, expected int }{ {2, 3, 5}, {-1, 1, 0}, {0, 0, 0} }
    for _, tt := range tests {
        result := Sum(tt.a, tt.b)
        if result != tt.expected { t.Errorf("Sum(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected) }
    }
}

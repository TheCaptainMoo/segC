fn fibonacci(n){
    if n == 0 {
        ret 0
    } elsif n == 1 {
        ret 1
    } else {
        ret fibonacci(n-1) + fibonacci(n-2)
    }
}

fn main() {
    ret fibonacci(7)
}
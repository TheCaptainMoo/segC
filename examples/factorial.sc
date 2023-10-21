fn factorial( n) {
    if n == 0 {
        ret 1
    }
    ret factorial(n - 1) * n;
}

fn main() {
    x = factorial(150)
    ret x
}
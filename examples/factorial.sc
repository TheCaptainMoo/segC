fn factorial(n)
{
    if n == 0
    {
        ret 1
    }
    ret factorial(n - 1) * n
}

fn main(arg0)
{
    x = factorial(arg0)
    ret x
}
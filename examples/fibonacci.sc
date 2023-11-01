fn fibonacci(n)
{
    if n < 2
    {
        ret n == 1
    }
    ret fibonacci(n-1) + fibonacci(n-2)
}

fn main(arg)
{
    ret fibonacci(arg)
}
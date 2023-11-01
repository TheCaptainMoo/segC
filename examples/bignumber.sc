fn big_number(n)
{
    if n < 1000
    {
        ret big_number(n + 1)
    }
    ret n
}

fn main()
{
    x = big_number(0)
    ret x
}
![img.png](resources/images/img.png)

![img.png](resources/images/for.png)

```lox
for (var a=0; a<10; a=a+1) {
    print a;
}

{
    var a = 0;
    while (a < 10) {
        {
            print a;
        }
        a = a + 1
    }
}
```
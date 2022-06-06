_🔨 GLox is an interpreter for a toy language "Lox" based on Go._

"Lox" is a toy language for learning, so it may not be very performant.

It combines functional and object-oriented features, and supports basic features such as arithmetic operations and control flow (you can see some Lox code in `resources/lox` folder).

The interpreter is not full-tested, so may be there are some bugs, once I find a bug, it will be fixed ASAP.

Usage:
```
cd ./cmd
go build -o GLox
./GLox -s "source_file_path"
```

---

todo: refactor error handing.

# DrawSudokus

This is a go command line utility drawing multiple sudokus on an A4 .pdf file. It leverages the `qqwing` sudoku generator.


## Usage
```
go run drawsudokus.go <command line options>
```

Viable command line options are:
```
  --help
        display this list
  -difficulty value
        one of simple, easy, intermediate, expert, any
  -nx int
        number of sudokus put horizontally (default 4)
  -ny int
        number of sudokus put vertically (default 3)
```

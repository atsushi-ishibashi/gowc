# gowc

## Installing
```
go get github.com/atsushi-ishibashi/gowc
```

## Usage
```
gowc
TotalStats
-------------------------
| Files | Lines | Bytes |
-------------------------
| 3     | 203   | 5037  |
-------------------------
Each files
--------------
| main.go    |
--------------
| 173 | 3538 |
--------------
-------------
| README.md |
-------------
| 1 | 6     |
-------------
-------------
| LICENSE   |
-------------
| 29 | 1493 |
-------------
```
```
// output specified file
gowc -f README.md
TotalStats
-------------------------
| Files | Lines | Bytes |
-------------------------
| 1     | 1     | 6     |
-------------------------
Each files
-------------
| README.md |
-------------
| 1 | 6     |
-------------
```
```
// exclude files
gowc -ex .*.md
TotalStats
-------------------------
| Files | Lines | Bytes |
-------------------------
| 2     | 202   | 5031  |
-------------------------
Each files
--------------
| main.go    |
--------------
| 173 | 3538 |
--------------
-------------
| LICENSE   |
-------------
| 29 | 1493 |
-------------
```

### Help
```
gowc -h

Usage of gowc:
  -ex string
    	regexp to exclude file name(optional)
  -f string
    	file path(optional)
```

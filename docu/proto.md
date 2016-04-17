# tiss lang concepts

### assignments, variables, objects

````
hash -> a
5 -> x
10 + x -> y
m[10] -> y
h["hash"] -> x
h.hash
````

### expressions

````
(x + 5) ^ (4 + x)
x // 2
x % 2
x / 2.3
x + 1
-x
!(-x)
1 -! x
1 +! x
3 - 3
e = 1
3 # e
[1 .. 4)
(-inf .. inf)
````

### selectors

````
x[4].block.x[5].y -> z
x[1 .. 4] -> x14
````

### blocks
````
BLOCK

END

BLOCK Block

END Block

BLOCK

END -> block
````

### versioning
````
UNIT Unit
VERSION 0.0.0

  IMPORT OtherUnit[2.3.4]  

INIT

CLOSE

END Unit
````

### statements

````
0 -> i
Block
Block(i, x, ret)
Block[0 -> x, ret -> z]
\Block 0 -> i
IF x THEN 0 -> i ELSIF ~x THEN 1 -> i ELSE -1 -> i END
WHILE x DO INC(i) END
REPEAT INC(i) UNTIL x;
CHOOSE x OF 1: Block1 OR 2: Block2 ELSE BlockElse END
````

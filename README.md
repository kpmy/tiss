# WebAssembly s-expr generator

For the good.

* [wasm specs](https://github.com/WebAssembly/spec/)

```
;; golang wasm generator github.com/kpmy/tiss
(module
	(type $t0
		(func))
	(func $fib
		(param $x i64)
		(result i64)
		(local $i i64)
		(return
			(i64.const 0)))
	(func $start
		(type $t0)
		(call $fib
			(i64.const 0)))
	(start $start))
```

*not so pretty :(*

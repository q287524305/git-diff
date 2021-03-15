package Const

const DiffFile = `
 diff --git a/test/test1.go b/test/test1.go
index 8001ad7..065789c 100644
--- a/test/test1.go
+++ b/test/test1.go
@@ -1,5 +1,7 @@
package test
 
-func test1() {
+import "fmt"
 
+func test1() {
+       fmt.Println("test1")
 }
`

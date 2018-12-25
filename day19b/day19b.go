package main

import "fmt"

/*
func main() {
	var var0, var1, var2, var4 int

	var0 = 1

	var1 = 983
	if var0 == 1 {
		var1 += 10550400
		var0 = 0
	}

	var2 = 1
	for {
		var4 = 1
		for {
			if (var2 * var4) == var1 {
				var0 += var2
			}
			var4++

			if var4 > var1 {
				break
			}
		}
		var2++
		if var2 > var1 {
			break // EXIT
		}
	}

	fmt.Println(var0)
}
*/

func main() {
	var var0, var1 int

	var0 = 1

	var1 = 983
	if var0 == 1 {
		var1 += 10550400
		var0 = 0
	}

	sum := 0
	for i := 1; i <= var1; i++ {
		if var1%i == 0 {
			sum += i
		}
	}

	fmt.Println(sum)
}

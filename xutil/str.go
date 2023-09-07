/*
 * @Date: 2023-09-07 16:40:46
 * @LastEditTime: 2023-09-07 16:40:47
 * @Description:
 */
package xutil

func GetLastChars(str string, num int) string {
	if len(str) > num {
		return str[len(str)-num:]
	}
	return str
}

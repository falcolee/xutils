package xfile

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"

	"github.com/falcolee/xutils/xgen"
	"github.com/falcolee/xutils/xtype"
)

const (
	// TypeAll list dir and file
	TypeAll int = iota
	// TypeDir list only dir
	TypeDir
	// TypeFile list only file
	TypeFile
)

// LsFile is list file info
type LsFile struct {
	Type int
	Path string
	Name string
}

// New opens a file for new and return fd
func New(fpath string) (*os.File, error) {
	return NewFile(fpath, false)
}

// NewFile opens a file and return fd
func NewFile(fpath string, isAppend bool) (*os.File, error) {
	dir, _ := filepath.Split(fpath)
	if dir != "" && !IsDir(dir) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	if isAppend {
		return os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	}

	return os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

// WriteText writes string data to file
func WriteText(fpath, text string) (err error) {
	return WriteByte(fpath, []byte(text))
}

// AppendText appends string data to file
func AppendText(fpath, text string) (err error) {
	return Append(fpath, []byte(text))
}

// Write writes bytes data to file
func WriteByte(fpath string, data []byte) (err error) {
	return writeFile(fpath, data, false)
}

// Append appends bytes data to file
func Append(fpath string, data []byte) (err error) {
	return writeFile(fpath, data, true)
}

// writeFile writes bytes data to file
func writeFile(fpath string, data []byte, isAppend bool) (err error) {
	fd, err := NewFile(fpath, isAppend)
	if err != nil {
		return
	}

	n, err := fd.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if e := fd.Close(); err == nil {
		err = e
	}

	return
}

// Temp ...
func Temp(filenames ...string) string {
	var filename string
	if len(filename) > 0 {
		filename = filenames[0]
	} else {
		filename = xgen.Nanoid()
	}
	return path.Join(os.TempDir(), filename)
}

// Read ...
func Read(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return ""
	}
	return string(b)
}

// Write ...
func Write(filepath, str string) error {
	return os.WriteFile(filepath, []byte(str), 0755)
}

// Copy ...
func Copy(src, dst string) error {
	return Write(dst, Read(src))
}

// Size ...
func Size(filepath string) int64 {
	f, err := os.Stat(filepath)
	if err != nil {
		return 0
	}
	return f.Size()
}

// SizeText ...
func SizeText(size int64) string {
	s := float64(size)
	units := []string{"B", "KB", "MB", "GB", "TB"}
	index := 0
	for s >= 1024 {
		s /= 1024
		index++
		if index == len(units)-1 {
			break
		}
	}
	return fmt.Sprintf("%.2f%s", s, units[index])
}

// LineCount ...
func LineCount(filepath string) int {
	f, err := os.Open(filepath)
	if err != nil {
		return 0
	}
	defer f.Close()
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	return count
}

// ReadLastNLines ...
func ReadLastNLines(filePath string, n int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// LineContent ...
func LineContent(filepath string, numbers ...int) map[int]string {
	f, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer f.Close()
	res := make(map[int]string)
	count := len(numbers)
	scanner := bufio.NewScanner(f)
	for number := 1; scanner.Scan(); number++ {
		if count == 0 || xtype.IsContains(numbers, number) {
			res[number] = scanner.Text()
		}
	}
	return res
}

// MineType ...
func MineType(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	return ReaderMineType(f)
}

// ReaderMineType ...
func ReaderMineType(r io.Reader) string {
	// 512 http/sniff.go sniffLen
	var buf [512]byte
	n, _ := io.ReadFull(r, buf[:])
	if n == 0 {
		return ""
	}
	return http.DetectContentType(buf[:n])
}

// WriteJSON write data to JSON file
func WriteJSON(filepath string, data interface{}, pretty ...bool) error {
	var (
		b   []byte
		err error
	)
	if len(pretty) > 0 {
		b, err = json.MarshalIndent(data, "", "    ")
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, b, 0664)
}

// ReadJSON read JSON file data
func ReadJSON(filepath string, data interface{}) error {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, data)
}

// Md5 file md5
func Md5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	var size int64 = 1024 * 1024
	hash := md5.New()
	if fi.Size() < size {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		hash.Write(data)
	} else {
		b := make([]byte, size)
		for {
			n, err := f.Read(b)
			if err != nil {
				break
			}
			hash.Write(b[:n])
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Shadow 根据当前文件名称取其影分身
// /tmp/abc.txt => /tmp/abc_1.txt
func Shadow(filepath string) string {
	ext := path.Ext(filepath)
	prefix := strings.TrimSuffix(filepath, ext)
	res := ""
	for i := 1; ; i++ {
		f := fmt.Sprintf("%s_%d%s", prefix, i, ext)
		if !IsExist(f) {
			res = f
			break
		}
	}
	return res
}

// ReverseRead 读取文件的最后 N 行
func ReverseRead(fileName string, lineNum uint) ([]string, error) {
	//打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	buff := make([]string, 0, 100)
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return buff, err
		}
		if char[0] == '\n' {
			offset-- //windows跳过'\r'
			if lineNum > 0 {
				buff = append(buff, lineStr)
				lineNum-- //到此读取完一行
			}
			lineStr = ""
			if lineNum == 0 {
				return buff, nil
			}
		} else {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	buff = append(buff, lineStr)
	return buff, nil
}

// ReadLines returns N lines of file
func ReadLines(fpath string, n int) (lines []string, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer fd.Close()
	nRead := 0
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		nRead++
		if n > 0 && nRead >= n {
			break
		}
	}

	err = scanner.Err()

	return
}

// ReadFirstLine returns first NOT empty line
func ReadFirstLine(fpath string) (line string, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if line != "" {
			return
		}
	}

	err = scanner.Err()

	return
}

// ReadLastLine returns last NOT empty line
func ReadLastLine(fpath string) (line string, err error) {
	fd, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		return
	}

	size := stat.Size()
	if size == 0 {
		return
	}

	var cursor int64
	data := make([]byte, 0)

	for {
		cursor--
		_, err = fd.Seek(cursor, io.SeekEnd)
		if err != nil {
			return
		}

		buf := make([]byte, 1)
		_, err = fd.Read(buf)
		if err != nil {
			return
		}

		if buf[0] != '\r' && buf[0] != '\n' {
			data = append([]byte{buf[0]}, data...)
		} else {
			if cursor != -1 && strings.TrimSpace(string(data)) != "" {
				break
			}
			data = make([]byte, 0)
		}

		if cursor == -size {
			break
		}
	}

	return string(data), nil
}

// ListDir lists dir without recursion
func ListDir(fpath string, ftype, n int) (ls []LsFile, err error) {
	if fpath == "" {
		fpath = "."
	}

	if !strings.HasSuffix(fpath, "/") {
		fpath += "/"
	}

	fd, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer fd.Close()
	fs, err := fd.Readdir(-1)
	if err != nil {
		return
	}

	for _, f := range fs {
		tpath := fpath + f.Name()
		if f.IsDir() {
			if ftype == TypeAll || ftype == TypeDir {
				ls = append(ls, LsFile{TypeDir, tpath, f.Name()})
				if n > 0 && len(ls) >= n {
					return
				}
			}
		} else {
			if ftype == TypeAll || ftype == TypeFile {
				ls = append(ls, LsFile{TypeFile, tpath, f.Name()})
				if n > 0 && len(ls) >= n {
					return
				}
			}
		}
	}

	return
}

// ListDirAll lists dir and children, filter by type, returns up to n
func ListDirAll(fpath string, ftype, n int) (ls []LsFile, err error) {
	if fpath == "" {
		fpath = "."
	}

	if !strings.HasSuffix(fpath, "/") {
		fpath += "/"
	}

	fd, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer fd.Close()
	fs, err := fd.Readdir(-1)
	if err != nil {
		return
	}

	for _, f := range fs {
		tpath := fpath + f.Name()
		if f.IsDir() {
			if ftype == TypeAll || ftype == TypeDir {
				ls = append(ls, LsFile{TypeDir, tpath, f.Name()})
				if n > 0 && len(ls) >= n {
					return
				}
			}
			tls, err := ListDirAll(tpath, ftype, n-len(ls))
			if err != nil {
				return ls, err
			}
			ls = append(ls, tls...)
			if n > 0 && len(ls) >= n {
				return ls, nil
			}
		} else {
			if ftype == TypeAll || ftype == TypeFile {
				ls = append(ls, LsFile{TypeFile, tpath, f.Name()})
				if n > 0 && len(ls) >= n {
					return
				}
			}
		}
	}

	return
}

// Chmod chmods to path without recursion
func Chmod(fpath string, mode os.FileMode) error {
	return os.Chmod(fpath, mode)
}

// ChmodAll chmods to path and children, returns the first error it encounters
func ChmodAll(root string, mode os.FileMode) error {
	return filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return Chmod(fpath, mode)
	})
}

// Chown chowns to path without recursion
func Chown(fpath string, uid, gid int) error {
	return os.Chown(fpath, uid, gid)
}

// ChownAll chowns to path and children, returns the first error it encounters
func ChownAll(root string, uid, gid int) error {
	return filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return Chown(fpath, uid, gid)
	})
}

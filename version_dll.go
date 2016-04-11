package windll

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	VersionDLL = &versionDLL{}
)

type versionDLL struct{}

func (v *versionDLL) ExtractProductVersion(filePath string) (string, error) {
	versionDll, err := syscall.LoadLibrary("Version.dll")
	if err != nil {
		return "", err
	}

	getVersionInfoSizeProc, err := syscall.GetProcAddress(versionDll, "GetFileVersionInfoSizeW")
	if err != nil {
		return "", err
	}

	getVersionInfoProc, err := syscall.GetProcAddress(versionDll, "GetFileVersionInfoW")
	if err != nil {
		return "", err
	}

	var nargs uintptr = 2
	var sizeInfoLength uintptr
	var dwHandle uintptr
	if ret, _, err := syscall.Syscall(
		uintptr(getVersionInfoSizeProc),
		nargs,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(filePath))),
		dwHandle,
		0); err != 0 {
		return "", fmt.Errorf("Call GetFileVersionInfoSizeW failed. For file '%s': %v", filePath, err)
	} else {
		sizeInfoLength = ret
	}

	nargs = 4
	data := make([]uint16, int(sizeInfoLength))
	if _, _, err := syscall.Syscall6(
		uintptr(getVersionInfoProc),
		nargs,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(filePath))),
		0,
		sizeInfoLength,
		uintptr(unsafe.Pointer(&data[0])),
		0,
		0); err != 0 {
		return "", fmt.Errorf("Call GetFileVersionInfoW failed. For file '%s': %v", filePath, err)
	}

	str := string(utf16.Decode(data))

	tmpRegexString := `ProductVersion[^0-9]*([0-9]+.[0-9]+.[0-9]+.[0-9]+)\w*`
	findProductVersionRegex := regexp.MustCompile(tmpRegexString)
	if !findProductVersionRegex.MatchString(str) {
		return "", fmt.Errorf("Unable to find the expected regex ('%s') inside string: %s", tmpRegexString, str)
	}
	submatches := findProductVersionRegex.FindStringSubmatch(str)

	return strings.Trim(submatches[1], " "), nil
}

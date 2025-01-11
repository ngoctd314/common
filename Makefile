gomodule:=github.com/ngoctd314/common

.PHONY: testfunc, testpkg, testfile, testcoveragepkg, testcovertotal, testall
# Description: test a function in a package
# Usage: make testfunc pkg=<path to package> func=<func name> 
# Example: make testfunc pkg=pkg/helper func=TestReplaceTemplateByTag
testfunc:
	go test -v -count=1 -run ${func} ${gomodule}/${pkg}

# Description: test all function in a package
# Usage: make testpkg pkg=<path to package>
# Example: make testpkg pkg=pkg/helper
testpkg:
	go test -v -shuffle=on -count=1 ${gomodule}/${pkg}

# Description: test all function in a file
# Usage: make testfile pkg=<path to package> file=<file name without extension> 
# Example: make testfile pkg=pkg/helper file=time_test 
define get_func_name
 	grep 'func Test' ${pkg}/${file}.go | \
	sed -E 's|func ||' | \
	sed -E 's|\((.*?)\)||' | \
	sed -E 's| \{||' | \
	paste -sd '|'
endef

testfile:
	go test -v -shuffle=on -count=1 -run "$(shell $(call get_func_name))" ${gomodule}/${pkg}

# Description: test coverage and show coverage result in browser  
# Usage: make testcoveragepkg pkg=<path to package>
# Example: make testcoveragepkg pkg=pkg/helper
testcoveragepkg:
	go test -cover -covermode=count -coverprofile=cover.out ${gomodule}/${pkg} -count=1  && \
	go tool cover -html=cover.out && \
	rm cover.out

testcovertotal:
	go test ./... -coverprofile cover.out && \
	go tool cover -func cover.out|grep "total:" && \
	rm cover.out

testall:
	go test -v ./...

.PHONY: lint
# tools: https://golangci-lint.run/welcome/install/
lint:
	golangci-lint run

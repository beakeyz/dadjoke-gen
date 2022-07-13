package build

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	GoOSWindows = "windows"
	GoOSLinux   = "linux"
)

var binaries = []string{"server", "client"}

func RunBuild() int {

	var opts BuildType = makeBuildTypeFromFlags()

	opts.version = "1.0"

	if opts.printGenVersion {
		fmt.Println(opts.version)
	}

	for _, cmd := range flag.Args() {
		switch cmd {
		case "setup":
			setup(opts.goos)

		case "server":
			if !opts.isDev {
				clean(opts)
			}

			if err := doBuild("server", "./pkg/cmd/server", opts); err != nil {
				log.Println(err)
				return 1
			}
		case "client":
			clean(opts)
			if err := doBuild("client", "./pkg/cmd/client", opts); err != nil {
				log.Println(err)
				return 1
			}
		case "build":
			for _, binary := range binaries {
				log.Println("building binaries", cmd)
				// Can't use filepath.Join here because filepath.Join calls filepath.Clean, which removes the `./` from this path, which upsets `go build`
				if err := doBuild(binary, fmt.Sprintf("./pkg/cmd/%s", binary), opts); err != nil {
					log.Println(err)
					return 1
				}
			}

		default:
			log.Println("Command not found ", cmd)
		}
	}
	return -1
}

func setup(goos string) {
	args := []string{"install", "-v"}
	if goos == GoOSWindows {
		args = append(args, "-buildmode=exe")
	}
	args = append(args, "./pkg/cmd/server")
	runPrint("go", args...)
}

func doBuild(binaryName, pkg string, opts BuildType) error {
	log.Println("building", binaryName, pkg)
	libcPart := ""
	if opts.libc != "" {
		libcPart = fmt.Sprintf("-%s", opts.libc)
	}
	binary := fmt.Sprintf("./bin/%s", binaryName)

	//don't include os/arch/libc in output path in dev environment
	if !opts.isDev {
		binary = fmt.Sprintf("./bin/%s-%s%s/%s", opts.goos, opts.goarch, libcPart, binaryName)
	}

	if opts.goos == GoOSWindows {
		binary += ".exe"
	}

	if !opts.isDev {
		rmr(binary, binary+".md5")
	}

	lf, err := ldflags(opts)
	if err != nil {
		return err
	}

	args := []string{"build", "-ldflags", lf}

	if opts.goos == GoOSWindows {
		// Work around a linking error on Windows: "export ordinal too large"
		args = append(args, "-buildmode=exe")
	}

	if len(opts.buildTags) > 0 {
		args = append(args, "-tags", strings.Join(opts.buildTags, ","))
	}

	args = append(args, "-o", binary)
	args = append(args, pkg)

	runPrint("go", args...)

	if opts.isDev {
		return nil
	}

	if err := setBuildEnv(opts); err != nil {
		return err
	}
	runPrint("go", "version")
	libcPart = ""
	if opts.libc != "" {
		libcPart = fmt.Sprintf("/%s", opts.libc)
	}
	fmt.Printf("Targeting %s/%s%s\n", opts.goos, opts.goarch, libcPart)

	// Create an md5 checksum of the binary, to be included in the archive for
	// automatic upgrades.
	return md5File(binary)
}

func ldflags(opts BuildType) (string, error) {
	buildStamp, err := buildStamp()
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	b.WriteString("-w")
	b.WriteString(fmt.Sprintf(" -X main.version=%s", opts.version))
	b.WriteString(fmt.Sprintf(" -X main.commit=%s", getGitSha()))
	b.WriteString(fmt.Sprintf(" -X main.buildstamp=%d", buildStamp))
	b.WriteString(fmt.Sprintf(" -X main.buildBranch=%s", getGitBranch()))
	if v := os.Getenv("LDFLAGS"); v != "" {
		b.WriteString(fmt.Sprintf(" -extldflags \"%s\"", v))
	}

	return b.String(), nil
}

func setBuildEnv(opts BuildType) error {
	if err := os.Setenv("GOOS", opts.goos); err != nil {
		return err
	}

	if opts.goos == GoOSWindows {
		// require windows >=7
		if err := os.Setenv("CGO_CFLAGS", "-D_WIN32_WINNT=0x0601"); err != nil {
			return err
		}
	}

	if opts.goarch != "amd64" || opts.goos != GoOSLinux {
		// needed for all other archs
		opts.cgo = true
	}

	if strings.HasPrefix(opts.goarch, "armv") {
		if err := os.Setenv("GOARCH", "arm"); err != nil {
			return err
		}

		if err := os.Setenv("GOARM", opts.goarch[4:]); err != nil {
			return err
		}
	} else {
		if err := os.Setenv("GOARCH", opts.goarch); err != nil {
			return err
		}
	}

	if opts.cgo {
		if err := os.Setenv("CGO_ENABLED", "1"); err != nil {
			return err
		}
	}

	if opts.gocc == "" {
		return nil
	}

	return os.Setenv("CC", opts.gocc)
}

func buildStamp() (int64, error) {
	// use SOURCE_DATE_EPOCH if set.
	if v, ok := os.LookupEnv("SOURCE_DATE_EPOCH"); ok {
		return strconv.ParseInt(v, 10, 64)
	}

	bs, err := runError("git", "show", "-s", "--format=%ct")
	if err != nil {
		return time.Now().Unix(), nil
	}

	return strconv.ParseInt(string(bs), 10, 64)
}

func clean(opts BuildType) {
	rmr("dist")
	rmr("tmp")
	rmr(filepath.Join(build.Default.GOPATH, fmt.Sprintf("pkg/%s_%s/github.com/beakeyz", opts.goos, opts.goarch)))
}

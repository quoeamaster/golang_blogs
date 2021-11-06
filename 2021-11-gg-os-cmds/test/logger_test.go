package test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"

	g "github.com/quoeamaster/golang_blogs/ggoscmds"
)

func Test_logToFile(t *testing.T) {
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vel risus commodo viverra maecenas accumsan lacus vel facilisis volutpat. Neque egestas congue quisque egestas. Diam phasellus vestibulum lorem sed risus. Sed egestas egestas fringilla phasellus. Ut venenatis tellus in metus vulputate eu scelerisque. At consectetur lorem donec massa sapien faucibus et molestie ac. Suspendisse ultrices gravida dictum fusce ut placerat orci nulla. Aliquet sagittis id consectetur purus ut faucibus pulvinar. At elementum eu facilisis sed. Egestas erat imperdiet sed euismod nisi porta. Rutrum quisque non tellus orci ac auctor augue mauris augue. Sit amet nisl purus in mollis nunc sed id.",
		"At varius vel pharetra vel turpis nunc eget lorem. Lobortis elementum nibh tellus molestie nunc non blandit. Ac turpis egestas maecenas pharetra convallis posuere morbi leo. Aliquet nibh praesent tristique magna sit amet purus. Dolor sit amet consectetur adipiscing elit pellentesque habitant morbi. Facilisi etiam dignissim diam quis enim. Euismod elementum nisi quis eleifend quam adipiscing vitae proin. Pellentesque habitant morbi tristique senectus. Placerat orci nulla pellentesque dignissim enim sit amet venenatis. Amet consectetur adipiscing elit ut aliquam purus sit amet luctus. Ullamcorper sit amet risus nullam eget felis eget. Turpis massa tincidunt dui ut ornare lectus. Rhoncus mattis rhoncus urna neque viverra justo nec ultrices dui. Aliquam purus sit amet luctus venenatis lectus magna fringilla urna. Tincidunt eget nullam non nisi est. Varius quam quisque id diam vel quam elementum pulvinar etiam. Porttitor leo a diam sollicitudin tempor id eu nisl nunc. Id ornare arcu odio ut sem nulla pharetra diam. Sodales ut etiam sit amet.",
		"Nunc vel risus commodo viverra maecenas accumsan lacus. Sodales ut etiam sit amet nisl purus in mollis. Ultrices sagittis orci a scelerisque purus. Justo nec ultrices dui sapien eget mi proin sed libero. Ultricies tristique nulla aliquet enim tortor at auctor urna nunc. Duis tristique sollicitudin nibh sit amet. Faucibus scelerisque eleifend donec pretium vulputate sapien nec. Vel quam elementum pulvinar etiam non quam lacus suspendisse faucibus. Mauris nunc congue nisi vitae suscipit. Semper eget duis at tellus at urna condimentum. Id nibh tortor id aliquet lectus proin nibh. Eget egestas purus viverra accumsan in. Quam elementum pulvinar etiam non quam lacus suspendisse faucibus.",
		"Vel orci porta non pulvinar neque laoreet suspendisse. Maecenas pharetra convallis posuere morbi. Ac felis donec et odio pellentesque diam volutpat commodo sed. Non odio euismod lacinia at quis risus sed vulputate. Ut diam quam nulla porttitor massa id neque. Lacus sed viverra tellus in hac habitasse platea dictumst vestibulum. Sed faucibus turpis in eu mi bibendum neque. Sed risus ultricies tristique nulla. Nunc sed velit dignissim sodales ut. Ac turpis egestas integer eget aliquet nibh praesent tristique. Vitae aliquet nec ullamcorper sit amet risus nullam eget felis. Ut porttitor leo a diam sollicitudin tempor. Leo integer malesuada nunc vel risus commodo.",
		"Hendrerit dolor magna eget est lorem ipsum dolor sit. Eget est lorem ipsum dolor sit amet consectetur. Mattis rhoncus urna neque viverra. Convallis posuere morbi leo urna molestie at elementum eu. Ac ut consequat semper viverra. Tincidunt augue interdum velit euismod in pellentesque massa placerat. In cursus turpis massa tincidunt dui ut. Sapien eget mi proin sed libero enim sed faucibus turpis. Vel pretium lectus quam id leo in vitae turpis. Faucibus turpis in eu mi bibendum neque egestas congue. Nulla pharetra diam sit amet nisl suscipit adipiscing bibendum est. Pellentesque sit amet porttitor eget dolor morbi non arcu risus. Laoreet sit amet cursus sit amet dictum sit amet. Lorem ipsum dolor sit amet. Maecenas sed enim ut sem viverra.",
	}

	i := 0
	for {
		if i > 15 {
			break
		}
		line := lines[rand.Intn(len(lines))]
		if err := g.LogToFile("./testing.log", line); err != nil {
			t.Logf("error in logging a line... %v\n", err)
		}
		i++
	}
	// verification
	if wordCountSuccess("./testing.log") == 0 {
		t.Errorf("[wordCountSuccess] expect testing.log has 15 rows of data... BUT was EMPTY~")
	}
	if wordCountSuccess_2("./testing.log") == 0 {
		t.Errorf("[wordCountSuccess_2] expect testing.log has 15 rows of data... BUT was EMPTY~")
	}
	if wordCountSuccess_3("./testing.log") == 0 {
		t.Errorf("[wordCountSuccess_3] expect testing.log has 15 rows of data... BUT was EMPTY~")
	}

	// housekeep
	os.Remove("./testing.log")
}

// wordCountSuccess - success way to run wc -l on a given file.
func wordCountSuccess(file string) int {
	s := bytes.NewBuffer(nil)
	c := exec.Cmd{
		Path:   "/usr/bin/wc",
		Args:   []string{"/usr/bin/wc", "-l", file},
		Env:    []string{"PATH=/usr/bin/|"},
		Stdout: s,
		Stderr: s,
	}
	if err := c.Run(); err != nil {
		fmt.Printf("failed to run 'wc -l', reason [%v]", err)
	}
	fmt.Printf("[wordCountSuccess] results from execution wc -l -> \n%v\n", s.String())
	return s.Len()
}

// wordCountSuccess_2 - another way to run wc -l.
// tricky point is /usr/bin/wc might be required if a simple "wc" was not recognized due to PATH settings...
func wordCountSuccess_2(file string) int {
	cmd := "wc"
	//cmd = "/usr/bin/wc"
	c := exec.Command(cmd, "-l", file)
	b, err := c.Output()
	if err != nil {
		fmt.Printf("failed to run the command, %v\n", err)
	} else {
		result := string(b)
		fmt.Printf("[wordCountSuccess_2] result ->\n%v\n", result)
		return len(result)
	}
	return 0
}

func wordCountSuccess_3(file string) int {
	s := bytes.NewBuffer(nil)
	c := exec.Cmd{
		Path:   "/usr/bin/wc",                       // doesn't recognize.... need full path
		Args:   []string{"NOT_USED_:)", "-l", file}, // the 1st argument is ignored... hence you will need a "placeholder" to occupy the 1st argument...
		Env:    []string{"PATH=/usr/bin/|"},         // depends on where your command is located, the PATH variable might need to be configured
		Stdout: s,
		Stderr: s,
	}
	if err := c.Run(); err != nil {
		fmt.Printf("failed to run 'wc -l', reason [%v]", err)
	}
	fmt.Printf("[wordCountSuccess_3] results from execution wc -l -> \n%v\n", s.String())
	return s.Len()
}

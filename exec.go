package envplate

import "os"

func SplitArgs() (flagArgs []string, execArgs []string) {

	split := len(os.Args)

	for idx, e := range os.Args {

		if e == "--" {
			split = idx
			break
		}

	}

	flagArgs = os.Args[0:split]

	if split < len(os.Args) {
		execArgs = os.Args[split+1 : len(os.Args)]
	} else {
		execArgs = []string{}
	}

	return flagArgs, execArgs

}

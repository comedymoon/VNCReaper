package banner

import (
    "fmt"
    "github.com/fatih/color"
)

const BannerText = `
         _       _______ _______ _______ _______ _______ _______ _______ 
|\     /( (    /(  ____ (  ____ |  ____ (  ___  |  ____ |  ____ (  ____ )
| )   ( |  \  ( | (    \/ (    )| (    \/ (   ) | (    )| (    \/ (    )|
| |   | |   \ | | |     | (____)| (__   | (___) | (____)| (__   | (____)|
( (   ) ) (\ \) | |     |     __)  __)  |  ___  |  _____)  __)  |     __)
 \ \_/ /| | \   | |     | (\ (  | (     | (   ) | (     | (     | (\ (   
  \   / | )  \  | (____/\ ) \ \_| (____/\ )   ( | )     | (____/\ ) \ \__
   \_/  |/    )_|_______//   \__(_______//     \|/      (_______//   \__/
`

var Version = "v1.1"
var Author = "tg:@nocommand/@lostcmd ❤"

func Print() {
    cyan := color.New(color.FgCyan).SprintFunc()
    yellow := color.New(color.FgYellow).SprintFunc()
    green := color.New(color.FgGreen).SprintFunc()

    fmt.Println(cyan(BannerText))
    fmt.Printf("%s %s by %s\n\n", yellow("VNCReaper"), Version, green(Author))
}
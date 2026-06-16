# vhs tape for termiedos — render with: make demo  (or: vhs demo.tape)
Output demo.gif

Set Shell "bash"
Set FontSize 15
Set Width 1280
Set Height 760
Set Padding 18
Set Theme "Dracula"

# build + launch (hidden so the GIF starts on the loaded UI)
Hide
Type "go build -o termiedos . && ./termiedos" Enter
Sleep 4s
Show

# scores: scroll through today's games
Sleep 1500ms
Down@350ms 4
Sleep 1s

# open the selected league → standings
Enter
Sleep 3s

# switch to the fixture tab
Tab
Sleep 2500ms

# back to scores
Type "b"
Sleep 1s

# open a game → detail
Down@350ms 2
Enter
Sleep 3s

# quit
Type "b"
Sleep 800ms
Type "q"
Sleep 1s

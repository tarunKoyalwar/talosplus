
## Directives/Modules

- Refer Below Table for available directives and their use

| Directive | Syntax | Description | 
| -- | -- | -- |
| **#dir** | `#dir:/path/to/directory` | Run Given Command in this directory |
| **#from** | `#from:@varname` | Get Data from variable(`@varname`) and pass as stdin to cmd |
| **#as**| `#as:@varname` | Export Output(Stdout Only) of Command to variable (`@varname`) |
| **#for** | `#for:@arr:@i` | For each line(`@i`) in variable(`@arr`) run new command (Similar to **interlace**) |
| **#ignore** | `#ignore`  | Ignore Output of this command While showing Output | 



## Variables 

Variables are like buffers/env-variable etc starting with `@` and are handled by golang and are thread-safe . All variables exported in script are saved to MongoDB thus it is possible to get output of a specific command in the middle of execution. Talosplus tries to ignore `Everything is a file` Linux Philosophy by abstracting file system and creating and deleting files at runtime based on the need. Below Table Contains Some operations that can be performed on variables.

A Particular operation can be done on variable by supplying operator within `{}`

|Operator| Use Case | Description |
| -- | -- | -- |
| **add** | `#as:@gvar{add}` | Append Output of command to `@gvar` variable |
| **unique** | `#as:@gvar{unique}` | Append output of command to `@gvar` but all values are unique |
| **file** | `@inscope{file}` | Create a Temp File with `@inscope` variable data and return path of that temp file |
| **!file** | `@outscope{!file}` | Same as `file` but it can be empty |



## - Special Cases

| Syntax | Example | Description |
| -- | -- | -- |
| `@outfile`   |  `subfinder ... -o @outfile`  | Create a temp file(`@outfile`) and use content of file as output instead of stdout |
| `@tempfile` | - | Create a temp file and return its path  |
| `@env` | `@env:DISCORD_TOKEN` | Get value of enviournment variable (Can also be done using `$`) |



## Writing Automation Scripts With Syntax
To leverage all features of Talosplus like Auto Scheduling etc . It is essential the written bash script follows the syntax . Example of such bash script can be found at [subenum.sh](examples/subenum.sh) . 

In detail guide of how to write such scripts and using the syntax can be found at [blog](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-2-c8cd72018922)
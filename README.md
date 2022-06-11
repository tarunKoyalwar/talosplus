<p align="center" >
<img src="static/banner.jpeg" width="800" height="396" >
</br>
</p>


<p align="center">
<a href="https://opensource.org/licenses/BSD-3-Clause"><img src="https://img.shields.io/badge/license-BSD-_red.svg"></a>
<a href="https://goreportcard.com/badge/github.com/tarunKoyalwar/talosplus"><img src="https://goreportcard.com/badge/github.com/tarunKoyalwar/talosplus"></a>
<a href="https://github.com/tarunKoyalwar/talosplus/releases"><img src="https://img.shields.io/github/release/tarunKoyalwar/talosplus"></a>
<a href="https://twitter.com/KoyalwarTarun"><img src="https://img.shields.io/twitter/follow/KoyalwarTarun.svg?logo=twitter"></a>
</p>

<p align="center">
 <a href="#screenshots">Screenshots</a> •
 <a href="#blogs">Blog</a> •
  <a href="#features">Features</a> •
  <a href="#installation-instructions">Installation</a> •
  <a href="#usage">Usage</a> 
</p>


Talosplus is tool to run bash scripts at faster rate by executing commands in parallel using goroutines and with some exceptional features like **Auto Scheduling, Filesystem Abstraction ,Stop/Resume, Buffers,Thread Safe ,Fail Safe, Serial + Parallel Execution, Notification Support** etc provided that script follows given Syntax and is integrated with **MongoDB** which provides lot of flexiblity similar to `bbrf` especially for Bug Hunters, Security Professionals etc.

# Blog / How To Guides

[Create Your Ultimate Bug Bounty Automation Without Nerdy Bash Skills](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-1-a78c2b109731)

- [Part 1](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-1-a78c2b109731)

- [Part 2](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-2-c8cd72018922)

- [Part 3](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-3-7ee2b353a781)

# Why ??

Why use this when bash scripts can be run directly ?? You can think of this like a middleware to run bash scripts . I wanted to create a perfect automation much like **@hakluke** . This project resolves all challenges and issues I faced while writing bash scripts and creating the perfect automation and makes it possible to leverage all important features with comments `Ex: #as:@nmapout, #from:@allsubs etc` . and adds a lot of additional features.

Even If you are a little intriqued, Consider reading [my blog](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-1-a78c2b109731) . Which describes how I overcame challenges I faced , how and when to use these comments ? and effective use of this project and detailed description of all its features like scheduling algo etc.


~~~
If you don't want to use of these comments  or features . Supplying your regular bash script
Will run every command it can find in parallel.
~~~

# Screenshots


- Sample Bash Script with Syntax at [here](static/script.png)

- Talosplus output at [here](static/cmdout.png)

- Custom Discord Notification at [here](static/notification.png)


# Features

These are oversimplified features to name from my blog.

- Auto Scheduling Commands at Runtime
- Intelligent Automation
- Filesystem Abstraction
- Discord Notification Support
- Thread Safe
- All Features of BBRF+ Others (MongoDB Backend)
- Easy Syntax
- Fail Safe && Condition Checks
- Stop /Resume(BETA) 
- No Compatiblity issues


The driving forces behind talosplus are **variables** and **directives** . These directives and variables abstract complex bash syntaxes and solve challanges with little syntax.

## Directives

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



- Special Cases

| Syntax | Example | Description |
| -- | -- | -- |
| `@outfile`   |  `subfinder ... -o @outfile`  | Create a temp file(`@outfile`) and use content of file as output instead of stdout |
| `@tempfile` | - | Create a temp file and return its path  |
| `@env` | `@env:DISCORD_TOKEN` | Get value of enviournment variable (Can also be done using `$`) |

# Installation Instructions

- Configure MongoDB Atlas or Install [MongoDB](https://www.mongodb.com/docs/manual/installation/).

- Install `libx11-dev` (Provides Clipboard Access)
  
  - On Debian Based distro ```sudo  apt install libx11-dev```

  - On ArchLinux Based distro ```sudo pacman -S libx11```


- Build From Source .

~~~sh
go install github.com/tarunKoyalwar/talosplus/cmd/talosplus@latest
~~~


Do Star the repo to show  your support.
Follow me on [github](https://github.com/tarunKoyalwar) / [twitter](https://twitter.com/KoyalwarTarun) to get latest updates on Talosplus.


Refer to Blog [Part 3](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-3-7ee2b353a781) for step by step instructions on using `talosplus` command in detail with examples.


# Limitations

1. Taloplus is just a parser tool and is not aware of bash syntax .

2. Each Command is sandboxed if you are using bash environment variables etc it won't work .It has to be variables

3. For Loops, IF etc Will Work But they can only be in a single line or newline should be escaped using `\`.


Saving Outputs to File/Environment Variables Entirely Defeats Purpose of This tool .
Read Blog or Refer to subenum.sh file before running any script file.


# Usage


Check Below Sample Video which Shows How I use talosplus for Subdomain Enumeration Automation using [subenum.sh](/examples/subenum.sh) 


[![asciicast](https://asciinema.org/a/qHeRefcO6WOPrWuNAnpcuICLf.svg)](https://asciinema.org/a/qHeRefcO6WOPrWuNAnpcuICLf)




Talosplus has every feature that would make it easy to write and run bash scripts . 

## Writing Automation Scripts With Syntax
To leverage all features of Talosplus like Auto Scheduling etc . It is essential the written bash script follows the syntax . Example of such bash script can be found at [subenum.sh](examples/subenum.sh) . 

In detail guide of how to write such scripts and using the syntax can be found at [blog](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-2-c8cd72018922)

## Usage in Detail

In depth details on running any scripts ,configs , interacting with db , storing and retrieving any subdomain(or any variable from bash script) etc. and much more can be found at [blog](https://medium.com/@zealousme/create-your-ultimate-bug-bounty-automation-without-nerdy-bash-skills-part-3-7ee2b353a781)

# Support

If you like `talosplus` and want to see it improve furthur or want me to create intresting projects , You can buy me a coffee 

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/B0B4CPU5V)

# Acknowledgment

Some Features are inspired by [@honoki/bbrf-client](https://github.com/honoki/bbrf-client)
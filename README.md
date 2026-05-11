# GitLab-User-Ripper
A GitLab user enumeration tool that uses concurrency to rip

> DISCLAIMER: This tool is intended for educational use only. Do not use this tool to attack any instance of GitLab you are not authorized to test, including gitlab.com.

This tool is inspired by [dpgg101/GitLabUserEnum](https://github.com/dpgg101/GitLabUserEnum), which is written in Python, and created for completing the *Attack GitLab* Section in the HTB Academy module **Attack Common Applications**.

GitLab allows for users on the instance to be enumerated by sending a HEAD HTTP request to `<URL>/<username>`. `200 OK` is return if the user exists, and `302 FOUND` is returned otherwise. This tool speeds up the process by using Go's worker pool model of concurrency, allowing multiple requests to be sent out simultaneously, drastically increasing the enumeration speed.

## Build
```
git clone https://github.com/RandomChugokujin/GitLab-User-Ripper.git
cd GitLab-User-Ripper
go build
```

## Usage
```
$ gitlab-user-ripper -h

 _____ _ _   __        _       _____                 _____ _
|   __|_| |_|  |   ___| |_ ___|  |  |___ ___ ___ ___| __  |_|___ ___ ___ ___
|  |  | |  _|  |__| .'| . |___|  |  |_ -| -_|  _|___|    -| | . | . | -_|  _|
|_____|_|_| |_____|__,|___|   |_____|___|___|_|     |__|__|_|  _|  _|___|_|
                                                            |_| |_|
[*] A GitLab user enumeration tool that rips
[*] Author: RandomChugokujin

Usage of ./gitlab-user-ripper:
  -f string
        Path to username file
  -t int
        Number of threads (default 50)
  -u string
        Base URL to scan (e.g., http://gitlab.local:8081)
  -v    Verbose Output
```

Example:
```
gitlab-user-ripper -u http://gitlab.inlanefreight.local -f /usr/share/seclists/Usernames/xato-net-10-million-usernames.txt -w 100
```

## Contributing
If you find any issues with the code or have suggestions for improvements, please feel free to open an issue or submit a pull request.

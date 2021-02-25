# go-evilarc
go-evilarc lets you create a zip file that contains files with directory traversal characters in their embedded path. Most commercial zip program (winzip, etc) will prevent extraction of zip files whose embedded files contain paths with directory traversal characters. However, many software development libraries do not include these same protection mechanisms (ex. Java, PHP, etc). If a program and/or library does not prevent directory traversal characters then go-evilarc can be used to generate zip files that, once extracted, will place a file at an arbitrary location on the target system.

## Usage
```bash
go-evilarc version v0.0.1
Usage: go-evilarc <input file>

Create archive containing a file with directory traversal

Options:
	-out <filename>         File to output archive to. Archive type is based off of file extension.
	                        Supported extesions are zip, jar, tar, tar.bz2, tar.gz and tgz
	-depth <int>            Number of directories to traverse (default: 8)
	-platform [win|unix]    OS platform for archive (default: win)
	-trav                   You can define a custom traversal vector to use
	-path                   Path to include in filename after traversal.
                            Ex: WINDOWS\\System32\\ or var/www/
```
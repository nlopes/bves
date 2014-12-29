This is the source code repository for bves.

bves is a command line utility to interface with Virtualbox.es boxes

## Usage
    usage: bves [<flags>] <command> [<flags>] [<args> ...]

    A cmdline interface to virtualbox.es.

    Flags:
      --help             Show help.
      -d, --debug        Enable debug mode.
      -t, --timeout=10s  Timeout for download.
      --version          Show application version.

    Commands:
      help [<command>]
        Show help for a command.

      list
        List boxes.

      show <id>
        Show box information.

      url <id>
        Show box URL

      clear-cache
        Clear local cache.


## Examples

`# bves list | head -n 10`

     1 Debian 7.3.0 64-bit Puppet 3.4.1 (Vagrant 1.4.0)
     2 OpenBSD 5.5 64-bit + Chef 11.16.0 + Puppet 3.4.2
     3 OpenBSD 5.4 64-bit + Chef 11.8.2 (150GB HDD)
     4 OpenBSD 5.3 64-bit (Vagrant 1.2)
     5 OpenBSD 5.3 64-bit
     6 Aegir-up Aegir (Debian Squeeze 6.0.4 64-bit)
     7 Aegir-up Debian (Debian Squeeze 6.0.4 64-bit)
     8 Aegir-up LAMP (Debian Squeeze 6.0.4 64-bit)
     9 AppScale 1.12.0 (Ubuntu Precise 12.04 64-bit)
    10 Arch Linux 64 (2014-06-20)

`# bves show 260`

    Name:     Windows 8.1 with IE11 (32bit)
    Details:  Windows 8.1 with IE11 (32bit)
	  The Microsoft Software License Terms for the IE VMs are included in the release notes and supersede any conflicting Windows license terms included in the VMs. By downloading and using this software, you agree to these license terms.
    Provider: Virtualbox
    URL:      http://aka.ms/vagrant-win81-ie11
    Size:     3584.0 MB


## How is this tool useful?
If like me you want to avoid going to the browser, you can plug bves into vagrant:

`# vagrant box add Debian7 $(./bves url 1)`

    Downloading or copying the box...
    Progress: 1% (Rate: 808k/s, Estimated time remaining: 0:06:03)


## License
BSD 2 Clause license

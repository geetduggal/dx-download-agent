# dx-download-agent
CLI tool to manage the download of large quantities of files from DNAnexus

**WARNING: This is an alpha version of this tool. It is currently in a specification/draft stage and it is likely incomplete. Please use at your own risk.**

## Quick Start

To get started with `dx-download-agent`, download the the latest pre-compiled binary from the release page.  The download agent accepts two files:

* `manifest_file`: A BZ2-compressed JSON manifest file that describes, at minimimum, a list of DNAnexus file IDs to download

For example, in the current working directory let an uncompressed `manifest.json.bz2` be:

```json
{ 
  "project-AAAA": [
    {
      "id": "file-XXXX",
      "name": "foo",
      "folder": "/path/to",
      "parts": [
        { "size": 10, "md5": 49302323 },
        { "size": 5,  "md5": 39239329 }
      ]
    },
    "..."
  ],
  "project-BBBB": [ "..." ]
}
```

To start a download process, first [generate a DNAnexus API token](https://wiki.dnanexus.com/Command-Line-Client/Login-and-Lgout#Authentication-Tokens) that is valid for a time period that you plan on downloading the files.  Store it in the following environment variable:

```bash
export DX_API_TOKEN=<INSERT API TOKEN HERE>
```

In the same directory, begin the download process with this command:

```
dx-download-agent probe-environment
```

This command will perfrom a series of initial checks but avoid downloads.  These checks include:

* Network connectivity and potential issues with it
* Whether you have enough space locally
* Approximate speeds of download rates
* Whether it looks like another download process is running (i.e. file sizes are changing, status files being updated).

```
dx-download-agent download exome_bams.manifest.json.bz2
Rate: 49.9 KBp/s
```

This command will also probe the environment and, if it doesn't appear another download process is running, it will start a download process within your terminal using the current working directory.

Once a download has begun, in a separate terminal in the same directory type:

```
dx-download-agent progress
```

and you will get a brief summary of the status the downloads:

```
95/1056 Mb (9%)
```

## Execution options

* `--max_threads` (integer): maximum # of concurrent threads to use when downloading files
* `--max_bandwidth`: (integer) Maximum download bandwidth of all files in KBp/s
* ...


## Additional notes

* Only objects of [class File](https://wiki.dnanexus.com/API-Specification-v1.0.0/Introduction-to-Data-Object-Classes) can be downloaded. 
* On DNAnexus, files are immutable and the same directory can contain multiple files of the same name.  If this occurs, files on a local POSIX filesystem will be appended with the DNAnexus file ID to ensure they are not overwritten.  
* In the case a directory and a file have the same name and share the same parent directory, a DNAnexus file ID will also be appended.  If the file name contains at least one character that is illegal on a POSIX system, the file will be named directly by its file ID on DNAnexus.
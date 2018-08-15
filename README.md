# dx-download-agent
CLI tool to manage the download of large quantities of files from DNAnexus

**WARNING: This is an alpha version of this tool. It is currently in a specification/draft stage and it is likely incomplete. Please use at your own risk.**

## Quick Start

To get started with `dx-download-agent`, download the the latest pre-compiled binary from the release page.  The download agent accepts two files with the following name:

* `manifest_file`: A YAML manifest file that describes, at minimimum, a list of DNAnexus file IDs to download
* `execution_params`: A YAML file describing parameters used for a particular execution of the download agent

For example, in the current working directory let `manifest.yaml` be:

```yaml
- id: file-XXXX
- id: file-YYYY
- id: file-ZZZZ
```
and let `execution.yaml` be:

```yaml
# If no execution YAML is provided, defaults are used.
# See "Execution file format" section below for more information

preserve_project_paths: true  # Mimics directory structure in project locally
max_threads: 4
max_bandwidth: 50 # KBp/s
```

To start a download process, first [generate a DNAnexus API token](https://wiki.dnanexus.com/Command-Line-Client/Login-and-Logout#Authentication-Tokens) that is valid for a time period that you plan on downloading the files.  Store it in the following environment variable:

```bash
export DX_API_TOKEN=<INSERT API TOKEN HERE>
```

In the same directory, begin the download process with this command:

```
dx-download-agent start .
Rate: 49.9 KBp/s
```

This command will start a download process within your terminal using the current working directory (the '.' in the command above).  If you wish to execute this command outside of the current directory and provide pointers to configuration files in a different location, you can execute the command like:

```
dx-download-agent start . --manifest_file=/path/to/manifest.yaml --execution_params=/path/to/execution.yaml
```


Once a download has begun, in a separate terminal, type:

```
dx-download-agent status .
```

and you will get a brief summary of the status the downloads:

```
Current state of download: RUNNING

Progress        # Files
-----------------------
  0-10%               1
 11-20%               1
 21-30%               0
 31-40%               0
 41-50%               1
 51-60%               0
 61-70%               0
 71-80%               0
 81-90%               0
 91-99%               0
 ----------------------
 COMPLETED            0
 PERMANENTLY FAILED   0
 AVERAGE RATE        45 KBp/s
```

For more detailed information on a per-file basis, the following command:

```
dx-download-agent status . --detailed-tsv
```

 will output a tab-separated file format to standard output with the following columns:

* DNAnexus file ID
* DNAnexus file name
* Status: in-progress, completed, permanently failed
* Percent downloaded
* Size downloaded in gigabytes
* Local path (relative to working directory)
* Additional info (e.g. reason for permanent failure)


## Managing downloads over long periods of time

Since downloading these files could occur across many days or even weeks, here are a few ways to manage the process over, for example, reboots of machines and other kinds of system maintenance.  If `dx-download-agent status .` returns:

```
Current state of download: STOPPED
```

then this means a new `start` command must be issued for the downloads to resume from where they left off.

After, for example, a system reboot, the user may wish to manually check on the status of the download and issue the `start` command again.  Optionally, the user can also create a system service for the download (e.g. via a [SystemD service file](https://www.devdungeon.com/content/creating-systemd-service-files)) to automatically handle restarts on failures and at system boot time.

## Manifest file format

The manifest file is a YAML file that is a list of dictionaries. Each dictionary item must have an `id` field at minimum as in the example above, but each item can also have other metadata:

* `id` (required): DNAnexus file ID
* `name`: File name on DNAnexus (used by `status` command)
* `desired_local_path`: Desired local path
* ...

## Execution file format

The execution file is completely optional

* `preserve_project_paths` (boolean, default true):  Mimic directory structure in project locally.  This overrides `desired_local_path` option.
* `max_threads` (integer): maximum # of concurrent threads to use when downloading files
* `max_bandwidth`: (integer) Maximum download bandwith of all files in KBp/s
* ...


## Additional notes

On DNAnexus, files are immutable and the same directory can contain multiple files of the same name.  If this occurs, files on a local POSIX filesystem will be appended with the DNAnexus file ID to ensure they are not overwritten.

# Technologies used

* [Go Grab](https://github.com/cavaliercoder/grab)

# Developing for dx-download-agent

Instructions for compilation via Go, adding tests, etc.
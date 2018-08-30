import argparse
import dxpy
from pprint import pprint
import json
import bz2


def fileID2manifest(fid, project):
    """
    Convert a single file ID to an entry in the manifest file
    Inputs: DNAnexus file and project ID
    Output: dictionary corresponding to manifest entry
    """

    fdetails = dxpy.api.file_describe(fid, input_params={'fields': {'id': True, 'name': True, 'folder': True, 'parts': True }})
    fdetails['parts'] = {pid: {k:v for k,v in pdetails.items() if k == "md5" or k == "size"} for pid, pdetails in fdetails['parts'].items()}
    return fdetails


def main():
    parser = argparse.ArgumentParser(description='Create a manifest file from a list of DNAnexus file IDs')
    parser.add_argument('id', nargs='+', help='a DNAnexus file ID')
    parser.add_argument('--project', help='Project ID: required to speed up API calls', required=True)
    parser.add_argument('--outfile', help='Name of the output file', default='manifest.json.bz2')


    args = parser.parse_args()

    manifest = {args.project: [fileID2manifest(fid, args.project) for fid in args.id]}
    with open(args.outfile, "w") as f:
        f.write(bz2.compress(json.dumps(manifest, indent=2, sort_keys=True)))

    print("Manifest file written to {}".format(args.outfile))

if __name__ == "__main__":
    main()

all: single_file

single_file.bin:
	dd if=/dev/random of=single_file.bin bs=1000000 count=100
	# Upload file to platform so that there are 10ish parts
	# dx upload --buffer-size 10000000 single_file.bin

single_file.manifest.json.bz2:
	python fileIDs2manifest.py file-FK41Q8Q0BbkqJxkX35p736BJ --project project-FK411f80BbkqJxkX35p73660 --outfile single_file.manifest.json.bz2

single_file: single_file.bin single_file.manifest.json.bz2

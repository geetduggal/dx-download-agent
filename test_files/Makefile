all: single_file

single_file.bin:
	dd if=/dev/random of=single_file.bin bs=1000000 count=100
	# Upload file to platform so that there are 10ish parts
	# dx upload --buffer-size 10000000 single_file.bin

other_file.bin:
	dd if=/dev/random of=other_file.bin bs=1000000 count=100

single_file.manifest.json.bz2:
	python fileIDs2manifest.py file-FK41Q8Q0BbkqJxkX35p736BJ --project project-FK411f80BbkqJxkX35p73660 --outfile single_file.manifest.json.bz2

single_file: single_file.bin single_file.manifest.json.bz2

two_files.manifest.json.bz2:
	python fileIDs2manifest.py file-FK41Q8Q0BbkqJxkX35p736BJ file-FK8j8100BbkjF42j4FBkg6qg --project project-FK411f80BbkqJxkX35p73660 --outfile two_files.manifest.json.bz2

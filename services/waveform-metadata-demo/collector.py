import os

from obspy.signal.quality_control import MSEEDMetadata

# IO Directories
INPUT = "/var/input"
OUTPUT = "/var/output/metadata.json"
ERROR = "/var/output/error.log"

# Get the filename from the input directory
files = [
  os.path.join(INPUT, file) for file in os.listdir(INPUT)
]

# Calculate the metadata
try:

  MD = MSEEDMetadata(
    files,
    add_c_segments=True,
    add_flags=True
  )

  # Write result to file
  with open(OUTPUT, "w") as outfile:
    outfile.write(MD.get_json_meta());

except Exception as e:
  
  # Write error to file
  with open(ERROR, "w") as outfile:
    outfile.write("Could not complete metric calculation: %s" % e);

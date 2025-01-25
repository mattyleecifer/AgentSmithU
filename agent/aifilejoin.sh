#!/bin/bash

# Function to escape special characters for XML
escape_xml() {
    sed 's/&/\&amp;/g; s/</\&lt;/g; s/>/\&gt;/g; s/"/\&quot;/g; s/'"'"'/\&#39;/g'
}

# Check if any files were passed as arguments
if [ $# -eq 0 ]; then
    read -p "Enter filenames separated by spaces (or press enter to skip): " filenames
else
    filenames="$@"
fi

# Split into an array in case there are multiple files
IFS=' ' read -ra file_list <<< "$filenames"

xml_content="<data>"

for filename in "${file_list[@]}"; do
    if [ -f "$filename" ]; then
        content=$(<"$filename")
        escaped_content=$(echo "$content" | escape_xml)
        xml_content+=$(cat << EOF
    <file name="$filename">
        $escaped_content
    </file>
EOF
)
    else
        echo "Warning: File '$filename' does not exist."
    fi
done

# Get request from user
read -p "Enter the request: " request
xml_content+=$(cat << EOF
    <request>
        $(echo "$request" | escape_xml)
    </request>
EOF
)

xml_content+="</data>"

# Output to a file, defaulting to output.xml if not provided
output_file="output.xml"
if [ -n "$1" ] && [ "${file_list[0]}" != "$1" ]; then
    # If the last argument is intended as the output file and it's different from the first file
    output_file="$1"
else
    read -p "Enter output filename (default: output.xml): " custom_output
    if [ -n "$custom_output" ]; then
        output_file="$custom_output"
    fi
fi

echo "$xml_content" > "$output_file"

echo "XML content has been written to $output_file"

### queryFS - recursively query a Linux/Unix file system by substring or permissions


###### Get the project, :
<pre>
  <code>
# Clone project
git clone https://github.com/rootVIII/queryFS.git

# Build and run:
cd &lt;project root&gt;
go build -o bin/queryfs
./bin/queryfs 

# Build binary in ~/go/bin (available in path) and run
cd &lt;project root&gt;
go install .
queryfs
  </code>
</pre>


###### Options:
<pre>
  <code>
# Required
-d     directory path to start searching from

# Optional (provide at least 1 option)
-s     display files/directories containing string
-p     display files/directories with matching permissions
-o     display files/directories with matching owner:group
  </code>
</pre>


###### Example Usage:
<pre>
  <code>
#  Search /home and all subdirectories for any file or directory path containing text &#34;my_file.txt&#34;
queryfs -d /home -s my_file.txt
  </code>
</pre>

<pre>
  <code>
#  Search entire file-system for any file with 777 permissions
queryfs -d / -p -rwxrwxrwx
  </code>
</pre>


<pre>
  <code>
#  Search entire file-system for any directory with 777 permissions
queryfs -d / -p drwxrwxrwx
  </code>
</pre>


<pre>
  <code>
#  Search /var and all subdirectories for any file with owner/group apache:apache
queryfs -d /var -o apache:apache
  </code>
</pre>


<pre>
  <code>
#  Only print results that have BOTH &#34;.py&#34; in path AND 0755 permissions:
queryfs -d /var -s .py -p -rwxr-xr-x
  </code>
</pre>


<pre>
  <code>
# Only print results that have &#34;.py&#34; in path, 755 permissions, AND apache:apache owner/group:
queryfs -d /var -s .py -p -rwxr-xr-x -o apache:apache
  </code>
</pre>
<hr>
This project was developed on Ubuntu 18.04.4 LTS
